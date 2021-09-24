package hyperscan_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/gohs/hyperscan"
)

//nolint:funlen
func TestPattern(t *testing.T) {
	Convey("Give a pattern", t, func() {
		Convey("When parse with flags", func() {
			p, err := hyperscan.ParsePattern(`/test/im`)

			So(err, ShouldBeNil)
			So(p, ShouldNotBeNil)
			So(p.Expression, ShouldEqual, "test")
			So(p.Flags, ShouldEqual, hyperscan.Caseless|hyperscan.MultiLine)

			So(string(p.Expression), ShouldEqual, "test")
			So(p.String(), ShouldEqual, `/test/im`)

			Convey("When pattern contains forward slash", func() {
				p, err := hyperscan.ParsePattern(`/te/st/im`)

				So(err, ShouldBeNil)
				So(p, ShouldNotBeNil)
				So(p.Expression, ShouldEqual, "te/st")
				So(p.Flags, ShouldEqual, hyperscan.Caseless|hyperscan.MultiLine)

				So(p.String(), ShouldEqual, "/te/st/im")
			})
		})

		Convey("When parse pattern with id and flags", func() {
			p, err := hyperscan.ParsePattern("3:/foobar/iu")

			So(err, ShouldBeNil)
			So(p.Id, ShouldEqual, 3)
			So(p.Expression, ShouldEqual, "foobar")
			So(p.Flags, ShouldEqual, hyperscan.Caseless|hyperscan.Utf8Mode)
		})

		Convey("When parse pattern with id, flags and extensions", func() {
			p, err := hyperscan.ParsePattern("3:/foobar/iu{min_offset=4,min_length=8}")
			So(err, ShouldBeNil)
			So(p.Id, ShouldEqual, 3)
			So(p.Expression, ShouldEqual, "foobar")
			So(p.Flags, ShouldEqual, hyperscan.Caseless|hyperscan.Utf8Mode)

			ext, err := p.Ext()
			So(err, ShouldBeNil)
			So(ext, ShouldResemble, new(hyperscan.ExprExt).With(hyperscan.MinOffset(4), hyperscan.MinLength(8)))

			So(p.String(), ShouldEqual, "3:/foobar/8i{min_offset=4,min_length=8}")
		})

		Convey("When parse with a lot of flags", func() {
			p, err := hyperscan.ParsePattern(`/test/ismoeupf`)

			So(err, ShouldBeNil)
			So(p, ShouldNotBeNil)
			So(p.Expression, ShouldEqual, "test")
			So(p.Flags, ShouldEqual, hyperscan.Caseless|hyperscan.DotAll|hyperscan.MultiLine|hyperscan.SingleMatch|
				hyperscan.AllowEmpty|hyperscan.Utf8Mode|hyperscan.UnicodeProperty|hyperscan.PrefilterMode)

			So(p.Flags.String(), ShouldEqual, "8HPVWims")
			So(p.String(), ShouldEqual, "/test/8HPVWims")
		})

		Convey("When parse without flags", func() {
			p, err := hyperscan.ParsePattern(`test`)

			So(err, ShouldBeNil)
			So(p, ShouldNotBeNil)
			So(p.Expression, ShouldEqual, "test")
			So(p.Flags, ShouldEqual, 0)

			Convey("When pattern contains forward slash", func() {
				p, err := hyperscan.ParsePattern(`te/st`)

				So(err, ShouldBeNil)
				So(p, ShouldNotBeNil)
				So(p.Expression, ShouldEqual, "te/st")
				So(p.Flags, ShouldEqual, 0)
			})
		})

		Convey("When pattern is valid", func() {
			p := hyperscan.Pattern{Expression: "test"}

			info, err := p.Info()

			So(err, ShouldBeNil)
			So(info, ShouldNotBeNil)
			So(info, ShouldResemble, &hyperscan.ExprInfo{MinWidth: 4, MaxWidth: 4})
			So(p.IsValid(), ShouldBeTrue)
		})

		Convey("When pattern is invalid", func() {
			p := hyperscan.Pattern{Expression: `\R`}

			info, err := p.Info()

			So(err, ShouldNotBeNil)
			So(info, ShouldBeNil)
			So(p.IsValid(), ShouldBeFalse)
		})

		Convey("When quote a string", func() {
			So(hyperscan.Quote("test"), ShouldEqual, "`test`")
			So(hyperscan.Quote("`can't backquote this`"), ShouldEqual, "\"`can't backquote this`\"")
		})
	})
}
