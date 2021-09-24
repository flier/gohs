package hs_test

import (
	"testing"
	"unsafe"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/gohs/internal/hs"
)

type testAllocator struct {
	memoryUsed  int
	memoryFreed []unsafe.Pointer
}

func (a *testAllocator) alloc(size uint) unsafe.Pointer {
	a.memoryUsed += int(size)

	return hs.AlignedAlloc(size)
}

func (a *testAllocator) free(ptr unsafe.Pointer) {
	a.memoryFreed = append(a.memoryFreed, ptr)

	hs.AlignedFree(ptr)
}

//nolint:funlen
func TestAllocator(t *testing.T) {
	Convey("Given the host platform", t, func() {
		platform, err := hs.PopulatePlatform()

		So(platform, ShouldNotBeNil)
		So(err, ShouldBeNil)

		a := &testAllocator{}

		Convey("Given a simple expression with allocator", func() {
			So(hs.SetMiscAllocator(a.alloc, a.free), ShouldBeNil)

			info, err := hs.ExpressionInfo("test", 0)

			So(info, ShouldNotBeNil)
			So(info, ShouldResemble, &hs.ExprInfo{
				MinWidth: 4,
				MaxWidth: 4,
			})
			So(err, ShouldBeNil)

			So(a.memoryUsed, ShouldBeGreaterThanOrEqualTo, 12)

			So(hs.ClearMiscAllocator(), ShouldBeNil)
		})

		Convey("Then create a stream database with allocator", func() {
			So(hs.SetDatabaseAllocator(a.alloc, a.free), ShouldBeNil)

			db, err := hs.Compile("test", 0, hs.StreamMode, platform)

			So(db, ShouldNotBeNil)
			So(err, ShouldBeNil)

			Convey("Get the database size", func() {
				size, err := hs.DatabaseSize(db)

				So(a.memoryUsed, ShouldBeGreaterThanOrEqualTo, size)
				So(err, ShouldBeNil)
			})

			Convey("Then create a scratch with allocator", func() {
				So(hs.SetScratchAllocator(a.alloc, a.free), ShouldBeNil)

				a.memoryUsed = 0

				s, err := hs.AllocScratch(db)

				So(s, ShouldNotBeNil)
				So(err, ShouldBeNil)

				Convey("Get the scratch size", func() {
					size, err := hs.ScratchSize(s)

					So(a.memoryUsed, ShouldBeGreaterThanOrEqualTo, size)
					So(err, ShouldBeNil)
				})

				Convey("Then open a stream", func() {
					So(hs.SetStreamAllocator(a.alloc, a.free), ShouldBeNil)

					a.memoryUsed = 0

					stream, err := hs.OpenStream(db, 0)

					So(stream, ShouldNotBeNil)
					So(err, ShouldBeNil)

					Convey("Get the stream size", func() {
						size, err := hs.StreamSize(db)

						So(a.memoryUsed, ShouldBeGreaterThanOrEqualTo, size)
						So(err, ShouldBeNil)
					})

					h := &hs.MatchRecorder{}

					Convey("Then close stream with allocator", func() {
						a.memoryFreed = nil

						So(hs.CloseStream(stream, s, h.Handle, nil), ShouldBeNil)

						So(hs.ClearStreamAllocator(), ShouldBeNil)
					})
				})

				Convey("Then free scratch with allocator", func() {
					a.memoryFreed = nil

					So(hs.FreeScratch(s), ShouldBeNil)

					So(a.memoryFreed, ShouldResemble, []unsafe.Pointer{unsafe.Pointer(s)})

					So(hs.ClearScratchAllocator(), ShouldBeNil)
				})
			})

			Convey("Then free database with allocator", func() {
				a.memoryFreed = nil

				So(hs.FreeDatabase(db), ShouldBeNil)

				So(a.memoryFreed, ShouldResemble, []unsafe.Pointer{unsafe.Pointer(db)})

				So(hs.ClearDatabaseAllocator(), ShouldBeNil)
			})
		})
	})
}
