package datasource

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func NewDatabase(username, pwd, host, name string, port uint16) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True", username, pwd, host, port, name)

	// Connect to the MySQL database.
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Test the database connection.
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
