package database

import (
	"database/sql"
	"fmt"
	"github.com/gmo-personal/coding_challenge/server/utils"
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

func closeRows(rows *sql.Rows) {
	err := rows.Close()
	if err != nil {
		utils.Logger.Println(err)
	}
}

// Begins db transaction.
func StartTransaction(db *sql.DB) (*sql.Tx, error) {
	return db.Begin()
}

// Commits or Rollbacks a db transaction.
func ResolveTransaction(tx *sql.Tx) {
	err := tx.Commit()
	if err != nil {
		utils.Logger.Println(err)
		err = tx.Rollback()
		if err != nil {
			utils.Logger.Println(err)
		}
	}
}
