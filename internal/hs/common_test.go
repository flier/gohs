package hs_test

import (
	"regexp"
	"testing"
	"unsafe"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/gohs/internal/hs"
)

func TestVersion(t *testing.T) {
	Convey("Given a HyperScan version", t, func() {
		ver := hs.Version()

		So(ver, ShouldNotBeEmpty)

		matched, err := regexp.MatchString(`^\d\.\d\.\d.*`, ver)

		So(err, ShouldBeNil)
		So(matched, ShouldBeTrue)
	})
}

var regexInfo = regexp.MustCompile(`^Version: (\d+\.\d+\.\d+) Features: ([\w\s]+)? Mode: (\w+)$`)

//nolint:funlen
func TestDatabase(t *testing.T) {
	Convey("Given a stream database", t, func() {
		platform, err := hs.PopulatePlatform()

		So(platform, ShouldNotBeNil)
		So(err, ShouldBeNil)

		db, err := hs.Compile("test", 0, hs.StreamMode, platform)

		So(db, ShouldNotBeNil)
		So(err, ShouldBeNil)

		Convey("Get the database info", func() {
			info, err := hs.DatabaseInfo(db)

			So(regexInfo.MatchString(info), ShouldBeTrue)
			So(err, ShouldBeNil)
		})

		Convey("Get the database size", func() {
			size, err := hs.DatabaseSize(db)

			So(size, ShouldBeGreaterThan, 800)
			So(err, ShouldBeNil)
		})

		Convey("Get the stream size", func() {
			size, err := hs.StreamSize(db)

			So(size, ShouldBeGreaterThan, 20)
			So(err, ShouldBeNil)
		})

		Convey("Get the stream size from a block database", func() {
			db, err := hs.Compile("test", 0, hs.BlockMode, platform)

			So(db, ShouldNotBeNil)
			So(err, ShouldBeNil)

			size, err := hs.StreamSize(db)

			So(size, ShouldEqual, 0)
			So(err, ShouldEqual, hs.ErrDatabaseModeError)
		})

		Convey("When serialize database", func() {
			data, err := hs.SerializeDatabase(db)

			So(data, ShouldNotBeNil)
			So(len(data), ShouldBeGreaterThan, 800)
			So(err, ShouldBeNil)

			Convey("Get the database info", func() {
				info, err := hs.SerializedDatabaseInfo(data)

				So(regexInfo.MatchString(info), ShouldBeTrue)
				So(err, ShouldBeNil)
			})

			Convey("Get the database size", func() {
				size, err := hs.SerializedDatabaseSize(data)

				So(size, ShouldBeGreaterThan, 800)
				So(err, ShouldBeNil)
			})

			Convey("Then deserialize database", func() {
				db, err := hs.DeserializeDatabase(data)

				So(db, ShouldNotBeNil)
				So(err, ShouldBeNil)

				Convey("Get the database info", func() {
					info, err := hs.DatabaseInfo(db)

					So(regexInfo.MatchString(info), ShouldBeTrue)
					So(err, ShouldBeNil)
				})
			})

			Convey("Then deserialize database to memory", func() {
				buf := make([]byte, 1000)
				db := hs.Database(unsafe.Pointer(&buf[0]))

				So(hs.DeserializeDatabaseAt(data, db), ShouldBeNil)

				Convey("Get the database info", func() {
					info, err := hs.DatabaseInfo(db)

					So(regexInfo.MatchString(info), ShouldBeTrue)
					So(err, ShouldBeNil)
				})
			})
		})

		So(hs.FreeDatabase(db), ShouldBeNil)
	})
}
