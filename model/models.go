package model

// Contains account info account ID, username, password, and plan as integer.
type Account struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Plan     int    `json:"plan"`
}

// Structure for dashboard page, including CSRF token to be injected.
type DashboardPage struct {
	CSRFToken string
}

// Structure for Index(login) page, including CSRF token to be injected.
type IndexPage struct {
	CSRFToken string
}
