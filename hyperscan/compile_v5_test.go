// +build !hyperscan_v4

package hyperscan

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

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
