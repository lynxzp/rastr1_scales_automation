package store

import (
	"collector/pkg/config"
	"database/sql"
	"errors"
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
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func SaveScale(id int, dataPerfAddr int, ip string, rs485addr int) {
	smt := "INSERT OR REPLACE INTO scales (id, ip, rs485addr, data_perf_addr) VALUES(?, ?, ?, ?)"
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
	log.Println("saved scales:", id, dataPerfAddr, ip, rs485addr)
}

func ClearScale(id int) {
	smt := "DELETE FROM scales WHERE id = ?"
	res, err := db.Exec(smt, id)
	if err != nil {
		log.Println("WW can't clear scales:", id, "with err:", err)
		return
	}
	affected, err := res.RowsAffected()
	if (err != nil) || (affected != 1) {
		log.Println("WW problem saving scales, err:", err, "rows affected:", affected)
		return
	}
	log.Println("cleared scales:", id)
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
		log.Println(id, ip, rs485addr, dataPerfAddr)
	}
	return scales, nil
}
