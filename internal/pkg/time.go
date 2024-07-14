package pkg

import (
	"strings"
	"time"
)

type Date time.Time

const layout = "2006-01-02"

func (t *Date) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = strings.Trim(s, "\"")

	parsedDob, err := time.Parse(layout, s)
	if err != nil {
		return err
	}

	*t = Date(parsedDob)
	return nil
}
