package hyperscan

import (
	"errors"
	"regexp"
	"testing"
	"unsafe"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	regexInfo = regexp.MustCompile(`^Version: \d\.\d\.\d Features: (NO)?AVX2 Mode: STREAM`)
)

func TestVersion(t *testing.T) {
	Convey("Given a HyperScan version", t, func() {
		ver := hsVersion()

		So(ver, ShouldNotBeEmpty)
		So(ver, ShouldStartWith, "4.")
	})
}

func TestDatabase(t *testing.T) {
	Convey("Given a stream database", t, func() {
		platform, err := hsPopulatePlatform()

		So(platform, ShouldNotBeNil)
		So(platform.info, ShouldNotBeNil)
		So(err, ShouldBeNil)

		db, err := hsCompile("test", 0, Stream, platform)

		So(db, ShouldNotBeNil)
		So(err, ShouldBeNil)

		Convey("Get the database info", func() {
			info, err := hsDatabaseInfo(db)

			So(regexInfo.MatchString(info), ShouldBeTrue)
			So(err, ShouldBeNil)
		})

		Convey("Get the database size", func() {
			size, err := hsDatabaseSize(db)

			So(size, ShouldEqual, 1000)
			So(err, ShouldBeNil)
		})

		Convey("Get the stream size", func() {
			size, err := hsStreamSize(db)

			So(size, ShouldEqual, 24)
			So(err, ShouldBeNil)
		})

		Convey("Get the stream size from a block database", func() {
			db, err := hsCompile("test", 0, Block, platform)

			So(db, ShouldNotBeNil)
			So(err, ShouldBeNil)

			size, err := hsStreamSize(db)

			So(size, ShouldEqual, 0)
			So(err, ShouldEqual, DatabaseModeError)
		})

		Convey("When serialize database", func() {
			data, err := hsSerializeDatabase(db)

			So(data, ShouldNotBeNil)
			So(len(data), ShouldEqual, 1000)
			So(err, ShouldBeNil)

			Convey("Get the database info", func() {
				info, err := hsSerializedDatabaseInfo(data)

				So(regexInfo.MatchString(info), ShouldBeTrue)
				So(err, ShouldBeNil)
			})

			Convey("Get the database size", func() {
				size, err := hsSerializedDatabaseSize(data)

				So(size, ShouldEqual, 1000)
				So(err, ShouldBeNil)
			})

			Convey("Then deserialize database", func() {
				db, err := hsDeserializeDatabase(data)

				So(db, ShouldNotBeNil)
				So(err, ShouldBeNil)

				Convey("Get the database info", func() {
					info, err := hsDatabaseInfo(db)

					So(regexInfo.MatchString(info), ShouldBeTrue)
					So(err, ShouldBeNil)
				})
			})

			Convey("Then deserialize database to memory", func() {
				buf := make([]byte, 1000)
				db := hsDatabase(unsafe.Pointer(&buf[0]))

				So(hsDeserializeDatabaseAt(data, db), ShouldBeNil)

				Convey("Get the database info", func() {
					info, err := hsDatabaseInfo(db)

					So(regexInfo.MatchString(info), ShouldBeTrue)
					So(err, ShouldBeNil)
				})
			})
		})

		So(hsFreeDatabase(db), ShouldBeNil)
	})
}

func TestCompile(t *testing.T) {
	Convey("Given a host platform", t, func() {
		platform, err := hsPopulatePlatform()

		So(platform, ShouldNotBeNil)
		So(platform.info, ShouldNotBeNil)
		So(err, ShouldBeNil)

		Convey("Compile a unsupported expression", func() {
			db, err := hsCompile(`\R`, 0, Stream, platform)

			So(db, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, `\R at index 0 not supported.`)

			So(hsFreeDatabase(db), ShouldBeNil)
		})

		Convey("Compile an empty expression", func() {
			db, err := hsCompile("", 0, Stream, platform)

			So(db, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "Pattern matches empty buffer; use HS_FLAG_ALLOWEMPTY to enable support.")

			So(hsFreeDatabase(db), ShouldBeNil)
		})

		Convey("Compile multi expressions", func() {
			db, err := hsCompileMulti([]string{`^\w+`, `\d+`, `\s+`}, nil, []uint{1, 2, 3}, Stream, platform)

			So(db, ShouldNotBeNil)
			So(err, ShouldBeNil)

			Convey("Get the database info", func() {
				info, err := hsDatabaseInfo(db)

				So(regexInfo.MatchString(info), ShouldBeTrue)
				So(err, ShouldBeNil)
			})

			So(hsFreeDatabase(db), ShouldBeNil)
		})
	})
}

func TestScratch(t *testing.T) {
	Convey("Given a block database", t, func() {
		platform, err := hsPopulatePlatform()

		So(platform, ShouldNotBeNil)
		So(platform.info, ShouldNotBeNil)
		So(err, ShouldBeNil)

		db, err := hsCompile("test", 0, Block, platform)

		So(db, ShouldNotBeNil)
		So(err, ShouldBeNil)

		Convey("Allocate a scratch", func() {
			s, err := hsAllocScratch(db)

			So(s, ShouldNotBeNil)
			So(err, ShouldBeNil)

			Convey("Get the scratch size", func() {
				size, err := hsScratchSize(s)

				So(size, ShouldBeGreaterThan, 1024)
				So(size, ShouldBeLessThan, 4096)
				So(err, ShouldBeNil)

				Convey("Clone the scratch", func() {
					s2, err := hsCloneScratch(s)

					So(s2, ShouldNotBeNil)
					So(err, ShouldBeNil)

					Convey("Cloned scrash should have same size", func() {
						size2, err := hsScratchSize(s2)

						So(size2, ShouldEqual, size)
						So(err, ShouldBeNil)
					})

					So(hsFreeScratch(s2), ShouldBeNil)
				})
			})

			So(hsFreeScratch(s), ShouldBeNil)
		})
	})
}

type matchEvent struct {
	id       uint
	from, to uint64
}

type scanHandler struct {
	matched []matchEvent
	err     error
}

func (h *scanHandler) handle(id uint, from, to uint64, flags uint, context interface{}) error {
	h.matched = append(h.matched, matchEvent{id, from, to})

	return h.err
}

func TestBlockScan(t *testing.T) {
	Convey("Given a block database", t, func() {
		platform, err := hsPopulatePlatform()

		So(platform, ShouldNotBeNil)
		So(platform.info, ShouldNotBeNil)
		So(err, ShouldBeNil)

		db, err := hsCompile("test", 0, Block, platform)

		So(db, ShouldNotBeNil)
		So(err, ShouldBeNil)

		s, err := hsAllocScratch(db)

		So(s, ShouldNotBeNil)
		So(err, ShouldBeNil)

		Convey("Scan block with pattern", func() {
			h := &scanHandler{}

			So(hsScan(db, []byte("abctestdef"), 0, s, h.handle, nil), ShouldBeNil)
			So(h.matched, ShouldResemble, []matchEvent{{0, 0, 7}})
		})

		Convey("Scan block without pattern", func() {
			h := &scanHandler{}

			So(hsScan(db, []byte("abcdef"), 0, s, h.handle, nil), ShouldBeNil)
			So(h.matched, ShouldBeEmpty)
		})

		Convey("Scan block with multi pattern", func() {
			h := &scanHandler{}

			So(hsScan(db, []byte("abctestdeftest"), 0, s, h.handle, nil), ShouldBeNil)
			So(h.matched, ShouldResemble, []matchEvent{{0, 0, 7}, {0, 0, 14}})
		})

		Convey("Scan block with multi pattern but terminated", func() {
			h := &scanHandler{err: errors.New("terminated")}

			So(hsScan(db, []byte("abctestdeftest"), 0, s, h.handle, nil), ShouldEqual, ScanTerminated)
			So(h.matched, ShouldResemble, []matchEvent{{0, 0, 7}})
		})

		So(hsFreeScratch(s), ShouldBeNil)
	})
}

func TestVectorScan(t *testing.T) {
	Convey("Given a block database", t, func() {
		platform, err := hsPopulatePlatform()

		So(platform, ShouldNotBeNil)
		So(platform.info, ShouldNotBeNil)
		So(err, ShouldBeNil)

		db, err := hsCompile("test", 0, Vectored, platform)

		So(db, ShouldNotBeNil)
		So(err, ShouldBeNil)

		s, err := hsAllocScratch(db)

		So(s, ShouldNotBeNil)
		So(err, ShouldBeNil)

		Convey("Scan multi block with pattern", func() {
			h := &scanHandler{}

			So(hsScanVector(db, [][]byte{[]byte("abctestdef"), []byte("abcdef")}, 0, s, h.handle, nil), ShouldBeNil)
			So(h.matched, ShouldResemble, []matchEvent{{0, 0, 7}})
		})

		Convey("Scan multi block without pattern", func() {
			h := &scanHandler{}

			So(hsScanVector(db, [][]byte{[]byte("123456"), []byte("abcdef")}, 0, s, h.handle, nil), ShouldBeNil)
			So(h.matched, ShouldBeEmpty)
		})

		Convey("Scan multi block with multi pattern", func() {
			h := &scanHandler{}

			So(hsScanVector(db, [][]byte{[]byte("abctestdef"), []byte("123test456")}, 0, s, h.handle, nil), ShouldBeNil)
			So(h.matched, ShouldResemble, []matchEvent{{0, 0, 7}, {0, 0, 17}})
		})

		Convey("Scan multi block with multi pattern but terminated", func() {
			h := &scanHandler{err: errors.New("terminated")}

			So(hsScanVector(db, [][]byte{[]byte("abctestdef"), []byte("123test456")}, 0, s, h.handle, nil), ShouldEqual, ScanTerminated)
			So(h.matched, ShouldResemble, []matchEvent{{0, 0, 7}})
		})

		So(hsFreeScratch(s), ShouldBeNil)
	})

}
