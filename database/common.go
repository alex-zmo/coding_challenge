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
	serverName := "host.docker.internal:3306"
	user := os.Getenv("MYSQL_USERNAME")
	password := os.Getenv("MYSQL_PASSWORD")
	schemaName := os.Getenv("MYSQL_SCHEMA")

	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, password, serverName, schemaName)
	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		utils.LogError(err)
		return nil, err
	}
	return db, nil
}
