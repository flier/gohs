package chimera_test

import (
	"fmt"

	"github.com/flier/gohs/chimera"
)

func ExampleMatch() {
	matched, err := chimera.Match(`foo.*`, []byte(`seafood`))
	fmt.Println(matched, err)
	matched, err = chimera.Match(`bar.*`, []byte(`seafood`))
	fmt.Println(matched, err)
	matched, err = chimera.Match(`a(b`, []byte(`seafood`))
	fmt.Println(matched, err)
	// Output:
	// true <nil>
	// false <nil>
	// false create block database, PCRE compilation failed: missing ).
}
