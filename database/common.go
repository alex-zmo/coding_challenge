package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

// Connects to the database with schema name teleport.
func InitDB() (db *sql.DB, err error) {
	serverName := os.Getenv("MYSQL_HOST") + ":" + os.Getenv("MYSQL_PORT")
	user := os.Getenv("MYSQL_USERNAME")
	password := os.Getenv("MYSQL_PASSWORD")
	schemaName := os.Getenv("MYSQL_SCHEMA")

	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, password, serverName, schemaName)
	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	return db, nil
}
