package hs

// #include <hs.h>
import "C"

type Scratch *C.hs_scratch_t

func AllocScratch(db Database) (Scratch, error) {
	var scratch *C.hs_scratch_t

	if ret := C.hs_alloc_scratch(db, &scratch); ret != C.HS_SUCCESS {
		return nil, Error(ret)
	}

	return scratch, nil
}

func ReallocScratch(db Database, scratch *Scratch) error {
	if ret := C.hs_alloc_scratch(db, (**C.struct_hs_scratch)(scratch)); ret != C.HS_SUCCESS {
		return Error(ret)
	}

	return nil
}

func CloneScratch(scratch Scratch) (Scratch, error) {
	var clone *C.hs_scratch_t

	if ret := C.hs_clone_scratch(scratch, &clone); ret != C.HS_SUCCESS {
		return nil, Error(ret)
	}

	return clone, nil
}

func ScratchSize(scratch Scratch) (int, error) {
	var size C.size_t

	if ret := C.hs_scratch_size(scratch, &size); ret != C.HS_SUCCESS {
		return 0, Error(ret)
	}

	return int(size), nil
}

func FreeScratch(scratch Scratch) error {
	if ret := C.hs_free_scratch(scratch); ret != C.HS_SUCCESS {
		return Error(ret)
	}

	return nil
}
