package hyperscan

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPattern(t *testing.T) {
	Convey("Give a pattern", t, func() {
		Convey("When parse with flags", func() {
			p, err := ParsePattern(`/test/im`)

			So(err, ShouldBeNil)
			So(p, ShouldNotBeNil)
			So(p.Expression, ShouldEqual, "test")
			So(p.Flags, ShouldEqual, Caseless|MultiLine)

			Convey("When pattern contains forward slash", func() {
				p, err := ParsePattern(`/te/st/im`)

				So(err, ShouldBeNil)
				So(p, ShouldNotBeNil)
				So(p.Expression, ShouldEqual, "te/st")
				So(p.Flags, ShouldEqual, Caseless|MultiLine)

				So(p.String(), ShouldEqual, "/te/st/im")
			})
		})

		Convey("When parse with a lot of flags", func() {
			p, err := ParsePattern(`/test/ismoeupf`)

			So(err, ShouldBeNil)
			So(p, ShouldNotBeNil)
			So(p.Expression, ShouldEqual, "test")
			So(p.Flags, ShouldEqual, Caseless|DotAll|MultiLine|SingleMatch|AllowEmpty|Utf8Mode|UnicodeProperty|PrefilterMode)

			So(p.Flags.String(), ShouldEqual, "efimopsu")
			So(p.String(), ShouldEqual, "/test/efimopsu")
		})

		Convey("When parse without flags", func() {
			p, err := ParsePattern(`test`)

			So(err, ShouldBeNil)
			So(p, ShouldNotBeNil)
			So(p.Expression, ShouldEqual, "test")
			So(p.Flags, ShouldEqual, 0)

			Convey("When pattern contains forward slash", func() {
				p, err := ParsePattern(`te/st`)

				So(err, ShouldBeNil)
				So(p, ShouldNotBeNil)
				So(p.Expression, ShouldEqual, "te/st")
				So(p.Flags, ShouldEqual, 0)
			})
		})

		Convey("When pattern is valid", func() {
			p := Pattern{Expression: "test"}

			info, err := p.Info()

			So(err, ShouldBeNil)
			So(info, ShouldNotBeNil)
			So(info, ShouldResemble, &ExprInfo{MinWidth: 4, MaxWidth: 4})
			So(p.IsValid(), ShouldBeTrue)
		})

		Convey("When pattern is invalid", func() {
			p := Pattern{Expression: `\R`}

			info, err := p.Info()

			So(err, ShouldNotBeNil)
			So(info, ShouldBeNil)
			So(p.IsValid(), ShouldBeFalse)
		})
	})
}

func TestDatabaseBuilder(t *testing.T) {
	Convey("Given a DatabaseBuilder", t, func() {
		b := DatabaseBuilder{}

		Convey("When build without patterns", func() {
			db, err := b.Build()

			So(db, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})

		Convey("When build with a simple expression", func() {
			db, err := b.AddPatterns("test").Build()

			So(err, ShouldBeNil)
			So(db, ShouldNotBeNil)

			So(db.Close(), ShouldBeNil)
		})

		Convey("When build with some complicated expression", func() {
			db, err := b.AddPatterns("test", EmailAddress, IPv4Address, CreditCard).Build()

			So(err, ShouldBeNil)
			So(db, ShouldNotBeNil)

			info, err := db.Info()

			So(err, ShouldBeNil)

			mode, err := info.Mode()

			So(err, ShouldBeNil)
			So(mode, ShouldEqual, Block)

			So(db.Close(), ShouldBeNil)
		})

		Convey("When build stream database with a simple expression", func() {
			b.Mode = Stream

			db, err := b.AddPatterns("test").Build()

			So(err, ShouldBeNil)
			So(db, ShouldNotBeNil)

			info, err := db.Info()

			So(err, ShouldBeNil)

			mode, err := info.Mode()

			So(err, ShouldBeNil)
			So(mode, ShouldEqual, Stream)

			So(db.Close(), ShouldBeNil)
		})

		Convey("When build vectored database with a simple expression", func() {
			b.Mode = Vectored

			db, err := b.AddPatterns("test").Build()

			So(err, ShouldBeNil)
			So(db, ShouldNotBeNil)

			info, err := db.Info()

			So(err, ShouldBeNil)

			mode, err := info.Mode()

			So(err, ShouldBeNil)
			So(mode, ShouldEqual, Vectored)

			So(db.Close(), ShouldBeNil)
		})
	})
}

func TestCompile(t *testing.T) {
	Convey("Given compile some expressions", t, func() {
		Convey("When compile a simple expression", func() {
			db, err := Compile("test")

			So(db, ShouldNotBeNil)
			So(err, ShouldBeNil)

			So(db.Close(), ShouldBeNil)
		})

		Convey("When compile a complex expression", func() {
			db, err := Compile(CreditCard)

			So(db, ShouldNotBeNil)
			So(err, ShouldBeNil)

			So(db.Close(), ShouldBeNil)
		})
	})
}
