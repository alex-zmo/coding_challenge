package database

import (
	"database/sql"
	"errors"
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

// Interface to take in *sql.Tx or *sql.DB.
type Caller interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// Begins db transaction.
func StartTransaction(db *sql.DB) (*sql.Tx, error) {
	return db.Begin()
}

// Commits or Rollbacks a db transaction.
func ResolveTransaction(tx *sql.Tx) error {
	err := tx.Commit()
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return err
		}
		return errors.New("transaction commit failed")
	}
	return nil
}
