package chimera_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/gohs/chimera"
)

func TestChimera(t *testing.T) {
	Convey("Given a chimera runtimes", t, func() {
		So(chimera.Version(), ShouldNotBeEmpty)
	})
}
