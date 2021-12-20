package piwigo

import (
	"strings"
	"time"
)

type TimeResult time.Time

func (c *TimeResult) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`) //get rid of "
	if value == "" || value == "null" {
		return nil
	}

	t, err := time.Parse("2006-01-02 15:04:05", value) //parse time
	if err != nil {
		return err
	}
	*c = TimeResult(t) //set result using the pointer
	return nil
}

func (c TimeResult) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(c).Format(time.RFC3339) + `"`), nil
}
