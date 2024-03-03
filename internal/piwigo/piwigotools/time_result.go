package piwigotools

import (
	"fmt"
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
	switch s := c.String(); s {
	case "":
		return []byte("null"), nil
	default:
		return []byte(`"` + s + `"`), nil
	}
}

func (c TimeResult) String() string {
	t := c.toTime()
	if t.IsZero() {
		return ""
	} else {
		return t.Format("2006-01-02 15:04:05")
	}
}

func (c TimeResult) toTime() time.Time {
	return time.Time(c)
}

func (c TimeResult) AgeAt(createdAt *TimeResult) string {
	var year, month, day, hour, minutes, seconds int
	a := c.toTime()
	if a.IsZero() {
		return ""
	}
	b := createdAt.toTime()
	if b.IsZero() {
		return ""
	}

	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = y2 - y1
	month = int(M2 - M1)
	day = d2 - d1
	hour = h2 - h1
	minutes = m2 - m1
	seconds = s2 - s1

	// Normalize negative values
	if seconds < 0 {
		seconds += 60
		minutes--
	}
	if minutes < 0 {
		minutes += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	switch {
	case year == 1:
		return fmt.Sprintf("%d months old", 12+month)
	case year > 1:
		return fmt.Sprintf("%d years old", year)
	case month == 1:
		return fmt.Sprintf("%d month old", month)
	case month > 1:
		return fmt.Sprintf("%d months old", month)
	case day > 1:
		return fmt.Sprintf("%d days old", day)
	default:
		return fmt.Sprintf("%d day, %d hour", day, hour)
	}
}
