package types

import (
	"strings"
	"time"
)

type FormDate struct {
	time.Time
}

func (t FormDate) MarshalJSON() ([]byte, error) {
	return []byte(t.Format(time.DateOnly)), nil
}

func (t *FormDate) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")

	date, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return err
	}
	t.Time = date
	return
}
