package hs_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/gohs/internal/hs"
)

func TestBlockScan(t *testing.T) {
	Convey("Given a block database", t, func() {
		platform, err := hs.PopulatePlatform()

		So(platform, ShouldNotBeNil)
		So(err, ShouldBeNil)

		db, err := hs.Compile("test", 0, hs.BlockMode, platform)

		So(db, ShouldNotBeNil)
		So(err, ShouldBeNil)

		s, err := hs.AllocScratch(db)

		So(s, ShouldNotBeNil)
		So(err, ShouldBeNil)

		h := &hs.MatchRecorder{}

		Convey("Scan block with pattern", func() {
			So(hs.Scan(db, []byte("abctestdef"), 0, s, h.Handle, nil), ShouldBeNil)
			So(h.Events, ShouldResemble, []hs.MatchEvent{{0, 0, 7, 0}})
		})

		Convey("Scan block without pattern", func() {
			So(hs.Scan(db, []byte("abcdef"), 0, s, h.Handle, nil), ShouldBeNil)
			So(h.Events, ShouldBeEmpty)
		})

		Convey("Scan block with multi pattern", func() {
			So(hs.Scan(db, []byte("abctestdeftest"), 0, s, h.Handle, nil), ShouldBeNil)
			So(h.Events, ShouldResemble, []hs.MatchEvent{{0, 0, 14, 0}})
		})

		Convey("Scan block with multi pattern but terminated", func() {
			h.Err = hs.ErrScanTerminated

			So(hs.Scan(db, []byte("abctestdeftest"), 0, s, h.Handle, nil), ShouldEqual, hs.ErrScanTerminated)
			So(h.Events, ShouldResemble, []hs.MatchEvent{{0, 0, 7, 0}})
		})

		Convey("Scan empty buffers", func() {
			So(hs.Scan(db, nil, 0, s, h.Handle, nil), ShouldEqual, hs.ErrInvalid)
			So(hs.Scan(db, []byte(""), 0, s, h.Handle, nil), ShouldBeNil)
		})

		So(hs.FreeScratch(s), ShouldBeNil)
		So(hs.FreeDatabase(db), ShouldBeNil)
	})
}

func TestVectorScan(t *testing.T) {
	Convey("Given a block database", t, func() {
		platform, err := hs.PopulatePlatform()

		So(platform, ShouldNotBeNil)
		So(err, ShouldBeNil)

		db, err := hs.Compile("test", 0, hs.VectoredMode, platform)

		So(db, ShouldNotBeNil)
		So(err, ShouldBeNil)

		s, err := hs.AllocScratch(db)

		So(s, ShouldNotBeNil)
		So(err, ShouldBeNil)

		h := &hs.MatchRecorder{}

		Convey("Scan multi block with pattern", func() {
			So(hs.ScanVector(db, [][]byte{[]byte("abctestdef"), []byte("abcdef")}, 0, s, h.Handle, nil), ShouldBeNil)
			So(h.Events, ShouldResemble, []hs.MatchEvent{{0, 0, 7, 0}})
		})

		Convey("Scan multi block without pattern", func() {
			So(hs.ScanVector(db, [][]byte{[]byte("123456"), []byte("abcdef")}, 0, s, h.Handle, nil), ShouldBeNil)
			So(h.Events, ShouldBeEmpty)
		})

		Convey("Scan multi block with multi pattern", func() {
			So(hs.ScanVector(db, [][]byte{[]byte("abctestdef"), []byte("123test456")}, 0, s, h.Handle, nil), ShouldBeNil)
			So(h.Events, ShouldResemble, []hs.MatchEvent{{0, 0, 17, 0}})
		})

		Convey("Scan multi block with multi pattern but terminated", func() {
			h.Err = hs.ErrScanTerminated

			So(hs.ScanVector(db, [][]byte{[]byte("abctestdef"), []byte("123test456")}, 0, s, h.Handle, nil),
				ShouldEqual, hs.ErrScanTerminated)
			So(h.Events, ShouldResemble, []hs.MatchEvent{{0, 0, 7, 0}})
		})

		Convey("Scan empty buffers", func() {
			So(hs.ScanVector(db, nil, 0, s, h.Handle, nil), ShouldEqual, hs.ErrInvalid)
			So(hs.ScanVector(db, [][]byte{}, 0, s, h.Handle, nil), ShouldBeNil)
			So(hs.ScanVector(db, [][]byte{[]byte(""), []byte("")}, 0, s, h.Handle, nil), ShouldBeNil)
		})

		So(hs.FreeScratch(s), ShouldBeNil)
	})
}
