//go:build chimera
// +build chimera

package chimera_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/gohs/chimera"
	"github.com/flier/gohs/hyperscan"
)

func TestCompileFlag(t *testing.T) {
	Convey("Given a compile flags", t, func() {
		flags := chimera.Caseless | chimera.DotAll | chimera.MultiLine | chimera.SingleMatch |
			chimera.Utf8Mode | chimera.UnicodeProperty

		So(flags.String(), ShouldEqual, "8HWims")

		Convey("When parse valid flags", func() {
			f, err := chimera.ParseCompileFlag("ismH8W")

			So(f, ShouldEqual, flags)
			So(err, ShouldBeNil)
		})

		Convey("When parse invalid flags", func() {
			f, err := chimera.ParseCompileFlag("abc")

			So(f, ShouldEqual, 0)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestDatabaseBuilder(t *testing.T) {
	Convey("Given a DatabaseBuilder", t, func() {
		b := chimera.DatabaseBuilder{}

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
	})
}

func TestCompile(t *testing.T) {
	Convey("Given compile some expressions", t, func() {
		Convey("When compile a simple expression", func() {
			db, err := chimera.Compile("test")

			So(db, ShouldNotBeNil)
			So(err, ShouldBeNil)

			So(db.Close(), ShouldBeNil)
		})

		Convey("When compile a complex expression", func() {
			db, err := chimera.Compile(hyperscan.CreditCard)

			So(db, ShouldNotBeNil)
			So(err, ShouldBeNil)

			So(db.Close(), ShouldBeNil)
		})
	})
}
