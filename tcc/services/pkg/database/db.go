package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func Connect2DB(databse string) (db *sql.DB, err error) {
	return sql.Open("mysql", "root:123456@tcp(mysql:3306)/"+databse)
}
