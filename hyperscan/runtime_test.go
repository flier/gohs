package hyperscan

import (
	"strings"
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
			So(bdb.FindStringIndex("abc123def456"), ShouldResemble, []int{3, 6})
		})

		Convey("When find the leftmost matched string", func() {
			So(bdb.FindString("abc123def456"), ShouldEqual, "123")
		})

		Convey("When find all the matched string index", func() {
			So(bdb.FindAllStringIndex("abc123def456", -1), ShouldResemble,
				[][]int{{3, 6}, {9, 12}})
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

func TestStreamScanner(t *testing.T) {
	Convey("Given a streaming database", t, func() {
		sdb, err := NewStreamDatabase(NewPattern(`abc`, SomLeftMost))

		So(err, ShouldBeNil)
		So(sdb, ShouldNotBeNil)

		Convey("When open a new stream", func() {
			var matches [][]uint64

			matched := func(id uint, from, to uint64, flags uint, context interface{}) error {
				matches = append(matches, []uint64{from, to})

				return nil
			}

			stream, err := sdb.Open(0, nil, matched, nil)

			So(err, ShouldBeNil)
			So(stream, ShouldNotBeNil)

			Convey("When scan a stream", func() {
				So(stream.Scan([]byte("123a")), ShouldBeNil)
				So(stream.Scan([]byte("b")), ShouldBeNil)
				So(stream.Scan([]byte("c456")), ShouldBeNil)
				So(stream.Close(), ShouldBeNil)

				So(matches, ShouldResemble, [][]uint64{{3, 6}})
			})
		})
	})
}

func TestStreamMatcher(t *testing.T) {
	Convey("Given a streaming database", t, func() {
		sdb, err := NewStreamDatabase(NewPattern(`\d+`, SomLeftMost))

		So(err, ShouldBeNil)
		So(sdb, ShouldNotBeNil)

		Convey("When scan a new stream", func() {
			r := strings.NewReader("foo123bar456")

			Convey("When `Match` a pattern", func() {
				So(sdb.Match(r), ShouldBeTrue)
			})

			Convey("When `Find` a pattern", func() {
				So(sdb.Find(r), ShouldResemble, []byte("123"))
			})

			Convey("When `FindIndex` a pattern", func() {
				So(sdb.FindIndex(r), ShouldResemble, []int{3, 6})
			})

			Convey("When `FindAll` a pattern", func() {
				So(sdb.FindAll(r, -1), ShouldResemble, [][]byte{[]byte("123"), []byte("456")})
			})

			Convey("When `FindAllIndex` a pattern", func() {
				So(sdb.FindAllIndex(r, -1), ShouldResemble, [][]int{{3, 6}, {9, 12}})
			})
		})
	})
}

func TestStreamCompressor(t *testing.T) {
	Convey("Given a streaming database", t, func() {
		sdb, err := NewStreamDatabase(NewPattern(`abc`, SomLeftMost))

		So(err, ShouldBeNil)
		So(sdb, ShouldNotBeNil)

		Convey("When open a new stream", func() {
			var matches [][]uint64

			matched := func(id uint, from, to uint64, flags uint, context interface{}) error {
				matches = append(matches, []uint64{from, to})

				return nil
			}

			stream, err := sdb.Open(0, nil, matched, nil)

			So(err, ShouldBeNil)
			So(stream, ShouldNotBeNil)

			defer stream.Close()

			So(stream.Scan([]byte("123a")), ShouldBeNil)

			Convey("When compress a stream", func() {
				buf, err := sdb.Compress(stream)

				So(err, ShouldBeNil)
				So(buf, ShouldNotBeNil)

				size, err := sdb.StreamSize()

				So(err, ShouldBeNil)
				So(len(buf), ShouldBeBetween, 0, size)

				Convey("When expand the stream", func() {
					stream2, err := sdb.Expand(buf, 0, nil, matched, nil)

					So(err, ShouldBeNil)
					So(stream2, ShouldNotBeNil)

					Convey("When scan a stream", func() {
						So(stream2.Scan([]byte("b")), ShouldBeNil)
						So(stream2.Scan([]byte("c456")), ShouldBeNil)
						So(stream2.Close(), ShouldBeNil)

						So(matches, ShouldResemble, [][]uint64{{3, 6}})
					})
				})

				Convey("When reset and expand the stream", func() {
					stream2, err := stream.Clone()

					So(err, ShouldBeNil)
					So(stream2, ShouldNotBeNil)

					stream2, err = sdb.ResetAndExpand(stream2, buf, 0, nil, matched, nil)

					So(err, ShouldBeNil)
					So(stream2, ShouldNotBeNil)

					Convey("When scan a stream", func() {
						So(stream2.Scan([]byte("b")), ShouldBeNil)
						So(stream2.Scan([]byte("c456")), ShouldBeNil)
						So(stream2.Close(), ShouldBeNil)

						So(matches, ShouldResemble, [][]uint64{{3, 6}})
					})
				})
			})
		})
	})
}
