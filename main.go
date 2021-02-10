package main

import (
	"database/sql"
	"github.com/gmo-personal/coding_challenge/database"
	"github.com/gmo-personal/coding_challenge/model"
	"github.com/gmo-personal/coding_challenge/server"
	"github.com/gmo-personal/coding_challenge/server/utils"
	"golang.org/x/crypto/bcrypt"
	"os"
)

func main() {
	// Initializes database connection
	db, err := database.InitDB()
	if err != nil {
		utils.Logger.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	// Creates account table if it does not already exist.
	err = database.CreateAccountTable(db)
	if err != nil {
		utils.Logger.Println(err)
		os.Exit(1)
	}
	// Creates metrics table if it does not already exist.
	err = database.CreateMetricTable(db)
	if err != nil {
		utils.Logger.Println(err)
		os.Exit(1)
	}
	// Adds base account if it does not already exist.
	err = AddAccountIfNotExists(db, "testacct-0000-0000-0000-000000000000", "t@gmail.com", "t", 0)
	if err != nil {
		utils.Logger.Println(err)
		os.Exit(1)
	}
	server.InitServer(db)
}

// Adds an account if the account doesnt already exist.
func AddAccountIfNotExists(db *sql.DB, id, username, password string, plan int) error {
	// Checks if base account already added
	existingAccount, err := database.SelectAccount(db, id)
	if existingAccount != nil {
		return nil
	}

	// Hardcoded Account for testing purposes. bcrypt library salts and hashes inputs.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// This is for fakeIOT test
	err = database.InsertAccount(db, &model.Account{
		ID:       id,
		Username: username,
		Password: string(hashedPassword),
		Plan:     plan,
	})
	return err
}
