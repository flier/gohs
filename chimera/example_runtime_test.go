package chimera_test

import (
	"fmt"

	"github.com/flier/gohs/chimera"
)

func ExampleBlockScanner() {
	p, err := chimera.ParsePattern(`foo(bar)+`)
	if err != nil {
		fmt.Println("parse pattern failed,", err)
		return
	}

	// Create new block database with pattern
	db, err := chimera.NewBlockDatabase(p)
	if err != nil {
		fmt.Println("create database failed,", err)
		return
	}
	defer db.Close()

	// Create new scratch for scanning
	s, err := chimera.NewScratch(db)
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

	handler := chimera.HandlerFunc(func(id uint, from, to uint64, flags uint,
		captured []*chimera.Capture, ctx interface{},
	) chimera.Callback {
		matches = append(matches, Match{from, to})
		return chimera.Continue
	})

	data := []byte("hello foobarbar!")

	// Scan data block with handler
	if err := db.Scan(data, s, handler, nil); err != nil {
		fmt.Println("database scan failed,", err)
		return
	}

	// chimera will reports all matches
	for _, m := range matches {
		fmt.Println("match [", m.from, ":", m.to, "]", string(data[m.from:m.to]))
	}

	// Output:
	// match [ 6 : 15 ] foobarbar
}
