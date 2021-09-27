package chimera_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/gohs/chimera"
)

func TestChimera(t *testing.T) {
	Convey("Given a chimera runtimes", t, func() {
		So(chimera.Version(), ShouldNotBeEmpty)
	})
}

func TestBaseDatabase(t *testing.T) {
	Convey("Given a block database", t, func() {
		So(chimera.ValidPlatform(), ShouldBeNil)

		bdb, err := chimera.NewBlockDatabase(&chimera.Pattern{Expression: "test"})

		So(err, ShouldBeNil)
		So(bdb, ShouldNotBeNil)

		Convey("When get size", func() {
			size, err := bdb.Size()

			So(err, ShouldBeNil)
			So(size, ShouldBeGreaterThan, 800)
		})

		Convey("When get info", func() {
			info, err := bdb.Info()

			So(err, ShouldBeNil)
			So(info, ShouldNotBeNil)

			_, _, _, err = info.Parse()
			So(err, ShouldBeNil)

			Convey("Then get version", func() {
				ver, err := info.Version()

				So(err, ShouldBeNil)
				So(chimera.Version(), ShouldStartWith, ver)
			})

			Convey("Then get mode", func() {
				mode, err := info.Mode()

				So(err, ShouldBeNil)
				So(mode, ShouldEqual, chimera.BlockMode)
			})
		})

		So(bdb.Close(), ShouldBeNil)
	})
}

func TestBlockDatabase(t *testing.T) {
	Convey("Give a block database", t, func() {
		bdb, err := chimera.NewBlockDatabase(&chimera.Pattern{Expression: "test"})

		So(err, ShouldBeNil)
		So(bdb, ShouldNotBeNil)

		Convey("When get info", func() {
			info, err := bdb.Info()

			So(err, ShouldBeNil)
			So(info, ShouldNotBeNil)

			_, _, _, err = info.Parse()
			So(err, ShouldBeNil)

			Convey("Then get mode", func() {
				mode, err := info.Mode()

				So(err, ShouldBeNil)
				So(mode, ShouldEqual, chimera.BlockMode)
			})
		})

		So(bdb.Close(), ShouldBeNil)
	})
}
