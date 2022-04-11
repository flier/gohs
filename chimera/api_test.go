//go:build chimera
// +build chimera

package chimera_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/gohs/chimera"
)

func TestMatch(t *testing.T) {
	Convey("Given a compatible API for regexp package", t, func() {
		Convey("When match a simple expression", func() {
			matched, err := chimera.MatchString("test", "abctestdef")

			So(matched, ShouldBeTrue)
			So(err, ShouldBeNil)
		})

		Convey("When match a invalid expression", func() {
			matched, err := chimera.MatchString(`(?R`, "abctestdef")

			So(matched, ShouldBeFalse)
			So(err, ShouldNotBeNil)
		})
	})
}
