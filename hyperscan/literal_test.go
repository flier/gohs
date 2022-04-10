//go:build !hyperscan_v4
// +build !hyperscan_v4

package hyperscan_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/gohs/hyperscan"
)

//nolint:funlen
func TestLiteral(t *testing.T) {
	Convey("Give a literal", t, func() {
		Convey("When parse with flags", func() {
			p, err := hyperscan.ParseLiteral(`/test/im`)

			So(err, ShouldBeNil)
			So(p, ShouldNotBeNil)
			So(p.Expression, ShouldEqual, "test")
			So(p.Flags, ShouldEqual, hyperscan.Caseless|hyperscan.MultiLine)

			So(p.Expression, ShouldEqual, "test")
			So(p.String(), ShouldEqual, `/test/im`)

			Convey("When literal contains regular grammar", func() {
				p, err := hyperscan.ParseLiteral(`/te?st/im`)

				So(err, ShouldBeNil)
				So(p, ShouldNotBeNil)
				So(p.Expression, ShouldEqual, "te?st")
				So(p.Flags, ShouldEqual, hyperscan.Caseless|hyperscan.MultiLine)

				So(p.String(), ShouldEqual, "/te?st/im")
			})
		})

		Convey("When parse literal with id and flags", func() {
			p, err := hyperscan.ParseLiteral("3:/foobar/iu")

			So(err, ShouldBeNil)
			So(p.Id, ShouldEqual, 3)
			So(p.Expression, ShouldEqual, "foobar")
			So(p.Flags, ShouldEqual, hyperscan.Caseless|hyperscan.Utf8Mode)
		})

		Convey("When parse with a lot of flags", func() {
			p, err := hyperscan.ParseLiteral(`/test/ismoeupf`)

			So(err, ShouldBeNil)
			So(p, ShouldNotBeNil)
			So(p.Expression, ShouldEqual, "test")
			So(p.Flags, ShouldEqual, hyperscan.Caseless|hyperscan.DotAll|hyperscan.MultiLine|hyperscan.SingleMatch|
				hyperscan.AllowEmpty|hyperscan.Utf8Mode|hyperscan.UnicodeProperty|hyperscan.PrefilterMode)

			So(p.Flags.String(), ShouldEqual, "8HPVWims")
			So(p.String(), ShouldEqual, "/test/8HPVWims")
		})

		Convey("When parse without flags", func() {
			p, err := hyperscan.ParseLiteral(`test`)

			So(err, ShouldBeNil)
			So(p, ShouldNotBeNil)
			So(p.Expression, ShouldEqual, "test")
			So(p.Flags, ShouldEqual, 0)

			Convey("When literal contains regular grammar", func() {
				p, err := hyperscan.ParseLiteral(`/te?st/im`)

				So(err, ShouldBeNil)
				So(p, ShouldNotBeNil)
				So(p.Expression, ShouldEqual, "te?st")
				So(p.Flags, ShouldEqual, hyperscan.Caseless|hyperscan.MultiLine)

				So(p.String(), ShouldEqual, "/te?st/im")
			})
		})

		Convey("When literal is valid", func() {
			p := hyperscan.NewLiteral("test")

			info, err := p.Info()

			So(err, ShouldBeNil)
			So(info, ShouldNotBeNil)
			So(info, ShouldResemble, &hyperscan.ExprInfo{MinWidth: 4, MaxWidth: 4})
			So(p.IsValid(), ShouldBeTrue)
		})

		Convey("When quote a string", func() {
			So(hyperscan.Quote("test"), ShouldEqual, "`test`")
			So(hyperscan.Quote("`can't backquote this`"), ShouldEqual, "\"`can't backquote this`\"")
		})
	})
}
