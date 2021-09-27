package ch_test

import (
	"testing"
	"unsafe"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/gohs/internal/ch"
	"github.com/flier/gohs/internal/hs"
)

type testAllocator struct {
	memoryUsed  int
	memoryFreed []unsafe.Pointer
}

func (a *testAllocator) alloc(size uint) unsafe.Pointer {
	a.memoryUsed += int(size)

	return ch.AlignedAlloc(size)
}

func (a *testAllocator) free(ptr unsafe.Pointer) {
	a.memoryFreed = append(a.memoryFreed, ptr)

	ch.AlignedFree(ptr)
}

//nolint:funlen
func TestAllocator(t *testing.T) {
	Convey("Given the host platform", t, func() {
		platform, err := hs.PopulatePlatform()

		So(platform, ShouldNotBeNil)
		So(err, ShouldBeNil)

		a := &testAllocator{}

		Convey("Then create a database with allocator", func() {
			So(ch.SetDatabaseAllocator(a.alloc, a.free), ShouldBeNil)

			db, err := ch.Compile("test", 0, ch.Groups, platform)

			So(db, ShouldNotBeNil)
			So(err, ShouldBeNil)

			Convey("Then get the database info with allocator", func() {
				So(ch.SetMiscAllocator(a.alloc, a.free), ShouldBeNil)

				s, err := ch.DatabaseInfo(db)
				So(err, ShouldBeNil)
				So(s, ShouldNotBeEmpty)

				So(a.memoryUsed, ShouldBeGreaterThanOrEqualTo, 12)

				So(ch.ClearMiscAllocator(), ShouldBeNil)
			})

			Convey("Get the database size", func() {
				size, err := ch.DatabaseSize(db)

				So(a.memoryUsed, ShouldBeGreaterThanOrEqualTo, size)
				So(err, ShouldBeNil)
			})

			Convey("Then create a scratch with allocator", func() {
				So(ch.SetScratchAllocator(a.alloc, a.free), ShouldBeNil)

				a.memoryUsed = 0

				s, err := ch.AllocScratch(db)

				So(s, ShouldNotBeNil)
				So(err, ShouldBeNil)

				Convey("Get the scratch size", func() {
					size, err := ch.ScratchSize(s)

					So(a.memoryUsed, ShouldBeGreaterThanOrEqualTo, size)
					So(err, ShouldBeNil)
				})

				Convey("Then free scratch with allocator", func() {
					a.memoryFreed = nil

					So(ch.FreeScratch(s), ShouldBeNil)

					So(a.memoryFreed[len(a.memoryFreed)-1], ShouldEqual, unsafe.Pointer(s))

					So(ch.ClearScratchAllocator(), ShouldBeNil)
				})
			})

			Convey("Then free database with allocator", func() {
				a.memoryFreed = nil

				So(ch.FreeDatabase(db), ShouldBeNil)

				So(a.memoryFreed, ShouldResemble, []unsafe.Pointer{unsafe.Pointer(db)})

				So(ch.ClearDatabaseAllocator(), ShouldBeNil)
			})
		})
	})
}
