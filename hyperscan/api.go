package hyperscan

import (
	"io"

	"github.com/hashicorp/go-multierror"
)

type matchEvent struct {
	id       uint
	from, to uint64
	flags    ScanFlag
}

func (e *matchEvent) Id() uint { return e.id }

func (e *matchEvent) From() uint64 { return e.from }

func (e *matchEvent) To() uint64 { return e.to }

func (e *matchEvent) Flags() ScanFlag { return e.flags }

type matchRecorder struct {
	matched []matchEvent
	err     error
}

func (h *matchRecorder) Handle(id uint, from, to uint64, flags uint, context interface{}) error {
	if len(h.matched) > 0 {
		tail := &h.matched[len(h.matched)-1]

		if tail.id == id && tail.from == from && tail.flags == ScanFlag(flags) && tail.to < to {
			tail.to = to

			return h.err
		}
	}

	h.matched = append(h.matched, matchEvent{id, from, to, ScanFlag(flags)})

	return h.err
}

func Match(pattern string, data []byte) (bool, error) {
	var result *multierror.Error

	if db, err := hsCompile(pattern, 0, BlockMode, nil); err != nil {
		result = multierror.Append(result, err)
	} else {
		if scratch, err := hsAllocScratch(db); err != nil {
			result = multierror.Append(result, err)
		} else {
			h := &matchRecorder{}

			if err = hsScan(db, data, 0, scratch, h.Handle, nil); err != nil {
				result = multierror.Append(result, err)
			}

			if err := hsFreeScratch(scratch); err != nil {
				result = multierror.Append(result, err)
			}

			return h.matched != nil, result.ErrorOrNil()
		}

		if err := hsFreeDatabase(db); err != nil {
			result = multierror.Append(result, err)
		}
	}

	return false, result.ErrorOrNil()
}

func MatchReader(pattern string, reader io.Reader) (bool, error) {
	var result *multierror.Error

	if db, err := hsCompile(pattern, 0, StreamMode, nil); err != nil {
		result = multierror.Append(result, err)
	} else {
		if scratch, err := hsAllocScratch(db); err != nil {
			result = multierror.Append(result, err)
		} else {
			if stream, err := hsOpenStream(db, 0); err != nil {
				result = multierror.Append(result, err)
			} else {
				buf := make([]byte, 4096)

				h := &matchRecorder{}

				for result == nil {
					if read, err := reader.Read(buf); err == io.EOF {
						break
					} else if err != nil {
						result = multierror.Append(result, err)
					} else if err := hsScanStream(stream, buf[:read], 0, scratch, h.Handle, nil); err != nil {
						result = multierror.Append(result, err)
					}
				}

				if err := hsCloseStream(stream, scratch, h.Handle, nil); err != nil {
					result = multierror.Append(result, err)
				}

				return h.matched != nil, result.ErrorOrNil()
			}

			if err := hsFreeScratch(scratch); err != nil {
				result = multierror.Append(result, err)
			}
		}

		if err := hsFreeDatabase(db); err != nil {
			result = multierror.Append(result, err)
		}
	}

	return false, result.ErrorOrNil()
}

func MatchString(pattern string, s string) (matched bool, err error) { return Match(pattern, []byte(s)) }
