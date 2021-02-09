package model

// Contains account info account ID, username, password, and plan as integer.
type Account struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Plan     int    `json:"plan"`
}

// Contains required info for dashboard, including metrics count and plan type.
type DashboardInfo struct {
	MetricsCount int `json:"metrics_count"`
	Plan         int `json:"plan"`
}

// Contains metrics info, including account ID, user ID, and timestamp.
type Metric struct {
	AccountID string `json:"account_id"`
	UserID    string `json:"user_id"`
	Timestamp string `json:"timestamp"`
}

// Structure for dashboard page, including CSRF token to be injected.
type DashboardPage struct {
	CSRFToken string
}

// Structure for Index(login) page, including CSRF token to be injected.
type IndexPage struct {
	CSRFToken string
}
