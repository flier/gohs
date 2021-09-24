package ch

import (
	"regexp"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestVersion(t *testing.T) {
	Convey("Given a Chimera version", t, func() {
		ver := Version()

		So(ver, ShouldNotBeEmpty)

		matched, err := regexp.MatchString(`^\d\.\d\.\d.*`, ver)

		So(err, ShouldBeNil)
		So(matched, ShouldBeTrue)
	})
}
