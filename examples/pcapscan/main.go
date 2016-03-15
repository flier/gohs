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
	"runtime/pprof"
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
	cpuprofile  = flag.String("cpuprofile", "", "write cpu profile to file")
	memprofile  = flag.String("memprofile", "", "write memory profile to this file")
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

// Class wrapping all state associated with the benchmark
type Benchmark struct {
	dbStreaming hyperscan.StreamDatabase // Hyperscan compiled database (streaming mode)
	dbBlock     hyperscan.BlockDatabase  // Hyperscan compiled database (block mode)
	scratch     *hyperscan.Scratch       // Hyperscan temporary scratch space (used in both modes)
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

func (b *Benchmark) decodePacket(pkt gopacket.Packet) (key *FiveTuple, payload []byte) {
	ipv4, ok := pkt.NetworkLayer().(*layers.IPv4)

	if !ok {
		return // Ignore packets that aren't IPv4
	}

	if ipv4.FragOffset != 0 || (ipv4.Flags&layers.IPv4MoreFragments) != 0 {
		return // Ignore fragmented packets.
	}

	var stream FiveTuple

	stream.protocol = ipv4.Protocol
	stream.srcAddr = ipv4.SrcIP
	stream.dstAddr = ipv4.DstIP

	switch t := pkt.TransportLayer().(type) {
	case *layers.TCP:
		stream.srcPort = uint16(t.SrcPort)
		stream.dstPort = uint16(t.DstPort)
		return &stream, t.Payload

	case *layers.UDP:
		stream.srcPort = uint16(t.SrcPort)
		stream.dstPort = uint16(t.DstPort)
		return &stream, t.Payload
	}

	return
}

// Read a set of streams from a pcap file
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

		key, payload := b.decodePacket(pkt)

		if key == nil || len(payload) == 0 {
			continue
		}

		var id int

		hash := key.Hash()

		if _id, exists := b.streamMap[hash]; exists {
			id = _id
		} else {
			id = len(b.streamMap)
			b.streamMap[hash] = id
		}

		b.packets = append(b.packets, payload)
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

// Return the number of bytes scanned
func (b *Benchmark) Bytes() (sum int) {
	for _, pkt := range b.packets {
		sum += len(pkt)
	}

	return
}

// Return the number of matches found.
func (b *Benchmark) Matches() int { return b.matchCount }

// Clear the number of matches found.
func (b *Benchmark) ClearMatches() { b.matchCount = 0 }

// Display some information about the compiled database and scanned data.
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

// Match event handler: called every time Hyperscan finds a match.
func (b *Benchmark) onMatch(id uint, from, to uint64, flags uint, context interface{}) error {
	b.matchCount += 1

	return nil
}

// Open a Hyperscan stream for each stream in stream_ids
func (b *Benchmark) OpenStreams() error {
	b.streams = make([]hyperscan.Stream, len(b.streamMap))

	var handler hyperscan.MatchHandler = b.onMatch

	for i := 0; i < len(b.streamMap); i++ {
		stream, err := b.dbStreaming.Open(0, b.scratch, handler, nil)

		if err != nil {
			return err
		}

		b.streams[i] = stream
	}

	return nil
}

// Close all open Hyperscan streams (potentially generating any end-anchored matches)
func (b *Benchmark) CloseStreams() error {
	for _, stream := range b.streams {
		if err := stream.Close(); err != nil {
			return err
		}
	}

	return nil
}

// Reset all open Hyperscan streams (potentially generating any end-anchored matches)
func (b *Benchmark) ResetStreams() error {
	for _, stream := range b.streams {
		if err := stream.Reset(); err != nil {
			return err
		}
	}

	return nil
}

// Scan each packet (in the ordering given in the PCAP file) through Hyperscan using the streaming interface.
func (b *Benchmark) ScanStreams() error {
	for i, pkt := range b.packets {
		if err := b.streams[b.streamIds[i]].Scan(pkt); err != nil {
			return err
		}
	}

	return nil
}

// Scan each packet (in the ordering given in the PCAP file) through Hyperscan using the block-mode interface.
func (b *Benchmark) ScanBlock() error {
	var scanner hyperscan.BlockScanner = b.dbBlock
	var handler hyperscan.MatchHandler = b.onMatch

	for _, pkt := range b.packets {
		if err := scanner.Scan(pkt, b.scratch, handler, nil); err != nil {
			return err
		}
	}

	return nil
}

// Simple timing class
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

	if *cpuprofile != "" {
		if f, err := os.Create(*cpuprofile); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Unable to start profiling, %s. Exiting.", err)
			os.Exit(-1)
		} else {
			pprof.StartCPUProfile(f)
			defer pprof.StopCPUProfile()
		}
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

	// Open streams.
	clock.Start()
	if err := bench.OpenStreams(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Unable to open stream, %s. Exiting.", err)
		os.Exit(-1)
	}
	clock.Stop()

	secsStreamingOpenClose += clock.Time()

	// Streaming mode scans.
	for i := 0; i < *repeatCount; i++ {
		if i > 0 {
			clock.Start()
			if err := bench.ResetStreams(); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Unable to reset stream, %s. Exiting.", err)
				os.Exit(-1)
			}
			clock.Stop()
			secsStreamingOpenClose += clock.Time()
		}

		// Scan all our packets in streaming mode.
		clock.Start()
		if err := bench.ScanStreams(); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Unable to scan packet, %s. Exiting.", err)
			os.Exit(-1)
		}
		clock.Stop()

		secsStreamingScan += clock.Time()
	}

	// Close streams.
	clock.Start()
	if err := bench.CloseStreams(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Unable to close stream, %s. Exiting.", err)
		os.Exit(-1)
	}
	clock.Stop()
	secsStreamingOpenClose += clock.Time()

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
		if err := bench.ScanBlock(); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Unable to scan packet, %s. Exiting.", err)
			os.Exit(-1)
		}
	}
	clock.Stop()

	secsScanBlock := clock.Time()

	if *memprofile != "" {
		if f, err := os.Create(*memprofile); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Unable to write memory profile to file, %s. Exiting.", err)
			os.Exit(-1)
		} else {
			pprof.WriteHeapProfile(f)
			f.Close()
		}
	}

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
