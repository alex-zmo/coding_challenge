package database

import (
	"database/sql"
	"fmt"
	"github.com/gmo-personal/coding_challenge/server/utils"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var db *sql.DB

// Inits the database with schema name teleport and all tables
func InitDatabase() {
	var err error
	// For Docker Image Usage
	serverName := "host.docker.internal:3306"
	user := os.Getenv("MYSQL_USERNAME")
	password := os.Getenv("MYSQL_PASSWORD")
	// For Local Running Usage
	//serverName := "localhost:3306"
	//user := "root"
	//password := ""
	schemaName := "teleport"

	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, password, serverName, schemaName)
	db, err = sql.Open("mysql", connectionString)
	if err != nil {
			panic(err)
	}
	err = createAccountTable()
	if err != nil {
			panic(err)
	}
}

// Closes rows on table.
func closeRows(rows *sql.Rows) {
	err := rows.Close()
	if err != nil {
			utils.LogError(err)
	}
}
