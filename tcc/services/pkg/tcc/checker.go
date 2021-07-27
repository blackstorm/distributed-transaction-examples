package tcc

import (
	"database/sql"
	"log"
)

type TCCChecker struct {
	tid string
	db  *sql.DB
}

func NewTCCChecker(tid string, db *sql.DB) *TCCChecker {
	return &TCCChecker{
		tid: tid,
		db:  db,
	}
}

func (t *TCCChecker) IsTried() (bool, error) {
	log.Print("tcc checker is tried")
	row := 0
	err := t.db.QueryRow("select count(1) from local_try_log where tx_no = ?", t.tid).Scan(&row)
	return row > 0, err
}

func (t *TCCChecker) IsCancel() (bool, error) {
	log.Print("tcc checker is cancel")
	row := 0
	err := t.db.QueryRow("select count(1) from local_cancel_log where tx_no = ?", t.tid).Scan(&row)
	return row > 0, err
}
