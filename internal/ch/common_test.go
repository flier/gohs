package ch_test

import (
	"regexp"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/gohs/internal/ch"
)

func TestVersion(t *testing.T) {
	Convey("Given a Chimera version", t, func() {
		ver := ch.Version()

		So(ver, ShouldNotBeEmpty)

		matched, err := regexp.MatchString(`^\d\.\d\.\d.*`, ver)

		So(err, ShouldBeNil)
		So(matched, ShouldBeTrue)
	})
}
