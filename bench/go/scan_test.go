package scan_test

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/flier/gohs/hyperscan"
)

var benchData = []struct{ name, re string }{
	{"Easy0", "ABCDEFGHIJKLMNOPQRSTUVWXYZ$"},
	{"Easy0i", "(?i)ABCDEFGHIJklmnopqrstuvwxyz$"},
	{"Easy1", "A[AB]B[BC]C[CD]D[DE]E[EF]F[FG]G[GH]H[HI]I[IJ]J$"},
	{"Medium", "[XYZ]ABCDEFGHIJKLMNOPQRSTUVWXYZ$"},
	{"Hard", "[ -~]*ABCDEFGHIJKLMNOPQRSTUVWXYZ$"},
	{"Hard1", "ABCD|CDEF|EFGH|GHIJ|IJKL|KLMN|MNOP|OPQR|QRST|STUV|UVWX|WXYZ"},
}

var benchSizes = []struct {
	name string
	n    int
}{
	{"16", 16},
	{"32", 32},
	{"1K", 1 << 10},
	{"32K", 32 << 10},
	{"1M", 1 << 20},
	{"32M", 32 << 20},
}

func testenv() string {
	return os.Getenv("GO_BUILDER_NAME")
}

var text []byte

func makeText(n int) []byte {
	if len(text) >= n {
		return text[:n]
	}
	text = make([]byte, n)
	x := ^uint32(0)
	for i := range text {
		x += x
		x ^= 1
		if int32(x) < 0 {
			x ^= 0x88888eef
		}
		if x%31 == 0 {
			text[i] = '\n'
		} else {
			text[i] = byte(x%(0x7E+1-0x20) + 0x20)
		}
	}
	return text
}

func BenchmarkHyperscanBlockScan(b *testing.B) {
	isRaceBuilder := strings.HasSuffix(testenv(), "-race")

	for _, data := range benchData {
		p := hyperscan.NewPattern(data.re, hyperscan.MultiLine)
		db, err := hyperscan.NewBlockDatabase(p)
		if err != nil {
			b.Fatalf("compile pattern %s: `%s`, %s", data.name, data.re, err)
		}

		s, err := hyperscan.NewScratch(db)
		if err != nil {
			b.Fatalf("create scratch, %s", err)
		}

		m := func(id uint, from, to uint64, flags uint, context interface{}) error {
			return hyperscan.ErrScanTerminated
		}

		for _, size := range benchSizes {
			if (isRaceBuilder || testing.Short()) && size.n > 1<<10 {
				continue
			}
			t := makeText(size.n)
			b.Run(data.name+"/"+size.name, func(b *testing.B) {
				b.SetBytes(int64(len(t)))
				for i := 0; i < b.N; i++ {
					if err = db.Scan(t, s, m, nil); err != nil {
						b.Fatalf("match, %s", err)
					}
				}
			})
		}
	}
}

const PageSize = 4096

func BenchmarkHyperscanStreamScan(b *testing.B) { //nolint: gocognit
	isRaceBuilder := strings.HasSuffix(testenv(), "-race")

	for _, data := range benchData {
		p := hyperscan.NewPattern(data.re, hyperscan.MultiLine)
		db, err := hyperscan.NewStreamDatabase(p)
		if err != nil {
			b.Fatalf("compile pattern %s: `%s`, %s", data.name, data.re, err)
		}

		s, err := hyperscan.NewScratch(db)
		if err != nil {
			b.Fatalf("create scratch, %s", err)
		}

		m := func(id uint, from, to uint64, flags uint, context interface{}) error {
			return hyperscan.ErrScanTerminated
		}

		for _, size := range benchSizes {
			if (isRaceBuilder || testing.Short()) && size.n > 1<<10 {
				continue
			}
			t := makeText(size.n)
			b.Run(data.name+"/"+size.name, func(b *testing.B) {
				b.SetBytes(int64(len(t)))
				for i := 0; i < b.N; i++ {
					st, err := db.Open(0, s, m, nil)
					if err != nil {
						b.Fatalf("open stream, %s", err)
					}
					for i := 0; i < len(t); i += PageSize {
						n := len(t) - i
						if n > PageSize {
							n = PageSize
						}
						if err = st.Scan(t[i : i+n]); err != nil {
							b.Fatalf("match, %s", err)
						}
					}
					st.Close()
				}
			})
		}
	}
}

func BenchmarkRegexpMatch(b *testing.B) {
	isRaceBuilder := strings.HasSuffix(testenv(), "-race")

	for _, data := range benchData {
		r := regexp.MustCompile(data.re)
		for _, size := range benchSizes {
			if (isRaceBuilder || testing.Short()) && size.n > 1<<10 {
				continue
			}
			t := makeText(size.n)
			b.Run(data.name+"/"+size.name, func(b *testing.B) {
				b.SetBytes(int64(len(t)))
				for i := 0; i < b.N; i++ {
					if r.Match(t) {
						b.Fatal("match!")
					}
				}
			})
		}
	}
}
