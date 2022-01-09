package reports

import (
	"collector/pkg/store"
	"log"
	"strconv"
	"sync"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func Count(start time.Time, end time.Time, shift int) map[string]int {
	reading.Lock()
	caps := make(map[string]int)
	all := store.ExportDataStructAnyTime(start, end)
	for record := range all {
		fr, sum := processRecord(record, shift)
		if sum == 0 {
			continue
		}
		caps[fr] = caps[fr] + sum
	}
	reading.Unlock()
	return caps
}

var (
	prevRecord map[string]store.DataRecord
	reading    sync.Mutex
)

func init() {
	prevRecord = make(map[string]store.DataRecord)
}

func processRecord(record store.DataRecord, shift int) (fraction string, sum int) {
	scaleStr := strconv.Itoa(record.Scale) /* + "_" + record.Fraction*/
	if record.Event == "start" {
		prevRecord[scaleStr] = record
		return "", 0
	}
	if record.Event == "periodic" {
		prev, ok := prevRecord[scaleStr]
		if !ok {
			prevRecord[scaleStr] = record
			log.Println("WW periodic save without start", record)
			return "", 0
		}
		sum = record.Accumulation - prev.Accumulation
		if sum < 0 {
			//todo: chek this
			prevRecord[scaleStr] = record
			//log.Println("WW the previous accumulation", prevRecordAccumulation[scaleStr], "is greater than the current one", record)
			return "", 0
		}
		if (shift == 0) || ((shift == record.Shift) && (shift == prev.Shift)) {
			if record.Fraction == prev.Fraction {
				return scaleStr + "_" + record.Fraction, sum
			}
		}
		return "", 0
	}
	log.Println("WW unexpected record", record)
	return "", 0
}
