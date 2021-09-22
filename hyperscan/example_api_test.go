package hyperscan_test

import (
	"fmt"
	"strings"

	"github.com/flier/gohs/hyperscan"
)

func ExampleMatch() {
	matched, err := hyperscan.Match(`foo.*`, []byte(`seafood`))
	fmt.Println(matched, err)
	matched, err = hyperscan.Match(`bar.*`, []byte(`seafood`))
	fmt.Println(matched, err)
	matched, err = hyperscan.Match(`a(b`, []byte(`seafood`))
	fmt.Println(matched, err)
	// Output:
	// true <nil>
	// false <nil>
	// false parse pattern, invalid pattern `a(b`, Missing close parenthesis for group started at index 1.
}

func ExampleMatchReader() {
	s := strings.NewReader(strings.Repeat("a", 4096) + `seafood`)
	matched, err := hyperscan.MatchReader(`foo.*`, s)
	fmt.Println(matched, err)
	matched, err = hyperscan.MatchReader(`bar.*`, s)
	fmt.Println(matched, err)
	matched, err = hyperscan.MatchReader(`a(b`, s)
	fmt.Println(matched, err)
	// Output:
	// true <nil>
	// false <nil>
	// false parse pattern, invalid pattern `a(b`, Missing close parenthesis for group started at index 1.
}
