package database

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	sqlx "github.com/jmoiron/sqlx"
)

var Engine *sqlx.DB

func Connect() error {
	var err error
	username := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	db := os.Getenv("DB_NAME")

	conn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", username, pass, host, db)

	Engine, err = sqlx.Open("mysql", conn)
	if err != nil {
		return err
	}

	return nil
}