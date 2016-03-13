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

	"github.com/flier/gohs/hyperscan"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	flagRepeatCount = flag.Int("n", 1, "Repeating PCAP scan several times")
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

	sdb, err := hyperscan.NewStreamDatabase(patterns...)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Could not compile patterns, %s", err)
		os.Exit(-1)
	}

	bdb, err := hyperscan.NewBlockDatabase(patterns...)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Could not compile patterns, %s", err)
		os.Exit(-1)
	}

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
}
