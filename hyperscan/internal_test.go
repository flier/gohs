package hyperscan

import (
	"errors"
	"regexp"
	"testing"
	"unsafe"

	. "github.com/smartystreets/goconvey/convey"
)

func TestVersion(t *testing.T) {
	Convey("Given a HyperScan version", t, func() {
		ver := hsVersion()

		So(ver, ShouldNotBeEmpty)

		matched, err := regexp.MatchString(`^\d\.\d\.\d.*`, ver)

		So(err, ShouldBeNil)
		So(matched, ShouldBeTrue)
	})
}

func TestModeFlag(t *testing.T) {
	Convey("Give a mode", t, func() {
		So(BlockMode.String(), ShouldEqual, "BLOCK")
		So(StreamMode.String(), ShouldEqual, "STREAM")
		So(VectoredMode.String(), ShouldEqual, "VECTORED")

		Convey("When combile mode with flags", func() {
			mode := StreamMode | SomHorizonLargeMode

			So(mode.String(), ShouldEqual, "STREAM")
		})

		Convey("When parse unknown mode", func() {
			m, err := ParseModeFlag("test")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "Unknown Mode: test")
			So(m, ShouldEqual, BlockMode)
		})
	})
}

func TestCompileFlag(t *testing.T) {
	Convey("Given a compile flags", t, func() {
		flags := Caseless | DotAll | MultiLine | SingleMatch | AllowEmpty | Utf8Mode | UnicodeProperty | PrefilterMode

		So(flags.String(), ShouldEqual, "efimopsu")

		Convey("When parse valid flags", func() {
			f, err := ParseCompileFlag("ifemopus")

			So(f, ShouldEqual, flags)
			So(err, ShouldBeNil)
		})

		Convey("When parse invalid flags", func() {
			f, err := ParseCompileFlag("abc")

			So(f, ShouldEqual, 0)
			So(err, ShouldNotBeNil)
		})
	})
}

type testAllocator struct {
	memoryUsed  int
	memoryFreed []unsafe.Pointer
}

func (a *testAllocator) alloc(size uint) unsafe.Pointer {
	a.memoryUsed += int(size)

	return hsDefaultAlloc(size)
}

func (a *testAllocator) free(ptr unsafe.Pointer) {
	a.memoryFreed = append(a.memoryFreed, ptr)

	hsDefaultFree(ptr)
}

func TestAllocator(t *testing.T) {
	Convey("Given the host platform", t, func() {
		platform, err := hsPopulatePlatform()

		So(platform, ShouldNotBeNil)
		So(err, ShouldBeNil)

		a := &testAllocator{}

		Convey("Given a simple expression with allocator", func() {
			So(hsSetMiscAllocator(a.alloc, a.free), ShouldBeNil)

			info, err := hsExpressionInfo("test", 0)

			So(info, ShouldNotBeNil)
			So(info, ShouldResemble, &ExprInfo{
				MinWidth: 4,
				MaxWidth: 4,
			})
			So(err, ShouldBeNil)

			So(a.memoryUsed, ShouldBeGreaterThanOrEqualTo, 12)

			So(hsClearMiscAllocator(), ShouldBeNil)
		})

		Convey("Then create a stream database with allocator", func() {
			So(hsSetDatabaseAllocator(a.alloc, a.free), ShouldBeNil)

			db, err := hsCompile("test", 0, StreamMode, platform)

			So(db, ShouldNotBeNil)
			So(err, ShouldBeNil)

			Convey("Get the database size", func() {
				size, err := hsDatabaseSize(db)

				So(a.memoryUsed, ShouldBeGreaterThanOrEqualTo, size)
				So(err, ShouldBeNil)
			})

			Convey("Then create a scratch with allocator", func() {
				So(hsSetScratchAllocator(a.alloc, a.free), ShouldBeNil)

				a.memoryUsed = 0

				s, err := hsAllocScratch(db)

				So(s, ShouldNotBeNil)
				So(err, ShouldBeNil)

				Convey("Get the scratch size", func() {
					size, err := hsScratchSize(s)

					So(a.memoryUsed, ShouldBeGreaterThanOrEqualTo, size)
					So(err, ShouldBeNil)
				})

				Convey("Then open a stream", func() {
					So(hsSetStreamAllocator(a.alloc, a.free), ShouldBeNil)

					a.memoryUsed = 0

					stream, err := hsOpenStream(db, 0)

					So(stream, ShouldNotBeNil)
					So(err, ShouldBeNil)

					Convey("Get the stream size", func() {
						size, err := hsStreamSize(db)

						So(a.memoryUsed, ShouldBeGreaterThanOrEqualTo, size)
						So(err, ShouldBeNil)
					})

					h := &matchRecorder{}

					Convey("Then close stream with allocator", func() {
						a.memoryFreed = nil

						So(hsCloseStream(stream, s, h.Handle, nil), ShouldBeNil)

						So(hsClearStreamAllocator(), ShouldBeNil)
					})
				})

				Convey("Then free scratch with allocator", func() {
					a.memoryFreed = nil

					So(hsFreeScratch(s), ShouldBeNil)

					So(a.memoryFreed, ShouldResemble, []unsafe.Pointer{unsafe.Pointer(s)})

					So(hsClearScratchAllocator(), ShouldBeNil)
				})
			})

			Convey("Then free database with allocator", func() {
				a.memoryFreed = nil

				So(hsFreeDatabase(db), ShouldBeNil)

				So(a.memoryFreed, ShouldResemble, []unsafe.Pointer{unsafe.Pointer(db)})

				So(hsClearDatabaseAllocator(), ShouldBeNil)
			})
		})
	})
}

func TestDatabase(t *testing.T) {
	Convey("Given a stream database", t, func() {
		platform, err := hsPopulatePlatform()

		So(platform, ShouldNotBeNil)
		So(err, ShouldBeNil)

		db, err := hsCompile("test", 0, StreamMode, platform)

		So(db, ShouldNotBeNil)
		So(err, ShouldBeNil)

		Convey("Get the database info", func() {
			info, err := hsDatabaseInfo(db)

			So(regexInfo.MatchString(info), ShouldBeTrue)
			So(err, ShouldBeNil)
		})

		Convey("Get the database size", func() {
			size, err := hsDatabaseSize(db)

			So(size, ShouldBeGreaterThan, 800)
			So(err, ShouldBeNil)
		})

		Convey("Get the stream size", func() {
			size, err := hsStreamSize(db)

			So(size, ShouldBeGreaterThan, 20)
			So(err, ShouldBeNil)
		})

		Convey("Get the stream size from a block database", func() {
			db, err := hsCompile("test", 0, BlockMode, platform)

			So(db, ShouldNotBeNil)
			So(err, ShouldBeNil)

			size, err := hsStreamSize(db)

			So(size, ShouldEqual, 0)
			So(err, ShouldEqual, ErrDatabaseModeError)
		})

		Convey("When serialize database", func() {
			data, err := hsSerializeDatabase(db)

			So(data, ShouldNotBeNil)
			So(len(data), ShouldBeGreaterThan, 800)
			So(err, ShouldBeNil)

			Convey("Get the database info", func() {
				info, err := hsSerializedDatabaseInfo(data)

				So(regexInfo.MatchString(info), ShouldBeTrue)
				So(err, ShouldBeNil)
			})

			Convey("Get the database size", func() {
				size, err := hsSerializedDatabaseSize(data)

				So(size, ShouldBeGreaterThan, 800)
				So(err, ShouldBeNil)
			})

			Convey("Then deserialize database", func() {
				db, err := hsDeserializeDatabase(data)

				So(db, ShouldNotBeNil)
				So(err, ShouldBeNil)

				Convey("Get the database info", func() {
					info, err := hsDatabaseInfo(db)

					So(regexInfo.MatchString(info), ShouldBeTrue)
					So(err, ShouldBeNil)
				})
			})

			Convey("Then deserialize database to memory", func() {
				buf := make([]byte, 1000)
				db := hsDatabase(unsafe.Pointer(&buf[0]))

				So(hsDeserializeDatabaseAt(data, db), ShouldBeNil)

				Convey("Get the database info", func() {
					info, err := hsDatabaseInfo(db)

					So(regexInfo.MatchString(info), ShouldBeTrue)
					So(err, ShouldBeNil)
				})
			})
		})

		So(hsFreeDatabase(db), ShouldBeNil)
	})
}

func TestCompileAPI(t *testing.T) {
	Convey("Given a host platform", t, func() {
		platform, err := hsPopulatePlatform()

		So(platform, ShouldNotBeNil)
		So(err, ShouldBeNil)

		Convey("When Compile a unsupported expression", func() {
			Convey("Then compile as stream", func() {
				db, err := hsCompile(`\R`, 0, StreamMode, platform)

				So(db, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, `\R at index 0 not supported.`)
			})

			Convey("Then compile as vector", func() {
				db, err := hsCompileMulti([]*Pattern{
					NewPattern(`\R`, Caseless),
				}, BlockMode, platform)

				So(db, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, `\R at index 0 not supported.`)
			})

			Convey("Then compile as extended vector", func() {
				db, err := hsCompileExtMulti([]string{`\R`}, []CompileFlag{Caseless}, []uint{1}, []ExprExt{{Flags: MinOffset, MinOffset: 10}}, VectoredMode, platform)

				So(db, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, `\R at index 0 not supported.`)
			})
		})

		Convey("Compile an empty expression", func() {
			db, err := hsCompile("", 0, StreamMode, platform)

			So(db, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "Pattern matches empty buffer; use HS_FLAG_ALLOWEMPTY to enable support.")

			So(hsFreeDatabase(db), ShouldBeNil)
		})

		Convey("Compile multi expressions", func() {
			db, err := hsCompileMulti([]*Pattern{
				NewPattern(`^\w+`, 0),
				NewPattern(`\d+`, 0),
				NewPattern(`\s+`, 0),
			}, StreamMode, platform)

			So(db, ShouldNotBeNil)
			So(err, ShouldBeNil)

			Convey("Get the database info", func() {
				info, err := hsDatabaseInfo(db)

				So(regexInfo.MatchString(info), ShouldBeTrue)
				So(err, ShouldBeNil)
			})

			So(hsFreeDatabase(db), ShouldBeNil)
		})

		Convey("Compile multi expressions with extension", func() {
			exts := []ExprExt{
				{Flags: MinOffset, MinOffset: 10},
				{Flags: MaxOffset, MaxOffset: 10},
				{Flags: MinLength, MinLength: 10},
				{Flags: EditDistance, EditDistance: 10},
			}
			db, err := hsCompileExtMulti([]string{`^\w+`, `\d+`, `\s+`}, nil, []uint{1, 2, 3}, exts, StreamMode, platform)

			So(db, ShouldNotBeNil)
			So(err, ShouldBeNil)

			Convey("Get the database info", func() {
				info, err := hsDatabaseInfo(db)

				So(regexInfo.MatchString(info), ShouldBeTrue)
				So(err, ShouldBeNil)
			})

			So(hsFreeDatabase(db), ShouldBeNil)
		})
	})
}

func TestExpression(t *testing.T) {
	Convey("Given a simple expression", t, func() {
		info, err := hsExpressionInfo("test", 0)

		So(info, ShouldNotBeNil)
		So(info, ShouldResemble, &ExprInfo{
			MinWidth: 4,
			MaxWidth: 4,
		})
		So(err, ShouldBeNil)
	})

	Convey("Given a credit card expression", t, func() {
		info, err := hsExpressionInfo(CreditCard, 0)

		So(info, ShouldNotBeNil)
		So(info, ShouldResemble, &ExprInfo{
			MinWidth: 13,
			MaxWidth: 16,
		})
		So(err, ShouldBeNil)
	})

	Convey("Given a expression match eod", t, func() {
		info, err := hsExpressionInfo("test$", 0)

		So(info, ShouldNotBeNil)
		So(info, ShouldResemble, &ExprInfo{
			MinWidth:        4,
			MaxWidth:        4,
			ReturnUnordered: true,
			AtEndOfData:     true,
			OnlyAtEndOfData: true,
		})
		So(err, ShouldBeNil)
	})
}

func TestScratch(t *testing.T) {
	Convey("Given a block database", t, func() {
		platform, err := hsPopulatePlatform()

		So(platform, ShouldNotBeNil)
		So(err, ShouldBeNil)

		db, err := hsCompile("test", 0, BlockMode, platform)

		So(db, ShouldNotBeNil)
		So(err, ShouldBeNil)

		Convey("Allocate a scratch", func() {
			s, err := hsAllocScratch(db)

			So(s, ShouldNotBeNil)
			So(err, ShouldBeNil)

			Convey("Get the scratch size", func() {
				size, err := hsScratchSize(s)

				So(size, ShouldBeGreaterThan, 1024)
				So(size, ShouldBeLessThan, 4096)
				So(err, ShouldBeNil)

				Convey("Clone the scratch", func() {
					s2, err := hsCloneScratch(s)

					So(s2, ShouldNotBeNil)
					So(err, ShouldBeNil)

					Convey("Cloned scrash should have same size", func() {
						size2, err := hsScratchSize(s2)

						So(size2, ShouldEqual, size)
						So(err, ShouldBeNil)
					})

					So(hsFreeScratch(s2), ShouldBeNil)
				})

				Convey("Reallocate the scratch with another database", func() {
					db2, err := hsCompile(EmailAddress, 0, BlockMode, platform)

					So(db, ShouldNotBeNil)
					So(err, ShouldBeNil)

					So(hsReallocScratch(db2, &s), ShouldBeNil)

					size2, err := hsScratchSize(s)

					So(size2, ShouldBeGreaterThan, size)
					So(err, ShouldBeNil)

					So(hsFreeDatabase(db2), ShouldBeNil)
				})
			})

			So(hsFreeScratch(s), ShouldBeNil)
		})

		So(hsFreeDatabase(db), ShouldBeNil)
	})
}

func TestBlockScan(t *testing.T) {
	Convey("Given a block database", t, func() {
		platform, err := hsPopulatePlatform()

		So(platform, ShouldNotBeNil)
		So(err, ShouldBeNil)

		db, err := hsCompile("test", 0, BlockMode, platform)

		So(db, ShouldNotBeNil)
		So(err, ShouldBeNil)

		s, err := hsAllocScratch(db)

		So(s, ShouldNotBeNil)
		So(err, ShouldBeNil)

		h := &matchRecorder{}

		Convey("Scan block with pattern", func() {
			So(hsScan(db, []byte("abctestdef"), 0, s, h.Handle, nil), ShouldBeNil)
			So(h.matched, ShouldResemble, []matchEvent{{0, 0, 7, 0}})
		})

		Convey("Scan block without pattern", func() {
			So(hsScan(db, []byte("abcdef"), 0, s, h.Handle, nil), ShouldBeNil)
			So(h.matched, ShouldBeEmpty)
		})

		Convey("Scan block with multi pattern", func() {
			So(hsScan(db, []byte("abctestdeftest"), 0, s, h.Handle, nil), ShouldBeNil)
			So(h.matched, ShouldResemble, []matchEvent{{0, 0, 14, 0}})
		})

		Convey("Scan block with multi pattern but terminated", func() {
			h.err = errors.New("terminated")

			So(hsScan(db, []byte("abctestdeftest"), 0, s, h.Handle, nil), ShouldEqual, ErrScanTerminated)
			So(h.matched, ShouldResemble, []matchEvent{{0, 0, 7, 0}})
		})

		So(hsFreeScratch(s), ShouldBeNil)
		So(hsFreeDatabase(db), ShouldBeNil)
	})
}

func TestVectorScan(t *testing.T) {
	Convey("Given a block database", t, func() {
		platform, err := hsPopulatePlatform()

		So(platform, ShouldNotBeNil)
		So(err, ShouldBeNil)

		db, err := hsCompile("test", 0, VectoredMode, platform)

		So(db, ShouldNotBeNil)
		So(err, ShouldBeNil)

		s, err := hsAllocScratch(db)

		So(s, ShouldNotBeNil)
		So(err, ShouldBeNil)

		h := &matchRecorder{}

		Convey("Scan multi block with pattern", func() {
			So(hsScanVector(db, [][]byte{[]byte("abctestdef"), []byte("abcdef")}, 0, s, h.Handle, nil), ShouldBeNil)
			So(h.matched, ShouldResemble, []matchEvent{{0, 0, 7, 0}})
		})

		Convey("Scan multi block without pattern", func() {
			So(hsScanVector(db, [][]byte{[]byte("123456"), []byte("abcdef")}, 0, s, h.Handle, nil), ShouldBeNil)
			So(h.matched, ShouldBeEmpty)
		})

		Convey("Scan multi block with multi pattern", func() {
			So(hsScanVector(db, [][]byte{[]byte("abctestdef"), []byte("123test456")}, 0, s, h.Handle, nil), ShouldBeNil)
			So(h.matched, ShouldResemble, []matchEvent{{0, 0, 17, 0}})
		})

		Convey("Scan multi block with multi pattern but terminated", func() {
			h.err = errors.New("terminated")

			So(hsScanVector(db, [][]byte{[]byte("abctestdef"), []byte("123test456")}, 0, s, h.Handle, nil), ShouldEqual, ErrScanTerminated)
			So(h.matched, ShouldResemble, []matchEvent{{0, 0, 7, 0}})
		})

		So(hsFreeScratch(s), ShouldBeNil)
	})
}

func TestStreamScan(t *testing.T) {
	Convey("Given a stream database", t, func() {
		platform, err := hsPopulatePlatform()

		So(platform, ShouldNotBeNil)
		So(err, ShouldBeNil)

		db, err := hsCompile("test", 0, StreamMode, platform)

		So(db, ShouldNotBeNil)
		So(err, ShouldBeNil)

		s, err := hsAllocScratch(db)

		So(s, ShouldNotBeNil)
		So(err, ShouldBeNil)

		Convey("Then open a stream", func() {
			stream, err := hsOpenStream(db, 0)

			So(stream, ShouldNotBeNil)
			So(err, ShouldBeNil)

			h := &matchRecorder{}

			Convey("Then scan a simple stream with first part", func() {
				So(hsScanStream(stream, []byte("abcte"), 0, s, h.Handle, nil), ShouldBeNil)
				So(h.matched, ShouldBeNil)

				Convey("When scan second part, should be matched", func() {
					So(hsScanStream(stream, []byte("stdef"), 0, s, h.Handle, nil), ShouldBeNil)
					So(h.matched, ShouldResemble, []matchEvent{{0, 0, 7, 0}})
				})

				Convey("Then copy the stream", func() {
					stream2, err := hsCopyStream(stream)

					So(stream2, ShouldNotBeNil)
					So(err, ShouldBeNil)

					Convey("When copied stream2 scan the second part, should be matched", func() {
						So(hsScanStream(stream2, []byte("stdef"), 0, s, h.Handle, nil), ShouldBeNil)
						So(h.matched, ShouldResemble, []matchEvent{{0, 0, 7, 0}})

						Convey("When copied stream2 scan the second part again, should not be matched", func() {
							h.matched = nil
							So(hsScanStream(stream2, []byte("stdef"), 0, s, h.Handle, nil), ShouldBeNil)
							So(h.matched, ShouldBeNil)

							Convey("When copy and reset stream2", func() {
								So(hsResetAndCopyStream(stream2, stream, s, h.Handle, nil), ShouldBeNil)

								Convey("When copied and reset stream2 scan the second part again, should be matched", func() {
									h.matched = nil
									So(hsScanStream(stream2, []byte("stdef"), 0, s, h.Handle, nil), ShouldBeNil)
									So(h.matched, ShouldResemble, []matchEvent{{0, 0, 7, 0}})
								})
							})
						})
					})

					So(hsCloseStream(stream2, s, h.Handle, nil), ShouldBeNil)
				})

				Convey("Then reset the stream", func() {
					So(hsResetStream(stream, 0, s, h.Handle, nil), ShouldBeNil)

					Convey("When scan the second part, should not be matched", func() {
						So(hsScanStream(stream, []byte("stdef"), 0, s, h.Handle, nil), ShouldBeNil)
						So(h.matched, ShouldBeNil)
					})
				})
			})

			So(hsCloseStream(stream, s, h.Handle, nil), ShouldBeNil)
		})

		So(hsFreeScratch(s), ShouldBeNil)
	})
}
