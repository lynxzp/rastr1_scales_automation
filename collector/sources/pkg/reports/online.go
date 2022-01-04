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
	all := store.ExportDataStructAnyTime(start, end, shift)
	for record := range all {
		fr, sum := processRecord(record)
		if sum == 0 {
			continue
		}
		caps[fr] = caps[fr] + sum
	}
	reading.Unlock()
	return caps
}

var (
	prevRecordAccumulation map[string]int
	reading                sync.Mutex
)

func init() {
	prevRecordAccumulation = make(map[string]int)
}

func processRecord(record store.DataRecord) (string, int) {
	scaleStr := strconv.Itoa(record.Scale) + "_" + record.Fraction
	if record.Event == "start" {
		prevRecordAccumulation[scaleStr] = record.Accumulation
		return "", 0
	}
	if record.Event == "periodic" {
		prev, ok := prevRecordAccumulation[scaleStr]
		if !ok {
			prevRecordAccumulation[scaleStr] = record.Accumulation
			log.Println("WW periodic save without start", record)
			return "", 0
		}
		sum := record.Accumulation - prev
		if sum < 0 {
			//todo: chek this
			prevRecordAccumulation[scaleStr] = record.Accumulation
			//log.Println("WW the previous accumulation", prevRecordAccumulation[scaleStr], "is greater than the current one", record)
			return "", 0
		}
		return scaleStr, sum
	}
	log.Println("WW unexpected record", record)
	return "", 0
}
