//go:build chimera
// +build chimera

package ch_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/gohs/internal/ch"
	"github.com/flier/gohs/internal/hs"
)

func TestBlockScan(t *testing.T) {
	Convey("Given a block database", t, func() {
		platform, err := hs.PopulatePlatform()

		So(platform, ShouldNotBeNil)
		So(err, ShouldBeNil)

		db, err := ch.Compile("test", 0, ch.Groups, platform)

		So(db, ShouldNotBeNil)
		So(err, ShouldBeNil)

		s, err := ch.AllocScratch(db)

		So(s, ShouldNotBeNil)
		So(err, ShouldBeNil)

		h := &ch.MatchRecorder{}

		Convey("Scan block with pattern", func() {
			So(ch.Scan(db, []byte("abctestdef"), 0, s, h.OnMatch, h.OnError, nil), ShouldBeNil)
			So(h.Events, ShouldResemble, []ch.MatchEvent{{0, 3, 7, 0, []*ch.Capture{{3, 7, []byte("test")}}}})
		})

		Convey("Scan block without pattern", func() {
			So(ch.Scan(db, []byte("abcdef"), 0, s, h.OnMatch, h.OnError, nil), ShouldBeNil)
			So(h.Events, ShouldBeEmpty)
		})

		Convey("Scan block with multi pattern", func() {
			So(ch.Scan(db, []byte("abctestdeftest"), 0, s, h.OnMatch, h.OnError, nil), ShouldBeNil)
			So(h.Events, ShouldResemble, []ch.MatchEvent{
				{0, 3, 7, 0, []*ch.Capture{{3, 7, []byte("test")}}},
				{0, 10, 14, 0, []*ch.Capture{{10, 14, []byte("test")}}},
			})
		})

		Convey("Scan block with multi pattern but terminated", func() {
			onMatch := func(id uint, from, to uint64, flags uint, captured []*ch.Capture, context interface{}) ch.Callback {
				return ch.Terminate
			}

			So(ch.Scan(db, []byte("abctestdeftest"), 0, s, onMatch, h.OnError, nil), ShouldEqual, ch.ErrScanTerminated)
		})

		Convey("Scan empty buffers", func() {
			So(ch.Scan(db, nil, 0, s, h.OnMatch, h.OnError, nil), ShouldEqual, ch.ErrInvalid)
			So(ch.Scan(db, []byte(""), 0, s, h.OnMatch, h.OnError, nil), ShouldBeNil)
		})

		So(ch.FreeScratch(s), ShouldBeNil)
		So(ch.FreeDatabase(db), ShouldBeNil)
	})
}
