package database

import (
	"database/sql"
)

// Selects metric count as part of a transaction if tx is not nil, otherwise executes regularly.
// Returns an -1 and error if count fails.
func CountMetrics(tx *sql.Tx, db *sql.DB, accountID string) (int, error) {
	countMetricStmt := `SELECT COUNT(account_id) FROM metric WHERE account_id = ? FOR UPDATE;`

	var result *sql.Rows
	var err error
	if tx == nil {
		result, err = db.Query(countMetricStmt, accountID)
	} else {
		result, err = tx.Query(countMetricStmt, accountID)
	}

	if err != nil {
		return -1, err
	}

	defer closeRows(result)
	var count int
	if result.Next() {
		err = result.Scan(&count)
		if err != nil {
			return -1, err
		}
	}
	return count, nil
}
