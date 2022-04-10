//go:build chimera
// +build chimera

package chimera_test

import (
	"fmt"
	"strings"

	"github.com/flier/gohs/chimera"
)

// This example demonstrates construct and match a pattern.
func ExamplePattern() {
	p := chimera.NewPattern(`foo.*bar`, chimera.Caseless)
	fmt.Println(p)

	db, err := chimera.NewBlockDatabase(p)
	fmt.Println(err)

	found := db.MatchString("fooxyzbarbar")
	fmt.Println(found)

	// Output:
	// /foo.*bar/i
	// <nil>
	// true
}

// This example demonstrates parsing pattern with id and flags.
func ExampleParsePattern() {
	p, err := chimera.ParsePattern("3:/foobar/i8")

	fmt.Println(err)
	fmt.Println(p.Id)
	fmt.Println(p.Expression)
	fmt.Println(p.Flags)

	// Output:
	// <nil>
	// 3
	// foobar
	// 8i
}

// This example demonstrates parsing patterns with comment.
func ExampleParsePatterns() {
	patterns, err := chimera.ParsePatterns(strings.NewReader(`
# empty line and comment will be skipped

1:/hatstand.*teakettle/s
2:/(hatstand|teakettle)/iH
3:/^.{10,20}hatstand/m
`))

	fmt.Println(err)

	for _, p := range patterns {
		fmt.Println(p)
	}

	// Output:
	// <nil>
	// 1:/hatstand.*teakettle/s
	// 2:/(hatstand|teakettle)/Hi
	// 3:/^.{10,20}hatstand/m
}
