package chimera_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/gohs/chimera"
)

//nolint:funlen
func TestPattern(t *testing.T) {
	Convey("Give a pattern", t, func() {
		Convey("When parse with flags", func() {
			p, err := chimera.ParsePattern(`/test/im`)

			So(err, ShouldBeNil)
			So(p, ShouldNotBeNil)
			So(p.Expression, ShouldEqual, "test")
			So(p.Flags, ShouldEqual, chimera.Caseless|chimera.MultiLine)

			So(p.Expression, ShouldEqual, "test")
			So(p.String(), ShouldEqual, `/test/im`)

			Convey("When pattern contains forward slash", func() {
				p, err := chimera.ParsePattern(`/te/st/im`)

				So(err, ShouldBeNil)
				So(p, ShouldNotBeNil)
				So(p.Expression, ShouldEqual, "te/st")
				So(p.Flags, ShouldEqual, chimera.Caseless|chimera.MultiLine)

				So(p.String(), ShouldEqual, "/te/st/im")
			})
		})

		Convey("When parse pattern with id and flags", func() {
			p, err := chimera.ParsePattern("3:/foobar/i8")

			So(err, ShouldBeNil)
			So(p.ID, ShouldEqual, 3)
			So(p.Expression, ShouldEqual, "foobar")
			So(p.Flags, ShouldEqual, chimera.Caseless|chimera.Utf8Mode)
		})

		Convey("When parse with a lot of flags", func() {
			p, err := chimera.ParsePattern(`/test/ismH8W`)

			So(err, ShouldBeNil)
			So(p, ShouldNotBeNil)
			So(p.Expression, ShouldEqual, "test")
			So(p.Flags, ShouldEqual, chimera.Caseless|chimera.DotAll|chimera.MultiLine|chimera.SingleMatch|
				chimera.Utf8Mode|chimera.UnicodeProperty)

			So(p.Flags.String(), ShouldEqual, "8HWims")
			So(p.String(), ShouldEqual, "/test/8HWims")
		})

		Convey("When parse without flags", func() {
			p, err := chimera.ParsePattern(`test`)

			So(err, ShouldBeNil)
			So(p, ShouldNotBeNil)
			So(p.Expression, ShouldEqual, "test")
			So(p.Flags, ShouldEqual, 0)

			Convey("When pattern contains forward slash", func() {
				p, err := chimera.ParsePattern(`te/st`)

				So(err, ShouldBeNil)
				So(p, ShouldNotBeNil)
				So(p.Expression, ShouldEqual, "te/st")
				So(p.Flags, ShouldEqual, 0)
			})
		})

		Convey("When quote a string", func() {
			So(chimera.Quote("test"), ShouldEqual, "`test`")
			So(chimera.Quote("`can't backquote this`"), ShouldEqual, "\"`can't backquote this`\"")
		})
	})
}
