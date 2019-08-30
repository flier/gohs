// +build !hyperscan_v4,!hyperscan_v5_1

package hyperscan

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLiteral(t *testing.T) {
	Convey("Give a literal", t, func() {
		Convey("When parse with flags", func() {
			p, err := ParseLiteral(`/test/im`)

			So(err, ShouldBeNil)
			So(p, ShouldNotBeNil)
			So(p.Expression, ShouldResemble, []byte("test"))
			So(p.Flags, ShouldEqual, Caseless|MultiLine)

			So(p.Expression, ShouldResemble, []byte("test"))
			So(p.String(), ShouldEqual, `/test/im`)

			Convey("When pattern contains forward slash", func() {
				p, err := ParseLiteral(`/te/st/im`)

				So(err, ShouldBeNil)
				So(p, ShouldNotBeNil)
				So(p.Expression, ShouldResemble, []byte("te/st"))
				So(p.Flags, ShouldEqual, Caseless|MultiLine)

				So(p.String(), ShouldEqual, "/te/st/im")
			})

			Convey("When pattern contains NULL", func() {
				p, err := ParseLiteral("/te\u0000st/im")

				So(err, ShouldBeNil)
				So(p, ShouldNotBeNil)
				So(p.Expression, ShouldResemble, []byte{116, 101, 0, 115, 116})
				So(p.Flags, ShouldEqual, Caseless|MultiLine)

				So(p.String(), ShouldEqual, "/te\\x00st/im")
			})
		})
	})
}

func TestLiteralDatabaseBuilder(t *testing.T) {
	Convey("Given a LiteralDatabaseBuilder", t, func() {
		b := LiteralDatabaseBuilder{
			Literals: []*Literal{
				NewLiteral([]byte("foo"), 0),
				NewLiteral([]byte("bar"), 0),
			},
		}

		db, err := b.Build()

		So(err, ShouldBeNil)
		So(db, ShouldNotBeNil)

		info, err := db.Info()

		So(err, ShouldBeNil)

		mode, err := info.Mode()

		So(err, ShouldBeNil)
		So(mode, ShouldEqual, BlockMode)

		So(db.Close(), ShouldBeNil)
	})
}
