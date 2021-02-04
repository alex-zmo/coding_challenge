package database

import (
	"errors"
	"github.com/gmo-personal/coding_challenge/model"
)

// Selects an account based on account ID as part of a transaction if tx is not nil, otherwise executes regularly.
func SelectAccount(id string) (*model.Account, error) {
	selectAccountStmt := `
		SELECT 
			id, 
			username,
			password,
			plan
		FROM account
		WHERE id = ? FOR UPDATE;`

	result, err := db.Query(selectAccountStmt, id)

	if err != nil {
		return nil, err
	}
	defer closeRows(result)

	account := &model.Account{}
	if result.Next() {
		err = result.Scan(
			&account.ID,
			&account.Username,
			&account.Password,
			&account.Plan)

		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("select account error")
	}
	return account, nil
}

// Selects an account based on username.
func SelectAccountByUsername(username string) (*model.Account, error) {
	selectAccountStmt := `
		SELECT 
			id, 
			username,
			password,
			plan
		FROM account
		WHERE username = ?;`

	result, err := db.Query(selectAccountStmt, username)
	if err != nil {
			return nil, err
	}
	defer closeRows(result)

	account := &model.Account{}
	if result.Next() {
		err = result.Scan(
			&account.ID,
			&account.Username,
			&account.Password,
			&account.Plan)

		if err != nil {
				return nil, err
		}
	} else {
			return nil, err
	}
	return account, nil
}
