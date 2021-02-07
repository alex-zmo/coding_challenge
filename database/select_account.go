package database

import (
	"database/sql"
	"github.com/gmo-personal/coding_challenge/model"
)

// Selects an account based on account ID.
func SelectAccount(db *sql.DB, id string) (*model.Account, error) {
	selectAccountStmt := `
		SELECT 
			id, 
			username,
			password,
			plan
		FROM account
		WHERE id = ? FOR UPDATE;`

	account := &model.Account{}
	err := db.QueryRow(selectAccountStmt, id).Scan(&account.ID,
		&account.Username,
		&account.Password,
		&account.Plan)
	if err != nil {
		return nil, err
	}
	return account, nil
}

// Selects an account based on username.
func SelectAccountByUsername(db *sql.DB, username string) (*model.Account, error) {
	selectAccountStmt := `
		SELECT 
			id, 
			username,
			password,
			plan
		FROM account
		WHERE username = ?;`

	account := &model.Account{}
	err := db.QueryRow(selectAccountStmt, username).Scan(&account.ID,
		&account.Username,
		&account.Password,
		&account.Plan)
	if err != nil {
		return nil, err
	}
	return account, nil
}
