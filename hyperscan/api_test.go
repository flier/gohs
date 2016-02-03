package hyperscan

import (
	"bytes"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMatch(t *testing.T) {
	Convey("Given a compatible API for regexp package", t, func() {
		Convey("When match a simple expression", func() {
			matched, err := MatchString("test", "abctestdef")

			So(matched, ShouldBeTrue)
			So(err, ShouldBeNil)
		})

		Convey("When match a invalid expression", func() {
			matched, err := MatchString(`\R`, "abctestdef")

			So(matched, ShouldBeFalse)
			So(err, ShouldNotBeNil)
		})

		Convey("When match a simple expression with io.Reader", func() {
			var buf bytes.Buffer

			buf.Write(make([]byte, 1024*1024))
			buf.WriteString("abctestdef")

			matched, err := MatchReader("test", &buf)

			So(matched, ShouldBeTrue)
			So(err, ShouldBeNil)
		})
	})
}
