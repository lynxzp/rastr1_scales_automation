package ucma

import (
	"collector/pkg/shift"
	"collector/pkg/store"
	"time"
)

func (ucma *Ucma) startSave() {
	store.SaveEvent(int(ucma.Id), int(ucma.DataAccumValue), "start", shift.GetCurrentShift(), ucma.Fraction)
	ucma.time = time.Now()
}

func (ucma *Ucma) periodicSave() {
	if ucma.DataAccumValue == 0 {
		return
	}
	n := time.Now()
	if ucma.time.Add(time.Second).Before(n) {
		ucma.time = n
		store.PeriodicSave(int(ucma.Id), int(ucma.DataAccumValue), "periodic", shift.GetCurrentShift(), ucma.Fraction)
	}
}
