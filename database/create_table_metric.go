package database

import "database/sql"

// Creates metrics table.
func CreateMetricTable(db *sql.DB) error {
	createMetricStmt := `CREATE TABLE IF NOT EXISTS metric (
		account_id varchar(36) NOT NULL,
		user_id VARCHAR(36) NOT NULL,
		time_stamp TIMESTAMP NOT NULL,
		FOREIGN KEY (account_id) REFERENCES account(id) ON DELETE CASCADE,
		UNIQUE(account_id, user_id)
	);`
	_, err := db.Exec(createMetricStmt)
	return err
}
