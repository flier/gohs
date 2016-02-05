package hyperscan

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFloatNumber(t *testing.T) {
	Convey("Given a compiled pattern", t, func() {
		db := MustCompile(FloatNumber)

		So(db, ShouldNotBeNil)
	})
}

func TestIPv4Address(t *testing.T) {
	Convey("Given a compiled pattern", t, func() {
		db := MustCompile(IPv4Address)

		So(db, ShouldNotBeNil)
	})
}

func TestEmailAddress(t *testing.T) {
	Convey("Given a compiled pattern", t, func() {
		db := MustCompile(EmailAddress)

		So(db, ShouldNotBeNil)
	})
}

func TestCreditCard(t *testing.T) {
	Convey("Given a compiled pattern", t, func() {
		db := MustCompile(CreditCard)

		So(db, ShouldNotBeNil)
	})
}
