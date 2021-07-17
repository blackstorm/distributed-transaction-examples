package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() *sql.DB {
	if db, err := sql.Open("mysql", "root:123456@tcp(mysql:3306)/xa"); err != nil {
		panic(err)
	} else {
		// db.SetConnMaxLifetime(10000)
		// db.SetMaxIdleConns(1)
		// db.SetMaxOpenConns(1)
		return db
	}
}
