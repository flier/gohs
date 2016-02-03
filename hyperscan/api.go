package hyperscan

import (
	"io"

	"github.com/hashicorp/go-multierror"
)

type matchEvent struct {
	id       uint
	from, to uint64
}

type matchHandler struct {
	matched []matchEvent
	err     error
}

func (h *matchHandler) handle(id uint, from, to uint64, flags uint, context interface{}) error {
	h.matched = append(h.matched, matchEvent{id, from, to})

	return h.err
}

func Match(pattern string, data []byte) (bool, error) {
	var result *multierror.Error

	if db, err := hsCompile(pattern, 0, Block, nil); err != nil {
		result = multierror.Append(result, err)
	} else {
		if scratch, err := hsAllocScratch(db); err != nil {
			result = multierror.Append(result, err)
		} else {
			h := &matchHandler{}

			if err = hsScan(db, data, 0, scratch, h, nil); err != nil {
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

	if db, err := hsCompile(pattern, 0, Stream, nil); err != nil {
		result = multierror.Append(result, err)
	} else {
		if scratch, err := hsAllocScratch(db); err != nil {
			result = multierror.Append(result, err)
		} else {
			if stream, err := hsOpenStream(db, 0); err != nil {
				result = multierror.Append(result, err)
			} else {
				buf := make([]byte, 4096)

				h := &matchHandler{}

				for result == nil {
					if read, err := reader.Read(buf); err == io.EOF {
						break
					} else if err != nil {
						result = multierror.Append(result, err)
					} else if err := hsScanStream(stream, buf[:read], 0, scratch, h, nil); err != nil {
						result = multierror.Append(result, err)
					}
				}

				if err := hsCloseStream(stream, scratch, h, nil); err != nil {
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
