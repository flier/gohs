package hs_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/gohs/hyperscan"
	"github.com/flier/gohs/internal/hs"
)

//nolint:funlen
func TestScratch(t *testing.T) {
	Convey("Given a block database", t, func() {
		platform, err := hs.PopulatePlatform()

		So(platform, ShouldNotBeNil)
		So(err, ShouldBeNil)

		db, err := hs.Compile("test", 0, hs.BlockMode, platform)

		So(db, ShouldNotBeNil)
		So(err, ShouldBeNil)

		Convey("Allocate a scratch", func() {
			s, err := hs.AllocScratch(db)

			So(s, ShouldNotBeNil)
			So(err, ShouldBeNil)

			Convey("Get the scratch size", func() {
				size, err := hs.ScratchSize(s)

				So(size, ShouldBeGreaterThan, 1024)
				So(size, ShouldBeLessThan, 4096)
				So(err, ShouldBeNil)

				Convey("Clone the scratch", func() {
					s2, err := hs.CloneScratch(s)

					So(s2, ShouldNotBeNil)
					So(err, ShouldBeNil)

					Convey("Cloned scrash should have same size", func() {
						size2, err := hs.ScratchSize(s2)

						So(size2, ShouldEqual, size)
						So(err, ShouldBeNil)
					})

					So(hs.FreeScratch(s2), ShouldBeNil)
				})

				Convey("Reallocate the scratch with another database", func() {
					db2, err := hs.Compile(hyperscan.EmailAddress, 0, hyperscan.BlockMode, platform)

					So(db, ShouldNotBeNil)
					So(err, ShouldBeNil)

					So(hs.ReallocScratch(db2, &s), ShouldBeNil)

					size2, err := hs.ScratchSize(s)

					So(size2, ShouldBeGreaterThan, size)
					So(err, ShouldBeNil)

					So(hs.FreeDatabase(db2), ShouldBeNil)
				})
			})

			So(hs.FreeScratch(s), ShouldBeNil)
		})

		So(hs.FreeDatabase(db), ShouldBeNil)
	})
}
