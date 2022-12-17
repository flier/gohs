package hyperscan

import (
	"runtime"

	"github.com/flier/gohs/internal/hs"
)

// Scratch is a Hyperscan scratch space.
type Scratch struct {
	s hs.Scratch
}

// NewScratch allocate a "scratch" space for use by Hyperscan.
// This is required for runtime use, and one scratch space per thread,
// or concurrent caller, is required.
func NewScratch(db Database) (*Scratch, error) {
	s, err := hs.AllocScratch(db.(database).c())
	if err != nil {
		return nil, err //nolint: wrapcheck
	}

	return &Scratch{s}, nil
}

// NewManagedScratch is a wrapper for NewScratch that sets
// a finalizer on the Scratch instance so that memory is freed
// once the object is no longer in use.
func NewManagedScratch(db Database) (*Scratch, error) {
	s, err := NewScratch(db)
	if err != nil {
		return nil, err
	}

	runtime.SetFinalizer(s, func(scratch *Scratch) {
		_ = scratch.Free()
	})
	return s, nil
}

// Size provides the size of the given scratch space.
func (s *Scratch) Size() (int, error) { return hs.ScratchSize(s.s) } //nolint: wrapcheck

// Realloc reallocate the scratch for another database.
func (s *Scratch) Realloc(db Database) error {
	r, _ := db.(database)

	return hs.ReallocScratch(r.c(), &s.s) //nolint: wrapcheck
}

// Clone allocate a scratch space that is a clone of an existing scratch space.
func (s *Scratch) Clone() (*Scratch, error) {
	cloned, err := hs.CloneScratch(s.s)
	if err != nil {
		return nil, err //nolint: wrapcheck
	}

	return &Scratch{cloned}, nil
}

// Free a scratch block previously allocated.
func (s *Scratch) Free() error { return hs.FreeScratch(s.s) } //nolint: wrapcheck
