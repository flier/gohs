package hyperscan

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBlockScanner(t *testing.T) {
	Convey("Given a block database", t, func() {
		bdb, err := NewBlockDatabase(NewPattern(`\d+`, SomLeftMost))

		So(err, ShouldBeNil)
		So(bdb, ShouldNotBeNil)

		Convey("When scan a string", func() {
			var matches [][]uint64

			matched := func(id uint, from, to uint64, flags uint, context interface{}) error {
				matches = append(matches, []uint64{from, to})

				return nil
			}

			err = bdb.Scan([]byte("abc123def456"), nil, matched, nil)

			So(err, ShouldBeNil)
			So(matches, ShouldResemble, [][]uint64{{3, 4}, {3, 5}, {3, 6}, {9, 10}, {9, 11}, {9, 12}})
		})
	})
}

func TestBlockMatcher(t *testing.T) {
	Convey("Given a block database", t, func() {
		bdb, err := NewBlockDatabase(NewPattern(`\d+`, SomLeftMost))

		So(err, ShouldBeNil)
		So(bdb, ShouldNotBeNil)

		Convey("When match the string", func() {
			So(bdb.MatchString("abc123def456"), ShouldBeTrue)
		})

		Convey("When find the leftmost matched string index", func() {
			So(bdb.FindStringIndex("abc123def456"), ShouldResemble, []int{3, 4})
		})

		Convey("When find the leftmost matched string", func() {
			So(bdb.FindString("abc123def456"), ShouldEqual, "1")
		})

		Convey("When find all the matched string index", func() {
			So(bdb.FindAllStringIndex("abc123def456", -1), ShouldResemble,
				[][]int{{3, 4}, {3, 5}, {3, 6}, {9, 10}, {9, 11}, {9, 12}})
		})

		Convey("When find all the matched string", func() {
			So(bdb.FindAllString("abc123def456", -1), ShouldResemble,
				[]string{"1", "12", "123", "4", "45", "456"})
		})

		Convey("When find all the first 4 matched string index", func() {
			So(bdb.FindAllStringIndex("abc123def456", 4), ShouldResemble,
				[][]int{{3, 4}, {3, 5}, {3, 6}, {9, 10}})
		})

		Convey("When find all the first 4 matched string", func() {
			So(bdb.FindAllString("abc123def456", 4), ShouldResemble,
				[]string{"1", "12", "123", "4"})
		})
	})
}
