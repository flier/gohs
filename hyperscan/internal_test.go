package hyperscan

import (
	"testing"
	"unsafe"

	. "github.com/smartystreets/goconvey/convey"
)

func TestVersion(t *testing.T) {
	Convey("Given a HyperScan version", t, func() {
		ver := hsVersion()

		So(ver, ShouldNotBeEmpty)
		So(ver, ShouldStartWith, "4.")
	})
}

func TestDatabase(t *testing.T) {
	Convey("Given a stream database", t, func() {
		platform, err := hsPopulatePlatform()

		So(platform, ShouldNotBeNil)
		So(platform.info, ShouldNotBeNil)
		So(err, ShouldBeNil)

		db, err := hsCompile("test", 0, Stream, platform)

		So(db, ShouldNotBeNil)
		So(err, ShouldBeNil)

		Convey("Get the database info", func() {
			info, err := hsDatabaseInfo(db)

			So(info, ShouldStartWith, "Version: 4.")
			So(info, ShouldEndWith, "Features:  AVX2 Mode: STREAM")
			So(err, ShouldBeNil)
		})

		Convey("Get the database size", func() {
			size, err := hsDatabaseSize(db)

			So(size, ShouldEqual, 1000)
			So(err, ShouldBeNil)
		})

		Convey("Get the stream size", func() {
			size, err := hsStreamSize(db)

			So(size, ShouldEqual, 24)
			So(err, ShouldBeNil)
		})

		Convey("Get the stream size from a block database", func() {
			db, err := hsCompile("test", 0, Block, platform)

			So(db, ShouldNotBeNil)
			So(err, ShouldBeNil)

			size, err := hsStreamSize(db)

			So(size, ShouldEqual, 0)
			So(err, ShouldEqual, DatabaseModeError)
		})

		Convey("When serialize database", func() {
			data, err := hsSerializeDatabase(db)

			So(data, ShouldNotBeNil)
			So(len(data), ShouldEqual, 1000)
			So(err, ShouldBeNil)

			Convey("Get the database info", func() {
				info, err := hsSerializedDatabaseInfo(data)

				So(info, ShouldStartWith, "Version: 4.")
				So(info, ShouldEndWith, "Features:  AVX2 Mode: STREAM")
				So(err, ShouldBeNil)
			})

			Convey("Get the database size", func() {
				size, err := hsSerializedDatabaseSize(data)

				So(size, ShouldEqual, 1000)
				So(err, ShouldBeNil)
			})

			Convey("Then deserialize database", func() {
				db, err := hsDeserializeDatabase(data)

				So(db, ShouldNotBeNil)
				So(err, ShouldBeNil)

				Convey("Get the database info", func() {
					info, err := hsDatabaseInfo(db)

					So(info, ShouldStartWith, "Version: 4.")
					So(info, ShouldEndWith, "Features:  AVX2 Mode: STREAM")
					So(err, ShouldBeNil)
				})
			})

			Convey("Then deserialize database to memory", func() {
				buf := make([]byte, 1000)
				db := hsDatabase(unsafe.Pointer(&buf[0]))

				So(hsDeserializeDatabaseAt(data, db), ShouldBeNil)

				Convey("Get the database info", func() {
					info, err := hsDatabaseInfo(db)

					So(info, ShouldStartWith, "Version: 4.")
					So(info, ShouldEndWith, "Features:  AVX2 Mode: STREAM")
					So(err, ShouldBeNil)
				})
			})
		})

		So(hsFreeDatabase(db), ShouldBeNil)
	})
}

func TestCompile(t *testing.T) {
	Convey("Compile a unsupported expression", t, func() {
		platform, err := hsPopulatePlatform()

		So(platform, ShouldNotBeNil)
		So(platform.info, ShouldNotBeNil)
		So(err, ShouldBeNil)

		db, err := hsCompile(`\R`, 0, Stream, platform)

		So(db, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, `\R at index 0 not supported.`)
	})

	Convey("Compile an empty expression", t, func() {
		platform, err := hsPopulatePlatform()

		So(platform, ShouldNotBeNil)
		So(platform.info, ShouldNotBeNil)
		So(err, ShouldBeNil)

		db, err := hsCompile("", 0, Stream, platform)

		So(db, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "Pattern matches empty buffer; use HS_FLAG_ALLOWEMPTY to enable support.")
	})
}
