package hyperscan_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/gohs/hyperscan"
)

func TestFloatNumber(t *testing.T) {
	Convey("Given a compiled pattern", t, func() {
		db := hyperscan.MustCompile(hyperscan.FloatNumber)

		So(db, ShouldNotBeNil)
	})
}

func TestIPv4Address(t *testing.T) {
	Convey("Given a compiled pattern", t, func() {
		db := hyperscan.MustCompile(hyperscan.IPv4Address)

		So(db, ShouldNotBeNil)
	})
}

func TestEmailAddress(t *testing.T) {
	Convey("Given a compiled pattern", t, func() {
		db := hyperscan.MustCompile(hyperscan.EmailAddress)

		So(db, ShouldNotBeNil)
	})
}

func TestCreditCard(t *testing.T) {
	Convey("Given a compiled pattern", t, func() {
		db := hyperscan.MustCompile(hyperscan.CreditCard)

		So(db, ShouldNotBeNil)
	})
}
