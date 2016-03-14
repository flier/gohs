/*
 * Hyperscan example program 2: pcapscan
 *
 * This example is a very simple packet scanning benchmark. It scans a given
 * PCAP file full of network traffic against a group of regular expressions and
 * returns some coarse performance measurements.  This example provides a quick
 * way to examine the performance achievable on a particular combination of
 * platform, pattern set and input data.
 *
 * Build instructions:
 *
 *     go build github.com/flier/gohs/examples/pcapscan
 *
 * Usage:
 *
 *     ./pcapscan [-n repeats] <pattern file> <pcap file>
 *
 * We recommend the use of a utility like 'taskset' on multiprocessor hosts to
 * pin execution to a single processor: this will remove processor migration
 * by the scheduler as a source of noise in the results.
 *
 */
package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"hash/fnv"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/flier/gohs/hyperscan"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	repeatCount = flag.Int("n", 1, "Repeating PCAP scan several times")
)

type FiveTuple struct {
	protocol         layers.IPProtocol
	srcAddr, dstAddr net.IP
	srcPort, dstPort uint16
	hash             uint64
}

func (t *FiveTuple) Hash() uint64 {
	if t.hash == 0 {
		h := fnv.New64a()

		binary.Write(h, binary.BigEndian, t.protocol)
		h.Write(t.srcAddr)
		h.Write(t.dstAddr)
		binary.Write(h, binary.BigEndian, t.srcPort)
		binary.Write(h, binary.BigEndian, t.dstPort)

		t.hash = h.Sum64()
	}

	return t.hash
}

type Benchmark struct {
	dbStreaming hyperscan.StreamDatabase // Hyperscan compiled database (streaming mode)
	dbBlock     hyperscan.BlockDatabase  // Hyperscan compiled database (block mode)
	scratch     hyperscan.Scratch        // Hyperscan temporary scratch space (used in both modes)
	packets     [][]byte                 // Packet data to be scanned.
	streamIds   []int                    // The stream ID to which each packet belongs
	streamMap   map[uint64]int           // Map used to construct stream_ids
	streams     []hyperscan.Stream       // Vector of Hyperscan stream state (used in streaming mode)
	matchCount  int                      // Count of matches found during scanning
}

func NewBenchmark(streaming hyperscan.StreamDatabase, block hyperscan.BlockDatabase) (*Benchmark, error) {
	scratch, err := hyperscan.NewScratch(streaming)

	if err != nil {
		return nil, fmt.Errorf("could not allocate scratch space, %s", err)
	}

	if err := scratch.Realloc(block); err != nil {
		return nil, fmt.Errorf("could not reallocate scratch space, %s", err)
	}

	return &Benchmark{
		dbStreaming: streaming,
		dbBlock:     block,
		scratch:     scratch,
		streamMap:   make(map[uint64]int),
	}, nil
}

func (b *Benchmark) ReadStreams(pcapFile string) (int, error) {
	h, err := pcap.OpenOffline(pcapFile)

	if err != nil {
		return 0, err
	}

	defer h.Close()

	s := gopacket.NewPacketSource(h, h.LinkType())
	count := 0

	for pkt := range s.Packets() {
		count += 1

		var key FiveTuple

		switch t := pkt.NetworkLayer().(type) {
		case *layers.IPv4:
			key.protocol = t.Protocol
			key.srcAddr = t.SrcIP
			key.dstAddr = t.DstIP
		case *layers.IPv6:
			key.protocol = t.NextHeader
			key.srcAddr = t.SrcIP
			key.dstAddr = t.DstIP
		default:
			continue
		}

		switch t := pkt.TransportLayer().(type) {
		case *layers.TCP:
			key.srcPort = uint16(t.SrcPort)
			key.dstPort = uint16(t.DstPort)
		case *layers.UDP:
			key.srcPort = uint16(t.SrcPort)
			key.dstPort = uint16(t.DstPort)
		default:
			continue
		}

		id := len(b.streamMap)
		hash := key.Hash()

		if _id, exists := b.streamMap[hash]; exists {
			id = _id
		} else {
			b.streamMap[hash] = id
		}

		b.packets = append(b.packets, pkt.TransportLayer().LayerPayload())
		b.streamIds = append(b.streamIds, id)
	}

	return count, nil
}

func (b *Benchmark) Close() {
	// Free scratch region
	b.scratch.Free()

	b.dbStreaming.Close()
	b.dbBlock.Close()
}

func (b *Benchmark) Bytes() (sum int) {
	for _, pkt := range b.packets {
		sum += len(pkt)
	}

	return
}

func (b *Benchmark) Matches() int { return b.matchCount }

func (b *Benchmark) ClearMatches() { b.matchCount = 0 }

func (b *Benchmark) DisplayStats() {
	numPackets := len(b.packets)
	numStreams := len(b.streamMap)
	numBytes := b.Bytes()

	fmt.Printf("%d packets in %d streams, totalling %d bytes\n", numPackets, numStreams, numBytes)
	fmt.Printf("Average packet length: %d bytes\n", numBytes/numPackets)
	fmt.Printf("Average stream length: %d bytes\n", numBytes/numStreams)

	if size, err := b.dbStreaming.Size(); err != nil {
		fmt.Printf("Error getting streaming mode Hyperscan database size, %s\n", err)
	} else {
		fmt.Printf("Streaming mode Hyperscan database size : %d bytes.\n", size)
	}

	if size, err := b.dbBlock.Size(); err != nil {
		fmt.Printf("Error getting block mode Hyperscan database size, %s\n", err)
	} else {
		fmt.Printf("Block mode Hyperscan database size : %d bytes.\n", size)
	}

	if size, err := b.dbStreaming.StreamSize(); err != nil {
		fmt.Printf("Error getting stream state size, %s\n", err)
	} else {
		fmt.Printf("Streaming mode Hyperscan stream state size : %d bytes (per stream).\n", size)
	}
}

func (b *Benchmark) onMatch(ctxt hyperscan.MatchContext, evt hyperscan.MatchEvent) error {
	b.matchCount += 1

	return nil
}

func (b *Benchmark) OpenStreams() {
	b.streams = make([]hyperscan.Stream, len(b.streamMap))

	for i := 0; i < len(b.streamMap); i++ {
		stream, err := b.dbStreaming.Open(0, b.scratch, hyperscan.MatchHandFunc(b.onMatch), nil)

		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Unable to open stream, %s. Exiting.", err)
			os.Exit(-1)
		}

		b.streams[i] = stream
	}
}

func (b *Benchmark) CloseStreams() {
	for _, stream := range b.streams {
		if err := stream.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Unable to close stream, %s. Exiting.", err)
			os.Exit(-1)
		}
	}
}

func (b *Benchmark) ScanStreams() {
	for i, pkt := range b.packets {
		if len(pkt) == 0 {
			continue
		}

		if err := b.streams[b.streamIds[i]].Scan(pkt); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Unable to scan packet, %s. Exiting.", err)
			os.Exit(-1)
		}
	}
}

func (b *Benchmark) ScanBlock() {
	for _, pkt := range b.packets {
		if len(pkt) == 0 {
			continue
		}

		if err := b.dbBlock.Scan(pkt, b.scratch, hyperscan.MatchHandFunc(b.onMatch), nil); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Unable to scan packet, %s. Exiting.", err)
			os.Exit(-1)
		}
	}
}

type Clock struct {
	start, stop time.Time
}

func (c *Clock) Start() { c.start = time.Now() }

func (c *Clock) Stop() { c.stop = time.Now() }

func (c *Clock) Time() time.Duration { return c.stop.Sub(c.start) }

func parseFile(filename string) (patterns []*hyperscan.Pattern) {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Can't open pattern file %s\n", filename)
		os.Exit(-1)
	}

	reader := bufio.NewReader(bytes.NewBuffer(data))
	eof := false
	lineno := 0

	for !eof {
		line, err := reader.ReadString('\n')

		switch err {
		case nil:
			// pass
		case io.EOF:
			eof = true
		default:
			fmt.Fprintf(os.Stderr, "ERROR: Can't read pattern file %s, %s\n", filename, err)
			os.Exit(-1)
		}

		line = strings.TrimSpace(line)
		lineno += 1

		// if line is empty, or a comment, we can skip it
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		// otherwise, it should be ID:PCRE, e.g.
		//  10001:/foobar/is
		strs := strings.SplitN(line, ":", 2)

		id, err := strconv.ParseInt(strs[0], 10, 64)

		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Could not parse id at line %d, %s", lineno, err)
			os.Exit(-1)
		}

		pattern, err := hyperscan.ParsePattern(strs[1])

		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Could not parse pattern at line %d, %s", lineno, err)
			os.Exit(-1)
		}

		pattern.Id = int(id)

		patterns = append(patterns, pattern)
	}

	return
}

/**
 * This function will read in the file with the specified name, with an
 * expression per line, ignoring lines starting with '#' and build a Hyperscan
 * database for it.
 */
func databasesFromFile(filename string) (hyperscan.StreamDatabase, hyperscan.BlockDatabase) {
	// do the actual file reading and string handling
	patterns := parseFile(filename)

	fmt.Printf("Compiling Hyperscan databases with %d patterns.\n", len(patterns))

	var clock Clock

	clock.Start()

	sdb, err := hyperscan.NewStreamDatabase(patterns...)

	clock.Stop()

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Could not compile patterns, %s", err)
		os.Exit(-1)
	}

	fmt.Printf("Hyperscan streaming mode database compiled in %.2f ms\n", clock.Time().Seconds()*1000)

	clock.Start()

	bdb, err := hyperscan.NewBlockDatabase(patterns...)

	clock.Stop()

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Could not compile patterns, %s", err)
		os.Exit(-1)
	}

	fmt.Printf("Hyperscan block mode database compiled in %.2f ms\n", clock.Time().Seconds()*1000)

	return sdb, bdb
}

func main() {
	flag.Parse()

	if flag.NArg() != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-n repeats] <pattern file> <pcap file>\n", os.Args[0])
		os.Exit(-1)
	}

	patternFile := flag.Arg(0)
	pcapFile := flag.Arg(1)

	// Read our pattern set in and build Hyperscan databases from it.
	fmt.Printf("Pattern file: %s\n", patternFile)

	dbStreaming, dbBlock := databasesFromFile(patternFile)

	bench, err := NewBenchmark(dbStreaming, dbBlock)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s, Exiting.", err)
		os.Exit(-1)
	}

	defer bench.Close()

	fmt.Printf("PCAP input file: %s\n", pcapFile)

	if read, err := bench.ReadStreams(pcapFile); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Unable to read packets from PCAP file, %s. Exiting.", err)
		os.Exit(-1)
	} else {
		fmt.Printf("Read %d of %d packets within %d stream\n", len(bench.packets), read, len(bench.streamMap))
	}

	if *repeatCount != 1 {
		fmt.Printf("Repeating PCAP scan %d times.\n", *repeatCount)
	}

	bench.DisplayStats()

	var clock Clock
	var secsStreamingScan, secsStreamingOpenClose time.Duration

	// Streaming mode scans.
	for i := 0; i < *repeatCount; i++ {
		// Open streams.
		clock.Start()
		bench.OpenStreams()
		clock.Stop()

		secsStreamingOpenClose += clock.Time()

		// Scan all our packets in streaming mode.
		clock.Start()
		bench.ScanStreams()
		clock.Stop()

		secsStreamingScan += clock.Time()

		// Close streams.
		clock.Start()
		bench.CloseStreams()
		clock.Stop()
		secsStreamingOpenClose += clock.Time()
	}

	// Collect data from streaming mode scans.
	bytes := bench.Bytes()
	tputStreamScanning := float64(bytes*8**repeatCount) / secsStreamingScan.Seconds()
	tputStreamOverhead := float64(bytes*8**repeatCount) / (secsStreamingScan + secsStreamingOpenClose).Seconds()
	matchesStream := bench.Matches()
	matchRateStream := float64(matchesStream) / (float64(bytes**repeatCount) / 1024.0) // matches per kilobyte

	// Scan all our packets in block mode.
	bench.ClearMatches()

	clock.Start()
	for i := 0; i < *repeatCount; i++ {
		bench.ScanBlock()
	}
	clock.Stop()

	secsScanBlock := clock.Time()

	// Collect data from block mode scans.
	tputBlockScanning := float64(bytes*8**repeatCount) / secsScanBlock.Seconds()
	matchesBlock := bench.Matches()
	matchRateBlock := float64(matchesBlock) / (float64(bytes**repeatCount) / 1024.0) // matches per kilobyte

	fmt.Println("Streaming mode:")
	fmt.Printf("  Total matches: %d\n", matchesStream)
	fmt.Printf("  Match rate:    %.4f matches/kilobyte\n", matchRateStream)
	fmt.Printf("  Throughput (with stream overhead): %.2f megabits/sec\n", tputStreamOverhead/1000000)
	fmt.Printf("  Throughput (no stream overhead):   %.2f megabits/sec\n", tputStreamScanning/1000000)
	fmt.Println("Block mode:")
	fmt.Printf("  Total matches: %d\n", matchesBlock)
	fmt.Printf("  Match rate:    %.4f matches/kilobyte\n", matchRateBlock)
	fmt.Printf("  Throughput:    %.2f megabits/sec\n", tputBlockScanning/1000000)

	if bytes < (2 * 1024 * 1024) {
		fmt.Println("")
		fmt.Println("WARNING: Input PCAP file is less than 2MB in size.")
		fmt.Println("This test may have been too short to calculate accurate results.")
	}
}
