package hyperscan_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/gohs/hyperscan"
)

func TestPlatform(t *testing.T) {
	Convey("Given a native platform", t, func() {
		p := hyperscan.PopulatePlatform()

		So(p, ShouldNotBeNil)
		So(p.Tune(), ShouldBeGreaterThan, hyperscan.Generic)
		So(p.CpuFeatures(), ShouldBeGreaterThanOrEqualTo, 0)

		So(p, ShouldResemble, hyperscan.NewPlatform(p.Tune(), p.CpuFeatures()))
	})
}
