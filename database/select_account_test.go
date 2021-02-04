package database

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

//TODO could add more functional tests on SelectAccountByUsername

// Happy Path
func TestSelectAccount(t *testing.T) {
	var mock sqlmock.Sqlmock
	db, mock = NewMock()

	query := `SELECT 
			id, 
			username,
			password,
			plan
		FROM account
		WHERE id \= \? FOR UPDATE;`

	id := "testacct-0000-0000-0000-000000000000"
	username := "test@gmail.com"
	password := "abc"
	plan := 0

	rows := sqlmock.NewRows([]string{"id", "username", "password", "plan"}).
		AddRow(id, username, password, plan)

	mock.ExpectQuery(query).WillReturnRows(rows)

	account, err := SelectAccount(id)
	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, account.ID, id)
	assert.Equal(t, account.Username, username)
	assert.Equal(t, account.Password, password)
	assert.Equal(t, account.Plan, plan)
}


// Unhappy Path
func TestSelectAccount2(t *testing.T) {
	var mock sqlmock.Sqlmock
	db, mock = NewMock()

	query := `SELECT 
			id, 
			username,
			password,
			plan
		FROM account
		WHERE id \= \? FOR UPDATE;`

	rows := sqlmock.NewRows([]string{"id", "username", "password", "plan"})

	mock.ExpectQuery(query).WillReturnRows(rows)

	account, err := SelectAccount("non-existant-account")
	assert.Nil(t, account)
	assert.Error(t, err)
}


