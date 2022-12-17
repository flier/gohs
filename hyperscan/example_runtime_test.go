package hyperscan_test

import (
	"fmt"

	"github.com/flier/gohs/hyperscan"
)

func ExampleBlockScanner() {
	// Pattern with `L` flag enable leftmost start of match reporting.
	p, err := hyperscan.ParsePattern(`/foo(bar)+/L`)
	if err != nil {
		fmt.Println("parse pattern failed,", err)
		return
	}

	// Create new block database with pattern
	db, err := hyperscan.NewBlockDatabase(p)
	if err != nil {
		fmt.Println("create database failed,", err)
		return
	}
	defer db.Close()

	// Create new scratch for scanning
	s, err := hyperscan.NewScratch(db)
	if err != nil {
		fmt.Println("create scratch failed,", err)
		return
	}

	defer func() {
		_ = s.Free()
	}()

	// Record matching text
	type Match struct {
		from uint64
		to   uint64
	}

	var matches []Match

	handler := hyperscan.MatchHandler(func(id uint, from, to uint64, flags uint, context interface{}) error {
		matches = append(matches, Match{from, to})
		return nil
	})

	data := []byte("hello foobarbar!")

	// Scan data block with handler
	if err := db.Scan(data, s, handler, nil); err != nil {
		fmt.Println("database scan failed,", err)
		return
	}

	// Hyperscan will reports all matches
	for _, m := range matches {
		fmt.Println("match [", m.from, ":", m.to, "]", string(data[m.from:m.to]))
	}

	// Output:
	// match [ 6 : 12 ] foobar
	// match [ 6 : 15 ] foobarbar
}

func ExampleVectoredScanner() {
	// Pattern with `L` flag enable leftmost start of match reporting.
	p, err := hyperscan.ParsePattern(`/foo(bar)+/L`)
	if err != nil {
		fmt.Println("parse pattern failed,", err)
		return
	}

	// Create new vectored database with pattern
	db, err := hyperscan.NewVectoredDatabase(p)
	if err != nil {
		fmt.Println("create database failed,", err)
		return
	}
	defer db.Close()

	// Create new scratch for scanning
	s, err := hyperscan.NewScratch(db)
	if err != nil {
		fmt.Println("create scratch failed,", err)
		return
	}

	defer func() {
		_ = s.Free()
	}()

	// Record matching text
	type Match struct {
		from uint64
		to   uint64
	}

	var matches []Match

	handler := hyperscan.MatchHandler(func(id uint, from, to uint64, flags uint, context interface{}) error {
		matches = append(matches, Match{from, to})
		return nil
	})

	data := []byte("hello foobarbar!")

	// Scan vectored data with handler
	if err := db.Scan([][]byte{data[:8], data[8:12], data[12:]}, s, handler, nil); err != nil {
		fmt.Println("database scan failed,", err)
		return
	}

	// Hyperscan will reports all matches
	for _, m := range matches {
		fmt.Println("match [", m.from, ":", m.to, "]", string(data[m.from:m.to]))
	}

	// Output:
	// match [ 6 : 12 ] foobar
	// match [ 6 : 15 ] foobarbar
}

func ExampleStreamScanner() { //nolint:funlen
	// Pattern with `L` flag enable leftmost start of match reporting.
	p, err := hyperscan.ParsePattern(`/foo(bar)+/L`)
	if err != nil {
		fmt.Println("parse pattern failed,", err)
		return
	}

	// Create new stream database with pattern
	db, err := hyperscan.NewStreamDatabase(p)
	if err != nil {
		fmt.Println("create database failed,", err)
		return
	}
	defer db.Close()

	// Create new scratch for scanning
	s, err := hyperscan.NewScratch(db)
	if err != nil {
		fmt.Println("create scratch failed,", err)
		return
	}

	defer func() {
		_ = s.Free()
	}()

	// Record matching text
	type Match struct {
		from uint64
		to   uint64
	}

	var matches []Match

	handler := hyperscan.MatchHandler(func(id uint, from, to uint64, flags uint, context interface{}) error {
		matches = append(matches, Match{from, to})
		return nil
	})

	data := []byte("hello foobarbar!")

	// Open stream with handler
	st, err := db.Open(0, s, handler, nil)
	if err != nil {
		fmt.Println("open streaming database failed,", err)
		return
	}

	// Scan data with stream
	for i := 0; i < len(data); i += 4 {
		start := i
		end := i + 4

		if end > len(data) {
			end = len(data)
		}

		if err = st.Scan(data[start:end]); err != nil {
			fmt.Println("streaming scan failed,", err)
			return
		}
	}

	// Close stream
	if err = st.Close(); err != nil {
		fmt.Println("streaming scan failed,", err)
		return
	}

	// Hyperscan will reports all matches
	for _, m := range matches {
		fmt.Println("match [", m.from, ":", m.to, "]", string(data[m.from:m.to]))
	}

	// Output:
	// match [ 6 : 12 ] foobar
	// match [ 6 : 15 ] foobarbar
}
