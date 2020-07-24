// +build !hyperscan_v4

package hyperscan

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLiteral(t *testing.T) {
	Convey("Give a literal", t, func() {
		Convey("When parse with flags", func() {
			p, err := ParseLiteral(`/test/im`)

			So(err, ShouldBeNil)
			So(p, ShouldNotBeNil)
			So(p.Expression, ShouldEqual, "test")
			So(p.Flags, ShouldEqual, Caseless|MultiLine)

			So(p.Expression.String(), ShouldEqual, "test")
			So(p.String(), ShouldEqual, `/test/im`)

			Convey("When literal contains regular grammer", func() {
				p, err := ParseLiteral(`/te?st/im`)

				So(err, ShouldBeNil)
				So(p, ShouldNotBeNil)
				So(p.Expression, ShouldEqual, "te?st")
				So(p.Flags, ShouldEqual, Caseless|MultiLine)

				So(p.String(), ShouldEqual, "/te?st/im")
			})
		})

		Convey("When parse literal with id and flags", func() {
			p, err := ParseLiteral("3:/foobar/iu")

			So(err, ShouldBeNil)
			So(p.Id, ShouldEqual, 3)
			So(p.Expression, ShouldEqual, "foobar")
			So(p.Flags, ShouldEqual, Caseless|Utf8Mode)
		})

		Convey("When parse with a lot of flags", func() {
			p, err := ParseLiteral(`/test/ismoeupf`)

			So(err, ShouldBeNil)
			So(p, ShouldNotBeNil)
			So(p.Expression, ShouldEqual, "test")
			So(p.Flags, ShouldEqual, Caseless|DotAll|MultiLine|SingleMatch|AllowEmpty|Utf8Mode|UnicodeProperty|PrefilterMode)

			So(p.Flags.String(), ShouldEqual, "8HPVWims")
			So(p.String(), ShouldEqual, "/test/8HPVWims")
		})

		Convey("When parse without flags", func() {
			p, err := ParseLiteral(`test`)

			So(err, ShouldBeNil)
			So(p, ShouldNotBeNil)
			So(p.Expression, ShouldEqual, "test")
			So(p.Flags, ShouldEqual, 0)

			Convey("When literal contains regular grammer", func() {
				p, err := ParseLiteral(`/te?st/im`)

				So(err, ShouldBeNil)
				So(p, ShouldNotBeNil)
				So(p.Expression, ShouldEqual, "te?st")
				So(p.Flags, ShouldEqual, Caseless|MultiLine)

				So(p.String(), ShouldEqual, "/te?st/im")
			})
		})

		Convey("When literal is valid", func() {
			p := NewLiteral("test")

			info, err := p.Info()

			So(err, ShouldBeNil)
			So(info, ShouldNotBeNil)
			So(info, ShouldResemble, &ExprInfo{MinWidth: 4, MaxWidth: 4})
			So(p.IsValid(), ShouldBeTrue)
		})

		Convey("When quote a string", func() {
			So(Quote("test"), ShouldEqual, "`test`")
			So(Quote("`can't backquote this`"), ShouldEqual, "\"`can't backquote this`\"")
		})
	})
}

func TestDatabaseBuilderV5(t *testing.T) {
	Convey("Given a DatabaseBuilder (v5)", t, func() {
		b := DatabaseBuilder{}
		Convey("When build with some combination expression", func() {
			db, err := b.AddExpressions("101:/abc/Q", "102:/def/Q", "/(101&102)/Co").Build()

			So(err, ShouldBeNil)
			So(db, ShouldNotBeNil)

			info, err := db.Info()

			So(err, ShouldBeNil)

			mode, err := info.Mode()

			So(err, ShouldBeNil)
			So(mode, ShouldEqual, BlockMode)

			So(db.Close(), ShouldBeNil)
		})
	})
}
