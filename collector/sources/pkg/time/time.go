package time

import (
	"log"
	"strconv"
)

type Time struct {
	hours   int
	minutes int
}

func (t *Time) UnmarshalJSON(b []byte) error {
	*t = Time{}
	var err error
	t.hours, err = strconv.Atoi(string(b[1:3]))
	if err != nil {
		log.Fatalln("EE Can't parse time", string(b))
	}
	t.minutes, err = strconv.Atoi(string(b[4:6]))
	if err != nil {
		log.Fatalln("EE Can't parse time", string(b))
	}
	return nil
}

func (t1 *Time) BeforeOrEqual(t2 Time) bool {
	if t2.hours > t1.hours {
		return true
	}
	if t1.hours > t2.hours {
		return false
	}
	if t1.minutes <= t2.minutes {
		return true
	}
	return false
}

func (t1 *Time) AfterOrEqual(t2 Time) bool {
	if t1.hours > t2.hours {
		return true
	}
	if t1.hours < t2.hours {
		return false
	}
	if t1.minutes >= t2.minutes {
		return true
	}
	return false
}
