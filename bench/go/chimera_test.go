package scan_test

import (
	"strings"
	"testing"

	"github.com/flier/gohs/chimera"
)

func BenchmarkChimeraBlockScan(b *testing.B) {
	isRaceBuilder := strings.HasSuffix(testenv(), "-race")

	for _, data := range benchData {
		p := chimera.NewPattern(data.re, chimera.MultiLine)
		db, err := chimera.NewBlockDatabase(p)
		if err != nil {
			b.Fatalf("compile pattern %s: `%s`, %s", data.name, data.re, err)
		}

		s, err := chimera.NewScratch(db)
		if err != nil {
			b.Fatalf("create scratch, %s", err)
		}

		m := chimera.HandlerFunc(func(id uint, from, to uint64, flags uint,
			captured []*chimera.Capture, context interface{},
		) chimera.Callback {
			return chimera.Terminate
		})

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
