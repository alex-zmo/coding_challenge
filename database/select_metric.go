package database

// Selects metric count as part of a transaction if tx is not nil, otherwise executes regularly.
// Returns an -1 and error if count fails.
func CountMetrics(c Caller, accountID string) (int, error) {
	countMetricStmt := `SELECT COUNT(account_id) FROM metric WHERE account_id = ? FOR UPDATE;`
	count := -1
	err := c.QueryRow(countMetricStmt, accountID).Scan(&count)
	return count, err
}
