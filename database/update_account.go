package database

import (
	"github.com/gmo-personal/coding_challenge/model"
)

// Updates account plan.
func UpdateAccount(account *model.Account) error {
	updateAccountStmt := `UPDATE account SET plan = ? WHERE id = ?;`
	_, err := db.Exec(
		updateAccountStmt,
		account.Plan,
		account.ID)
	return err
}
