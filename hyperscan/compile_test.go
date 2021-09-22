//nolint:funlen
package hyperscan_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/gohs/hyperscan"
)

func TestPattern(t *testing.T) {
	Convey("Give a pattern", t, func() {
		Convey("When parse with flags", func() {
			p, err := hyperscan.ParsePattern(`/test/im`)

			So(err, ShouldBeNil)
			So(p, ShouldNotBeNil)
			So(p.Expression, ShouldEqual, "test")
			So(p.Flags, ShouldEqual, hyperscan.Caseless|hyperscan.MultiLine)

			So(p.Expression.String(), ShouldEqual, "test")
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

func TestDatabaseBuilder(t *testing.T) {
	Convey("Given a DatabaseBuilder", t, func() {
		b := hyperscan.DatabaseBuilder{}

		Convey("When build without patterns", func() {
			db, err := b.Build()

			So(db, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})

		Convey("When build with a simple expression", func() {
			db, err := b.AddExpressions("test").Build()

			So(err, ShouldBeNil)
			So(db, ShouldNotBeNil)

			So(db.Close(), ShouldBeNil)
		})

		Convey("When build with some complicated expression", func() {
			db, err := b.AddExpressions("test", hyperscan.EmailAddress, hyperscan.IPv4Address, hyperscan.CreditCard).Build()

			So(err, ShouldBeNil)
			So(db, ShouldNotBeNil)

			info, err := db.Info()

			So(err, ShouldBeNil)

			mode, err := info.Mode()

			So(err, ShouldBeNil)
			So(mode, ShouldEqual, hyperscan.BlockMode)

			So(db.Close(), ShouldBeNil)
		})

		Convey("When build stream database with a simple expression", func() {
			b.Mode = hyperscan.StreamMode

			db, err := b.AddExpressionWithFlags("test", hyperscan.Caseless).Build()

			So(err, ShouldBeNil)
			So(db, ShouldNotBeNil)

			info, err := db.Info()

			So(err, ShouldBeNil)

			mode, err := info.Mode()

			So(err, ShouldBeNil)
			So(mode, ShouldEqual, hyperscan.StreamMode)

			So(db.Close(), ShouldBeNil)
		})

		Convey("When build vectored database with a simple expression", func() {
			b.Mode = hyperscan.VectoredMode

			db, err := b.AddExpressions("test").Build()

			So(err, ShouldBeNil)
			So(db, ShouldNotBeNil)

			info, err := db.Info()

			So(err, ShouldBeNil)

			mode, err := info.Mode()

			So(err, ShouldBeNil)
			So(mode, ShouldEqual, hyperscan.VectoredMode)

			So(db.Close(), ShouldBeNil)
		})
	})
}

func TestCompile(t *testing.T) {
	Convey("Given compile some expressions", t, func() {
		Convey("When compile a simple expression", func() {
			db, err := hyperscan.Compile("test")

			So(db, ShouldNotBeNil)
			So(err, ShouldBeNil)

			So(db.Close(), ShouldBeNil)
		})

		Convey("When compile a complex expression", func() {
			db, err := hyperscan.Compile(hyperscan.CreditCard)

			So(db, ShouldNotBeNil)
			So(err, ShouldBeNil)

			So(db.Close(), ShouldBeNil)
		})
	})
}

func TestPlatform(t *testing.T) {
	Convey("Given a native platform", t, func() {
		p := hyperscan.PopulatePlatform()

		So(p, ShouldNotBeNil)
		So(p.Tune(), ShouldBeGreaterThan, hyperscan.Generic)
		So(p.CpuFeatures(), ShouldBeGreaterThanOrEqualTo, 0)

		So(p, ShouldResemble, hyperscan.NewPlatform(p.Tune(), p.CpuFeatures()))
	})
}
