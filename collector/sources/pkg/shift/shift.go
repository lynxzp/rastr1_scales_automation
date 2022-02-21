package shift

import (
	"collector/pkg/config"
	"collector/pkg/time"
)

//var (
//	t12, t23, t31 time.Time
//)

//func init() {
//	var err error
//	t12, err = time.Parse("15:04", config.Cfg.Shift12ChangeTime)
//	if err != nil {
//		log.Fatalln("EE can't parse time changing shift 1 and 2", err)
//	}
//	t23, err = time.Parse("15:04", config.Cfg.Shift23ChangeTime)
//	if err != nil {
//		log.Fatalln("EE can't parse time changing shift 2 and 3", err)
//	}
//	t31, err = time.Parse("15:04", config.Cfg.Shift31ChangeTime)
//	if err != nil {
//		log.Fatalln("EE can't parse time changing shift 3 and 1", err)
//	}

//}

func GetCurrentShift(scaleNum int8) int {
	//t := time.Now()
	//if t.Before(t12) {
	//	return 1
	//}
	//if t.Before(t23) {
	//	return 2
	//}
	//if t.Before(t31) {
	//	return 3
	//}
	//return 0
	t := time.Now()
	for i := range config.Cfg.Shifts[scaleNum] {
		if t.AfterOrEqual(config.Cfg.Shifts[scaleNum][i].Start) && t.BeforeOrEqual(config.Cfg.Shifts[scaleNum][i].Finish) {
			return config.Cfg.Shifts[scaleNum][i].Number
		}
	}
	return 0
}
