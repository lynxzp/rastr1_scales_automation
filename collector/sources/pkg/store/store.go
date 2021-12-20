package store

import (
	"collector/pkg/config"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"runtime"
)

const (
	file string = "db.sqlite"
)

var db *sql.DB

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

	sqlStmt := `CREATE TABLE if not exists scales (id INTEGER, ip TEXT, rs485addr INTEGER, data_perf_addr INTEGER)`
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

func SaveScale(id int, dataPerfAddr int, ip string, rs485addr int) {
	ClearScale(id)
	smt := "INSERT INTO scales (id, ip, rs485addr, data_perf_addr) VALUES(?, ?, ?, ?)"
	res, err := db.Exec(smt, id, ip, rs485addr, dataPerfAddr)
	if err != nil {
		log.Println("WW can't save scales:", id, ip, rs485addr, dataPerfAddr, "with err:", err)
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
}

func ReadScales() ([]Scale, error) {
	scales := make([]Scale, config.ScalesNums)
	rows, err := db.Query(`SELECT * FROM scales`)
	if err != nil {
		return nil, errors.New("can't perform sql request")
	}
	defer rows.Close()

	var id, rs485addr, dataPerfAddr int
	var ip string
	for rows.Next() {
		err = rows.Scan(&id, &ip, &rs485addr, &dataPerfAddr)
		if err != nil {
			continue
		}
		scales[id].Rs485addr = rs485addr
		scales[id].Ip = ip
		scales[id].DataPerfAddr = dataPerfAddr
	}
	return scales, nil
}

func SaveEvent(scale int, accumulation int, event string, shift int, fraction string) {
	smt := "INSERT INTO data (scale, accumulation, event, shift, fraction, datetime) VALUES (?, ?, ?, ?, ?, datetime('now'))"
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
	SaveEvent(scale, accumulation, event, shift, fraction)

	smt := "SELECT event, shift, fraction, datetime  FROM data WHERE scale = ? ORDER BY datetime DESC LIMIT 2"
	rows, err := db.Query(smt, scale)
	if err != nil {
		log.Println("WW can't load 2 last with err:", scale, err)
		return
	}
	defer rows.Close()
	var events string
	var shifts int
	var fractions string
	var dateTimeString string

	if rows.Next() == false {
		log.Println("EE: this is impossible, no lines after writing into db")
		return
	}
	err = rows.Scan(&events, &shifts, &fractions, &dateTimeString)
	if err != nil {
		log.Println("EE: read from db unexpected lines", err)
		return
	}
	if (events != event) || (shifts != shift) || (fractions != fraction) {
		return
	}
	if rows.Next() == false {
		log.Println("WW: only 1 saved value for scale", scale, "?")
		return
	}
	err = rows.Scan(&events, &shifts, &fractions, &dateTimeString)
	if err != nil {
		log.Println("EE: read from db unexpected lines")
		return
	}
	if (events != event) || (shifts != shift) || (fractions != fraction) {
		return
	}
	rows.Close()

	smt = "DELETE FROM data WHERE datetime = ? AND scale = ? AND event = ?"
	res, err := db.Exec(smt, dateTimeString, scale, event)
	if err != nil {
		log.Println("WW can't delete double save:", scale, accumulation, event, shift, fraction, "with err:", err)
		return
	}
	affected, err := res.RowsAffected()
	if (err != nil) || (affected != 1) {
		log.Println("WW problem deleting double save, err:", err, "rows affected:", affected)
		return
	}

}

func ExportData(sep string) chan string {
	ret := make(chan string)
	go exportBackground(ret, sep)
	return ret
}

func exportBackground(c chan string, s string) {
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
		c <- fmt.Sprintf("%d%s%d%s%s%s%d%s%s%s%s\r\n", scale, s, accumulation, s, event, s, shift, s, fraction, s, datetime)
	}
}
