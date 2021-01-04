package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func Init(dataSource string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
