package xa

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/rs/xid"
)

type XA struct {
	XID     string
	Current string
	db      *sql.DB
}

func NewXA(db *sql.DB) *XA {
	return &XA{
		db:      db,
		XID:     xid.New().String(),
		Current: "",
	}
}

func (x *XA) Start() error {
	x.Current = "start"
	return Start(x.db, x.XID)
}

func (x *XA) End() error {
	x.Current = "end"
	return End(x.db, x.XID)
}

func (x *XA) Rollback() error {
	x.Current = "rollback"
	return Rollback(x.db, x.XID)
}

func (x *XA) Commit() error {
	x.Current = "commit"
	return Commit(x.db, x.XID)
}

func (x *XA) Prepare() error {
	x.Current = "prepare"
	return Prepare(x.db, x.XID)
}

func Start(db *sql.DB, xid string) error {
	log.Printf("Start xid=%s", xid)
	_, err := db.Exec(fmt.Sprintf("XA START '%s'", xid))
	return err
}

func End(db *sql.DB, xid string) error {
	log.Printf("End xid=%s", xid)
	_, err := db.Exec(fmt.Sprintf("XA END '%s'", xid))
	return err
}

func Rollback(db *sql.DB, xid string) error {
	log.Printf("ROLLBACK xid=%s", xid)
	_, err := db.Exec(fmt.Sprintf("XA ROLLBACK '%s'", xid))
	return err
}

func Commit(db *sql.DB, xid string) error {
	log.Printf("COMMIT xid=%s", xid)
	_, err := db.Exec(fmt.Sprintf("XA COMMIT '%s'", xid))
	return err
}

func Prepare(db *sql.DB, xid string) error {
	log.Printf("PREPARE xid=%s", xid)
	_, err := db.Exec(fmt.Sprintf("XA PREPARE '%s'", xid))
	return err
}
