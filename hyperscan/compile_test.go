//nolint:funlen
package hyperscan_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/gohs/hyperscan"
)

func TestCompileFlag(t *testing.T) {
	Convey("Given a compile flags", t, func() {
		flags := hyperscan.Caseless | hyperscan.DotAll | hyperscan.MultiLine | hyperscan.SingleMatch |
			hyperscan.AllowEmpty | hyperscan.Utf8Mode | hyperscan.UnicodeProperty | hyperscan.PrefilterMode

		So(flags.String(), ShouldEqual, "8HPVWims")

		Convey("When parse valid flags", func() {
			f, err := hyperscan.ParseCompileFlag("ifemopus")

			So(f, ShouldEqual, flags)
			So(err, ShouldBeNil)
		})

		Convey("When parse invalid flags", func() {
			f, err := hyperscan.ParseCompileFlag("abc")

			So(f, ShouldEqual, 0)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestModeFlag(t *testing.T) {
	Convey("Give a mode", t, func() {
		So(hyperscan.BlockMode.String(), ShouldEqual, "BLOCK")
		So(hyperscan.StreamMode.String(), ShouldEqual, "STREAM")
		So(hyperscan.VectoredMode.String(), ShouldEqual, "VECTORED")

		Convey("When combile mode with flags", func() {
			mode := hyperscan.StreamMode | hyperscan.SomHorizonLargeMode

			So(mode.String(), ShouldEqual, "STREAM")
		})

		Convey("When parse unknown mode", func() {
			m, err := hyperscan.ParseModeFlag("test")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "database mode test")
			So(m, ShouldEqual, hyperscan.BlockMode)
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
