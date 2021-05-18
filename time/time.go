package time

import (
	"time"

	"github.com/pkg/errors"
)

// ToString convert time to string
//
func ToString(t time.Time) string {
	return t.Format(time.RFC3339)
}

// FromString convert string to time
//
func FromString(s string) (time.Time, error) {
	t, err := time.Parse(
		time.RFC3339,
		s)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "from string "+s)
	}
	return t, nil
}
