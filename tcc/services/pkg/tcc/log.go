package tcc

import (
	"database/sql"
	"fmt"
)

type TCCLoger struct {
	tid string
	db  *sql.DB
}

func NewTCCLoger(tid string, db *sql.DB) *TCCLoger {
	return &TCCLoger{
		tid: tid,
		db:  db,
	}
}

func (t *TCCLoger) LogTry() error {
	_, err := t.db.Exec(fmt.Sprintf("insert into local_try_log values('%s')", t.tid))
	return err
}

func (t *TCCLoger) LogCancel() error {
	_, err := t.db.Exec(fmt.Sprintf("insert into local_cancel_log values('%s')", t.tid))
	return err
}

func (t *TCCLoger) LogConfirm() error {
	_, err := t.db.Exec(fmt.Sprintf("insert into local_confirm_log values('%s')", t.tid))
	return err
}
