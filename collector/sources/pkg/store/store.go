package store

import (
	"collector/pkg/config"
	scales_naming_grouping "collector/pkg/scales-naming-grouping"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"runtime"
	"time"
)

const (
	file string = "db.sqlite"
)

var db *sql.DB

type DataRecord struct {
	Scale        int
	Accumulation int
	Event        string
	Shift        int
	Fraction     string
	Datetime     string
}

func init() {
	var err error
	db, err = sql.Open("sqlite3", file)
	if err != nil {
		log.Fatalln("Can't open db ", file, " with error ", err)
	}
	runtime.SetFinalizer(db, func(db *sql.DB) {
		db.Close()
		log.Println("Database save closed")
	})

	sqlStmt := `CREATE TABLE if not exists scales (id INTEGER, ip TEXT, rs485addr INTEGER, data_perf_addr INTEGER, fraction TEXT)`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("EE %q: %s\n", err, sqlStmt)
		return
	}

	sqlStmt = `CREATE TABLE if not exists data (scale INTEGER, accumulation INTEGER, event TEXT, shift INTEGER, fraction TEXT, datetime INTEGER)`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("EE %q: %s\n", err, sqlStmt)
		return
	}

	sqlStmt = `CREATE INDEX IF NOT EXISTS index_datetime ON data(datetime)`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("EE %q: %s\n", err, sqlStmt)
		return
	}
}

func SaveScale(id int, dataPerfAddr int, ip string, rs485addr int, fraction string) {
	ClearScale(id)
	smt := "INSERT INTO scales (id, ip, rs485addr, data_perf_addr, fraction) VALUES(?, ?, ?, ?, ?)"
	res, err := db.Exec(smt, id, ip, rs485addr, dataPerfAddr, fraction)
	if err != nil {
		log.Println("WW can't save scales:", id, ip, rs485addr, dataPerfAddr, fraction, "with err:", err)
		return
	}
	affected, err := res.RowsAffected()
	if (err != nil) || (affected != 1) {
		log.Println("WW problem saving scales, err:", err, "rows affected:", affected)
		return
	}
}

func ClearScale(id int) {
	smt := "DELETE FROM scales WHERE id = ?"
	res, err := db.Exec(smt, id)
	if err != nil {
		log.Println("WW can't clear scales:", id, "with err:", err)
		return
	}
	affected, err := res.RowsAffected()
	if (err != nil) || (affected > 1) {
		log.Println("WW problem clearing scales, err:", err, "rows affected:", affected)
		return
	}
}

type Scale struct {
	Id           int
	DataPerfAddr int
	Ip           string
	Rs485addr    int
	Fraction     string
}

func ReadScales() ([]Scale, error) {
	scales := make([]Scale, config.ScalesNums)
	rows, err := db.Query(`SELECT * FROM scales`)
	if err != nil {
		return nil, errors.New("can't perform sql request")
	}
	defer rows.Close()

	var id, rs485addr, dataPerfAddr int
	var ip, fraction string
	for rows.Next() {
		err = rows.Scan(&id, &ip, &rs485addr, &dataPerfAddr, &fraction)
		if err != nil {
			continue
		}
		scales[id].Rs485addr = rs485addr
		scales[id].Ip = ip
		scales[id].DataPerfAddr = dataPerfAddr
		scales[id].Fraction = fraction
	}
	return scales, nil
}

func SaveScaleFraction(scale int, fraction string) {
	smt := "UPDATE scales SET fraction = ? WHERE id = ?"
	res, err := db.Exec(smt, fraction, scale)
	if err != nil {
		log.Println("WW can't save scale fraction:", scale, fraction, "with err:", err)
		return
	}
	affected, err := res.RowsAffected()
	if (err != nil) || (affected != 1) {
		log.Println("WW problem saving scale fraction, err:", err, "rows affected:", affected)
		return
	}
}

func SaveEvent(scale int, accumulation int, event string, shift int, fraction string) {
	smt := "INSERT INTO data (scale, accumulation, event, shift, fraction, datetime) VALUES (?, ?, ?, ?, ?, datetime('now','localtime'))"
	res, err := db.Exec(smt, scale, accumulation, event, shift, fraction)
	if err != nil {
		log.Println("WW can't save scales:", scale, accumulation, event, shift, fraction, "with err:", err)
		return
	}
	affected, err := res.RowsAffected()
	if (err != nil) || (affected != 1) {
		log.Println("WW problem saving scales, err:", err, "rows affected:", affected)
		return
	}
}

func PeriodicSave(scale int, accumulation int, event string, shift int, fraction string) {

	smt := "SELECT event, shift, fraction, datetime FROM data WHERE date(datetime) = date('now','localtime') AND scale = ? ORDER BY datetime DESC LIMIT 1"
	rows, err := db.Query(smt, scale)
	if err != nil {
		log.Println("WW can't load 2 last with err:", scale, err)
		return
	}
	defer rows.Close()
	var eventDb string
	var shiftDb int
	var fractionDb string
	var datetimeDb string

	if rows.Next() == false {
		SaveEvent(scale, accumulation, "start", shift, fraction)
		return
	}

	err = rows.Scan(&eventDb, &shiftDb, &fractionDb, &datetimeDb)
	if err != nil {
		log.Println("EE: read from db unexpected lines", err)
		SaveEvent(scale, accumulation, "start", shift, fraction)
		return
	}
	rows.Close()

	if (eventDb == "start") && (shift == shiftDb) && (fraction == fractionDb) {
		SaveEvent(scale, accumulation, event, shift, fraction)
		return
	}
	if (eventDb == event) && (shift == shiftDb) && (fraction == fractionDb) {
		// delete and save new one
		smt = "BEGIN TRANSACTION; " +
			"DELETE FROM data WHERE datetime = ? AND scale = ? AND event = ? AND fraction = ? AND shift = ?; " +
			"INSERT INTO data (scale, accumulation, event, shift, fraction, datetime) VALUES (?, ?, ?, ?, ?, datetime('now','localtime')); " +
			"COMMIT;"
		res, err := db.Exec(smt, datetimeDb, scale, eventDb, fractionDb, shiftDb, scale, accumulation, event, shift, fraction)
		if err != nil {
			log.Println("WW store: transaction failed:", err)
			return
		}
		affected, err := res.RowsAffected()
		if (err != nil) || (affected != 1) {
			log.Println("WW problem deleting old value. Err:", err, "rows affected:", affected)
		}
		return
	}
	// something changed, need to start new save line
	SaveEvent(scale, accumulation, "start", shift, fraction)
}

func ExportData(sep string) chan string {
	ret := make(chan string)
	go exportBackground(ret, sep)
	return ret
}

func exportBackground(c chan string, s string) {
	log.Println(s)
	defer close(c)
	smt := "SELECT * FROM data"
	rows, err := db.Query(smt)
	if err != nil {
		log.Println("WW can't load any data err:", err)
		return
	}
	c <- "scale" + s + "accumulation" + s + "event" + s + "shift" + s + "fraction" + s + "datetime\r\n"
	defer rows.Close()
	for rows.Next() {
		var scale, accumulation, shift int
		var event, fraction, datetime string
		err = rows.Scan(&scale, &accumulation, &event, &shift, &fraction, &datetime)
		if err != nil {
			log.Println("WW", err)
			return
		}
		c <- fmt.Sprintf("%s%s%d%s%s%s%d%s%s%s%s\r\n", scales_naming_grouping.Name[scale], s, accumulation, s, event, s, shift, s, fraction, s, datetime)
	}
}

func ExportDataStruct(scale int, start time.Time, finish time.Time) chan DataRecord {
	ret := make(chan DataRecord)
	go exportBackgroundStruct(ret, scale, start, finish)
	return ret
}

func exportBackgroundStruct(c chan DataRecord, scale int, start time.Time, finish time.Time) {
	defer close(c)
	smt := "SELECT * FROM data WHERE scale = ? AND datetime > ? AND datetime < ?"
	rows, err := db.Query(smt, scale, start, finish)
	if err != nil {
		log.Println("WW can't load any data err:", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var scale, accumulation, shift int
		var event, fraction, datetime string
		err = rows.Scan(&scale, &accumulation, &event, &shift, &fraction, &datetime)
		if err != nil {
			log.Println("WW", err)
			return
		}
		c <- DataRecord{scale, accumulation, event, shift, fraction, datetime}
	}
}

func ExportDataStructAnyTime(start time.Time, finish time.Time) chan DataRecord {
	ret := make(chan DataRecord)
	go exportBackgroundStructAnyTime(ret, start, finish)
	return ret
}

func exportBackgroundStructAnyTime(c chan DataRecord, start time.Time, finish time.Time) {
	defer close(c)
	smt := "SELECT * FROM data WHERE datetime > ? AND datetime < ?"
	rows, err := db.Query(smt, start, finish)
	if err != nil {
		log.Println("WW can't load any data err:", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var scale, accumulation, shift int
		var event, fraction, datetime string
		err = rows.Scan(&scale, &accumulation, &event, &shift, &fraction, &datetime)
		if err != nil {
			log.Println("WW", err)
			return
		}
		c <- DataRecord{scale, accumulation, event, shift, fraction, datetime}
	}
}
