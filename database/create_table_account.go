package database

// Creates account table.
func createAccountTable() error {
	createAccountStmt := `CREATE TABLE IF NOT EXISTS account (
		id VARCHAR(36) NOT NULL ,
		username VARCHAR(256) NOT NULL,
		password VARBINARY(1024) NOT NULL,
		plan INT DEFAULT 0, 
		PRIMARY KEY (id),
		UNIQUE(id),
		UNIQUE(username)
	);`
	_, err := db.Exec(createAccountStmt)
	return err
}
