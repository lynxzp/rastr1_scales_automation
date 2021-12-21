package time

import (
	"log"
	"testing"
)

type times struct {
	t1 Time
	t2 Time
}

var set []times

func push(a string, b string) {
	t := times{}
	t.t1.UnmarshalJSON([]byte(`"` + a + `"`))
	t.t2.UnmarshalJSON([]byte(`"` + b + `"`))
	set = append(set, t)
}

func TestTime_BeforeOrEqual(t *testing.T) {
	push("01:00", "02:00")
	push("01:00", "03:00")
	push("01:00", "10:00")
	push("01:00", "01:00")
	push("01:00", "01:02")
	push("00:00", "00:00")
	push("00:20", "00:20")
	push("01:00", "10:59")
	for i := range set {
		log.Println(set[i].t1, set[i].t2)
		if set[i].t1.BeforeOrEqual(set[i].t2) == false {
			t.Error(set[i].t1, "is before or equal", set[i].t2)
		}
	}
}
