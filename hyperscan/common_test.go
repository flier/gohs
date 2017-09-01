package hyperscan

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBaseDatabase(t *testing.T) {
	Convey("Given a block database", t, func() {
		So(ValidPlatform(), ShouldBeNil)

		bdb, err := NewBlockDatabase(&Pattern{Expression: "test"})

		So(err, ShouldBeNil)
		So(bdb, ShouldNotBeNil)

		Convey("When get size", func() {
			size, err := bdb.Size()

			So(err, ShouldBeNil)
			So(size, ShouldBeGreaterThan, 800)
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
			So(len(data), ShouldBeGreaterThan, 800)

			Convey("When get size", func() {
				size, err := SerializedDatabaseSize(data)

				So(err, ShouldBeNil)
				So(size, ShouldBeGreaterThan, 800)
			})

			Convey("When get info", func() {
				info, err := SerializedDatabaseInfo(data)

				So(err, ShouldBeNil)
				So(info, ShouldNotBeNil)

				So(regexInfo.MatchString(info.String()), ShouldBeTrue)
			})

			Convey("Then deserialize database", func() {
				db, err := UnmarshalBlockDatabase(data)

				So(err, ShouldBeNil)
				So(db, ShouldNotBeNil)

				Convey("When get info", func() {
					info, err := db.Info()

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

			Convey("Then deserialize database in place", func() {
				So(bdb.Unmarshal(data), ShouldBeNil)

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
		vdb, err := NewVectoredDatabase(&Pattern{Expression: "test"})

		So(err, ShouldBeNil)
		So(vdb, ShouldNotBeNil)

		Convey("When get info", func() {
			info, err := vdb.Info()

			So(err, ShouldBeNil)
			So(info, ShouldNotBeNil)

			So(regexInfo.MatchString(info.String()), ShouldBeTrue)

			Convey("Then get mode", func() {
				mode, err := info.Mode()

				So(err, ShouldBeNil)
				So(mode, ShouldEqual, VectoredMode)
			})
		})

		So(vdb.Close(), ShouldBeNil)
	})
}

func TestStreamDatabase(t *testing.T) {
	Convey("Give a stream database", t, func() {
		sdb, err := NewStreamDatabase(&Pattern{Expression: "test"})

		So(err, ShouldBeNil)
		So(sdb, ShouldNotBeNil)

		Convey("When get info", func() {
			info, err := sdb.Info()

			So(err, ShouldBeNil)
			So(info, ShouldNotBeNil)

			So(regexInfo.MatchString(info.String()), ShouldBeTrue)

			Convey("Then get mode", func() {
				mode, err := info.Mode()

				So(err, ShouldBeNil)
				So(mode, ShouldEqual, StreamMode)
			})
		})

		Convey("When get stream size", func() {
			size, err := sdb.StreamSize()

			So(err, ShouldBeNil)
			So(size, ShouldBeGreaterThan, 20)
		})

		So(sdb.Close(), ShouldBeNil)
	})
}
