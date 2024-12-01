package chimera

import (
	"fmt"

	"github.com/flier/gohs/internal/ch"
)

// Match reports whether the byte slice b contains any match of the regular expression pattern.
func Match(pattern string, data []byte) (bool, error) {
	p, err := ParsePattern(pattern)
	if err != nil {
		return false, fmt.Errorf("parse pattern, %w", err)
	}

	p.Flags |= SingleMatch

	db, err := NewBlockDatabase(p)
	if err != nil {
		return false, fmt.Errorf("create block database, %w", err)
	}
	defer db.Close()

	s, err := NewScratch(db)
	if err != nil {
		return false, fmt.Errorf("create scratch, %w", err)
	}

	defer func() {
		_ = s.Free()
	}()

	h := &ch.MatchRecorder{}

	if err = db.Scan(data, s, h, nil); err != nil {
		return false, err //nolint: wrapcheck
	}

	return h.Matched(), h.Err
}

// MatchString reports whether the string s contains any match of the regular expression pattern.
func MatchString(pattern, s string) (matched bool, err error) {
	return Match(pattern, []byte(s))
}
