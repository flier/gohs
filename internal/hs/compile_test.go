package hs_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/gohs/hyperscan"
	"github.com/flier/gohs/internal/hs"
)

//nolint:funlen
func TestCompileAPI(t *testing.T) {
	Convey("Given a host platform", t, func() {
		platform, err := hs.PopulatePlatform()

		So(platform, ShouldNotBeNil)
		So(err, ShouldBeNil)

		Convey("When Compile a unsupported expression", func() {
			Convey("Then compile as stream", func() {
				db, err := hs.Compile(`\R`, 0, hs.StreamMode, platform)

				So(db, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, `\R at index 0 not supported.`)
			})

			Convey("Then compile as vector", func() {
				db, err := hs.CompileMulti(hyperscan.NewPattern(`\R`, hs.Caseless), hs.BlockMode, platform)

				So(db, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, `\R at index 0 not supported.`)
			})

			Convey("Then compile as extended vector", func() {
				db, err := hs.CompileMulti(hyperscan.NewPattern(`\R`, hs.Caseless, hyperscan.MinOffset(10)),
					hs.VectoredMode, platform)

				So(db, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, `\R at index 0 not supported.`)
			})
		})

		Convey("Compile an empty expression", func() {
			db, err := hs.Compile("", 0, hs.StreamMode, platform)

			So(db, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "Pattern matches empty buffer; use HS_FLAG_ALLOWEMPTY to enable support.")

			So(hs.FreeDatabase(db), ShouldBeNil)
		})

		Convey("Compile multi expressions", func() {
			db, err := hs.CompileMulti(hyperscan.Patterns{
				hyperscan.NewPattern(`^\w+`, 0),
				hyperscan.NewPattern(`\d+`, 0),
				hyperscan.NewPattern(`\s+`, 0),
			}, hs.StreamMode, platform)

			So(db, ShouldNotBeNil)
			So(err, ShouldBeNil)

			Convey("Get the database info", func() {
				info, err := hs.DatabaseInfo(db)

				So(regexInfo.MatchString(info), ShouldBeTrue)
				So(err, ShouldBeNil)
			})

			So(hs.FreeDatabase(db), ShouldBeNil)
		})

		Convey("Compile multi expressions with extension", func() {
			db, err := hs.CompileMulti(hyperscan.Patterns{
				hyperscan.NewPattern(`^\w+`, 0, hyperscan.MinOffset(10)),
				hyperscan.NewPattern(`\d+`, 0, hyperscan.MaxOffset(10)),
				hyperscan.NewPattern(`\s+`, 0, hyperscan.MinLength(10)),
			}, hs.StreamMode, platform)

			So(db, ShouldNotBeNil)
			So(err, ShouldBeNil)

			Convey("Get the database info", func() {
				info, err := hs.DatabaseInfo(db)

				So(regexInfo.MatchString(info), ShouldBeTrue)
				So(err, ShouldBeNil)
			})

			So(hs.FreeDatabase(db), ShouldBeNil)
		})
	})
}

func TestExpression(t *testing.T) {
	Convey("Given a simple expression", t, func() {
		info, err := hs.ExpressionInfo("test", 0)

		So(info, ShouldNotBeNil)
		So(info, ShouldResemble, &hs.ExprInfo{
			MinWidth: 4,
			MaxWidth: 4,
		})
		So(err, ShouldBeNil)
	})

	Convey("Given a credit card expression", t, func() {
		info, err := hs.ExpressionInfo(hyperscan.CreditCard, 0)

		So(info, ShouldNotBeNil)
		So(info, ShouldResemble, &hs.ExprInfo{
			MinWidth: 13,
			MaxWidth: 16,
		})
		So(err, ShouldBeNil)
	})

	Convey("Given a expression match eod", t, func() {
		info, err := hs.ExpressionInfo("test$", 0)

		So(info, ShouldNotBeNil)
		So(info, ShouldResemble, &hs.ExprInfo{
			MinWidth:        4,
			MaxWidth:        4,
			ReturnUnordered: true,
			AtEndOfData:     true,
			OnlyAtEndOfData: true,
		})
		So(err, ShouldBeNil)
	})
}
