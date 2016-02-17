package hyperscan

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBaseDatabase(t *testing.T) {
	Convey("Given a block database", t, func() {
		bdb, err := NewBlockDatabase(&Pattern{Expression: "test"})

		So(err, ShouldBeNil)
		So(bdb, ShouldNotBeNil)

		Convey("When get size", func() {
			size, err := bdb.Size()

			So(err, ShouldBeNil)
			So(size, ShouldEqual, 1000)
		})

		Convey("When get info", func() {
			info, err := bdb.Info()

			So(err, ShouldBeNil)
			So(info, ShouldNotBeNil)

			So(regexInfo.MatchString(info.String()), ShouldBeTrue)

			Convey("Then get version", func() {
				ver, err := info.Version()

				So(err, ShouldBeNil)
				So(Version(), ShouldStartWith, ver)
			})

			Convey("Then get mode", func() {
				mode, err := info.Mode()

				So(err, ShouldBeNil)
				So(mode, ShouldEqual, BlockMode)
			})
		})

		Convey("When serialize database", func() {
			data, err := bdb.Marshal()

			So(err, ShouldBeNil)
			So(len(data), ShouldEqual, 1000)

			Convey("When get size", func() {
				size, err := DatabaseSize(data)

				So(err, ShouldBeNil)
				So(size, ShouldEqual, 1000)
			})

			Convey("When get info", func() {
				info, err := DatabaseInfo(data)

				So(err, ShouldBeNil)
				So(info, ShouldNotBeNil)

				So(regexInfo.MatchString(info.String()), ShouldBeTrue)
			})

			Convey("Then deserialize database", func() {
				db, err := UnmarshalDatabase(data)

				So(err, ShouldBeNil)
				So(db, ShouldNotBeNil)

				Convey("When get info", func() {
					info, err := bdb.Info()

					So(err, ShouldBeNil)
					So(info, ShouldNotBeNil)

					So(regexInfo.MatchString(info.String()), ShouldBeTrue)

					Convey("Then get version", func() {
						ver, err := info.Version()

						So(err, ShouldBeNil)
						So(Version(), ShouldStartWith, ver)
					})
				})

				So(db.Close(), ShouldBeNil)
			})
		})

		So(bdb.Close(), ShouldBeNil)
	})
}

func TestBlockDatabase(t *testing.T) {
	Convey("Give a block database", t, func() {
		bdb, err := NewBlockDatabase(&Pattern{Expression: "test"})

		So(err, ShouldBeNil)
		So(bdb, ShouldNotBeNil)

		Convey("When get info", func() {
			info, err := bdb.Info()

			So(err, ShouldBeNil)
			So(info, ShouldNotBeNil)

			So(regexInfo.MatchString(info.String()), ShouldBeTrue)

			Convey("Then get mode", func() {
				mode, err := info.Mode()

				So(err, ShouldBeNil)
				So(mode, ShouldEqual, BlockMode)
			})
		})

		So(bdb.Close(), ShouldBeNil)
	})
}

func TestVectoredDatabase(t *testing.T) {
	Convey("Give a vectored database", t, func() {
		bdb, err := NewVectoredDatabase(&Pattern{Expression: "test"})

		So(err, ShouldBeNil)
		So(bdb, ShouldNotBeNil)

		Convey("When get info", func() {
			info, err := bdb.Info()

			So(err, ShouldBeNil)
			So(info, ShouldNotBeNil)

			So(regexInfo.MatchString(info.String()), ShouldBeTrue)

			Convey("Then get mode", func() {
				mode, err := info.Mode()

				So(err, ShouldBeNil)
				So(mode, ShouldEqual, VectoredMode)
			})
		})

		So(bdb.Close(), ShouldBeNil)
	})
}

func TestStreamDatabase(t *testing.T) {
	Convey("Give a stream database", t, func() {
		bdb, err := NewStreamDatabase(&Pattern{Expression: "test"})

		So(err, ShouldBeNil)
		So(bdb, ShouldNotBeNil)

		Convey("When get info", func() {
			info, err := bdb.Info()

			So(err, ShouldBeNil)
			So(info, ShouldNotBeNil)

			So(regexInfo.MatchString(info.String()), ShouldBeTrue)

			Convey("Then get mode", func() {
				mode, err := info.Mode()

				So(err, ShouldBeNil)
				So(mode, ShouldEqual, StreamMode)
			})
		})

		So(bdb.Close(), ShouldBeNil)
	})
}
