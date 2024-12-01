package ch

// #include <ch.h>
import "C"

type Scratch *C.ch_scratch_t

func AllocScratch(db Database) (Scratch, error) {
	var scratch *C.ch_scratch_t

	if ret := C.ch_alloc_scratch(db, &scratch); ret != C.CH_SUCCESS {
		return nil, Error(ret)
	}

	return scratch, nil
}

func ReallocScratch(db Database, scratch *Scratch) error {
	if ret := C.ch_alloc_scratch(db, (**C.struct_ch_scratch)(scratch)); ret != C.CH_SUCCESS {
		return Error(ret)
	}

	return nil
}

func CloneScratch(scratch Scratch) (Scratch, error) {
	var clone *C.ch_scratch_t

	if ret := C.ch_clone_scratch(scratch, &clone); ret != C.CH_SUCCESS {
		return nil, Error(ret)
	}

	return clone, nil
}

func ScratchSize(scratch Scratch) (int, error) {
	var size C.size_t

	if ret := C.ch_scratch_size(scratch, &size); ret != C.CH_SUCCESS {
		return 0, Error(ret)
	}

	return int(size), nil
}

func FreeScratch(scratch Scratch) error {
	if ret := C.ch_free_scratch(scratch); ret != C.CH_SUCCESS {
		return Error(ret)
	}

	return nil
}
