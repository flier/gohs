//go:build chimera
// +build chimera

package chimera_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/gohs/chimera"
)

type (
	BlockDatabaseConstructor func(patterns ...*chimera.Pattern) (chimera.BlockDatabase, error)
)

var blockDatabaseConstructors = map[string]BlockDatabaseConstructor{
	"normal":  chimera.NewBlockDatabase,
	"managed": chimera.NewManagedBlockDatabase,
}

func TestBlockScanner(t *testing.T) {
	for dbType, dbConstructor := range blockDatabaseConstructors {
		Convey("Given a "+dbType+" block database", t, func() {
			bdb, err := dbConstructor(chimera.NewPattern(`\d+`, 0)) // nolint: scopelint

			So(err, ShouldBeNil)
			So(bdb, ShouldNotBeNil)

			Convey("When scan a string", func() {
				var matches [][]uint64

				matched := func(id uint, from, to uint64, flags uint,
					captured []*chimera.Capture, context interface{},
				) chimera.Callback {
					matches = append(matches, []uint64{from, to})

					return chimera.Continue
				}

				err = bdb.Scan([]byte("abc123def456"), nil, chimera.HandlerFunc(matched), nil)

				So(err, ShouldBeNil)
				So(matches, ShouldResemble, [][]uint64{{3, 6}, {9, 12}})
			})
		})
	}
}

func TestBlockMatcher(t *testing.T) {
	for dbType, dbConstructor := range blockDatabaseConstructors {
		Convey("Given a "+dbType+" block database", t, func() {
			bdb, err := dbConstructor(chimera.NewPattern(`\d+`, 0)) // nolint: scopelint

			So(err, ShouldBeNil)
			So(bdb, ShouldNotBeNil)

			Convey("When match the string", func() {
				So(bdb.MatchString("123"), ShouldBeTrue)
				So(bdb.MatchString("abc123def456"), ShouldBeTrue)
			})

			Convey("When find the leftmost matched string index", func() {
				So(bdb.FindStringIndex("123"), ShouldResemble, []int{0, 3})
				So(bdb.FindStringIndex("abc123def456"), ShouldResemble, []int{3, 6})
			})

			Convey("When find the leftmost matched string", func() {
				So(bdb.FindString("123"), ShouldEqual, "123")
				So(bdb.FindString("abc123def456"), ShouldEqual, "123")
			})

			Convey("When find all the matched string index", func() {
				So(bdb.FindAllStringIndex("123", -1), ShouldResemble, [][]int{{0, 3}})
				So(bdb.FindAllStringIndex("abc123def456", -1), ShouldResemble, [][]int{{3, 6}, {9, 12}})
			})

			Convey("When find all the matched string", func() {
				So(bdb.FindAllString("abc123def456", -1), ShouldResemble,
					[]string{"123", "456"})
			})

			Convey("When find all the first 4 matched string index", func() {
				So(bdb.FindAllStringIndex("abc123def456", 1), ShouldResemble,
					[][]int{{3, 6}})
			})

			Convey("When find all the first 4 matched string", func() {
				So(bdb.FindAllString("abc123def456", 1), ShouldResemble,
					[]string{"123"})
			})
		})
	}
}
