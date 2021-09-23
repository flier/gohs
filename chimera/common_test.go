package chimera_test

import (
	"testing"

	"github.com/flier/gohs/chimera"
	. "github.com/smartystreets/goconvey/convey"
)

func TestChimera(t *testing.T) {
	Convey("Given a chimera runtimes", t, func() {
		So(chimera.Version(), ShouldNotBeEmpty)
	})
}
