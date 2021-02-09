package database

import (
	"database/sql"
	"github.com/gmo-personal/coding_challenge/model"
)

// Selects an account based on account ID as part of a transaction if tx is not nil, otherwise executes regularly.
func SelectAccount(tx *sql.Tx, db *sql.DB, id string) (*model.Account, error) {
	selectAccountStmt := `
		SELECT 
			id, 
			username,
			password,
			plan
		FROM account
		WHERE id = ? FOR UPDATE;`

	var err error
	account := &model.Account{}
	if tx == nil {
		err = db.QueryRow(selectAccountStmt, id).Scan(&account.ID,
			&account.Username,
			&account.Password,
			&account.Plan)
	} else {
		err = tx.QueryRow(selectAccountStmt, id).Scan(&account.ID,
			&account.Username,
			&account.Password,
			&account.Plan)
	}

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
