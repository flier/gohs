package ch_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/gohs/chimera"
	"github.com/flier/gohs/internal/ch"
	"github.com/flier/gohs/internal/hs"
)

func TestCompileAPI(t *testing.T) {
	Convey("Given a host platform", t, func() {
		platform, err := hs.PopulatePlatform()

		So(platform, ShouldNotBeNil)
		So(err, ShouldBeNil)

		Convey("Compile an empty expression", func() {
			db, err := ch.Compile("", 0, ch.Groups, platform)

			So(err, ShouldBeNil)
			So(db, ShouldNotBeNil)

			So(ch.FreeDatabase(db), ShouldBeNil)
		})

		Convey("Compile multi expressions", func() {
			db, err := ch.CompileMulti(chimera.Patterns{
				chimera.NewPattern(`^\w+`, 0),
				chimera.NewPattern(`\d+`, 0),
				chimera.NewPattern(`\s+`, 0),
			}, ch.Groups, platform)

			So(db, ShouldNotBeNil)
			So(err, ShouldBeNil)

			Convey("Get the database info", func() {
				_, err := ch.DatabaseInfo(db)

				So(err, ShouldBeNil)
			})

			So(ch.FreeDatabase(db), ShouldBeNil)
		})

		Convey("Compile multi expressions with extension", func() {
			db, err := ch.CompileExtMulti(chimera.Patterns{
				chimera.NewPattern(`^\w+`, 0),
				chimera.NewPattern(`\d+`, 0),
				chimera.NewPattern(`\s+`, 0),
			}, ch.Groups, platform, 1, 1)

			So(db, ShouldNotBeNil)
			So(err, ShouldBeNil)

			Convey("Get the database info", func() {
				_, err := ch.DatabaseInfo(db)

				So(err, ShouldBeNil)
			})

			So(ch.FreeDatabase(db), ShouldBeNil)
		})
	})
}
