package hs_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/gohs/internal/hs"
)

//nolint:funlen
func TestStreamScan(t *testing.T) {
	Convey("Given a stream database", t, func() {
		platform, err := hs.PopulatePlatform()

		So(platform, ShouldNotBeNil)
		So(err, ShouldBeNil)

		db, err := hs.Compile("test", 0, hs.StreamMode, platform)

		So(db, ShouldNotBeNil)
		So(err, ShouldBeNil)

		s, err := hs.AllocScratch(db)

		So(s, ShouldNotBeNil)
		So(err, ShouldBeNil)

		Convey("Then open a stream", func() {
			stream, err := hs.OpenStream(db, 0)

			So(stream, ShouldNotBeNil)
			So(err, ShouldBeNil)

			h := &hs.MatchRecorder{}

			Convey("Then scan a simple stream with first part", func() {
				So(hs.ScanStream(stream, []byte("abcte"), 0, s, h.Handle, nil), ShouldBeNil)
				So(h.Events, ShouldBeNil)

				Convey("When scan second part, should be matched", func() {
					So(hs.ScanStream(stream, []byte("stdef"), 0, s, h.Handle, nil), ShouldBeNil)
					So(h.Events, ShouldResemble, []hs.MatchEvent{{0, 0, 7, 0}})
				})

				Convey("Then copy the stream", func() {
					stream2, err := hs.CopyStream(stream)

					So(stream2, ShouldNotBeNil)
					So(err, ShouldBeNil)

					Convey("When copied stream2 scan the second part, should be matched", func() {
						So(hs.ScanStream(stream2, []byte("stdef"), 0, s, h.Handle, nil), ShouldBeNil)
						So(h.Events, ShouldResemble, []hs.MatchEvent{{0, 0, 7, 0}})

						Convey("When copied stream2 scan the second part again, should not be matched", func() {
							h.Events = nil
							So(hs.ScanStream(stream2, []byte("stdef"), 0, s, h.Handle, nil), ShouldBeNil)
							So(h.Events, ShouldBeNil)

							Convey("When copy and reset stream2", func() {
								So(hs.ResetAndCopyStream(stream2, stream, s, h.Handle, nil), ShouldBeNil)

								Convey("When copied and reset stream2 scan the second part again, should be matched", func() {
									h.Events = nil
									So(hs.ScanStream(stream2, []byte("stdef"), 0, s, h.Handle, nil), ShouldBeNil)
									So(h.Events, ShouldResemble, []hs.MatchEvent{{0, 0, 7, 0}})
								})
							})
						})
					})

					So(hs.CloseStream(stream2, s, h.Handle, nil), ShouldBeNil)
				})

				Convey("Then reset the stream", func() {
					So(hs.ResetStream(stream, 0, s, h.Handle, nil), ShouldBeNil)

					Convey("When scan the second part, should not be matched", func() {
						So(hs.ScanStream(stream, []byte("stdef"), 0, s, h.Handle, nil), ShouldBeNil)
						So(h.Events, ShouldBeNil)
					})

					Convey("When scan empty buffers", func() {
						So(hs.ScanStream(stream, nil, 0, s, h.Handle, nil), ShouldEqual, hs.ErrInvalid)
						So(hs.ScanStream(stream, []byte(""), 0, s, h.Handle, nil), ShouldBeNil)
					})
				})
			})

			So(hs.CloseStream(stream, s, h.Handle, nil), ShouldBeNil)
		})

		So(hs.FreeScratch(s), ShouldBeNil)
	})
}
