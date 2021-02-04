package main

import (
	"github.com/gmo-personal/coding_challenge/database"
	"github.com/gmo-personal/coding_challenge/model"
	"github.com/gmo-personal/coding_challenge/server"
	"github.com/gmo-personal/coding_challenge/server/utils"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	database.InitDatabase()
	err := AddBaseAccountIfNotExists()
	if err != nil {
			utils.LogError(err)
	}
	server.InitServer()
}

func AddBaseAccountIfNotExists() error{
	accountID := "testacct-0000-0000-0000-000000000000"
	accountUsername := "t@gmail.com"

	// Checks if base account already added
	existingAccount, err := database.SelectAccount(accountID)
	if existingAccount != nil {
		return nil
	}

	// Hardcoded Account for testing purposes. bcrypt library salts and hashes inputs.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("t"), bcrypt.DefaultCost)
	if err != nil {
			return err
	}

	// This is for fakeIOT test
	err = database.InsertAccount(&model.Account{
		ID:       accountID,
		Username: accountUsername,
		Password: string(hashedPassword),
		Plan:     0,
	})
	return err
}
