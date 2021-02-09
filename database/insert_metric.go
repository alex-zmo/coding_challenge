package database

import (
	"database/sql"
	"github.com/gmo-personal/coding_challenge/model"
)

// Inserts metric as part of a transaction if tx is not nil, otherwise executes regularly.
// Updates duplicate unique keys (account_id, user_id) with new timestamp.
func InsertMetric(tx *sql.Tx, db *sql.DB, metricInfo *model.Metric) error {
	insertMetricStmt := `INSERT INTO metric (
		account_id,
		user_id,
		time_stamp
	) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE time_stamp = ?;`

	var err error
	if tx == nil {
		_, err = db.Exec(
			insertMetricStmt,
			metricInfo.AccountID,
			metricInfo.UserID,
			metricInfo.Timestamp,
			metricInfo.Timestamp)
	} else {
		_, err = tx.Exec(
			insertMetricStmt,
			metricInfo.AccountID,
			metricInfo.UserID,
			metricInfo.Timestamp,
			metricInfo.Timestamp)
	}
	return err
}
