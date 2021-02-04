package database

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gmo-personal/coding_challenge/model"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

// Happy Path
func TestInsertUser(t *testing.T) {
	var mock sqlmock.Sqlmock
	db, mock = NewMock()

	stmt := `INSERT INTO account \(
		id,
		username,
		password,
		plan
	\) VALUES \(\?, \?, \?, \?\);`

	hashedPassword := "abc"
	mock.ExpectExec(stmt).WithArgs(
		"testacct-0000-0000-0000-000000000000",
		"test@gmail.com",
		hashedPassword,
		0,
	).WillReturnResult(sqlmock.NewResult(1, 1))

	accInfo := &model.Account{
		ID:      "testacct-0000-0000-0000-000000000000",
		Username: "test@gmail.com",
		Password: hashedPassword,
		Plan:     0,
	}

	err := InsertAccount(accInfo)
	assert.NoError(t, err)
}

// Unhappy Path
func TestInsertUser2(t *testing.T) {
	var mock sqlmock.Sqlmock
	db, mock = NewMock()

	stmt := `INSERT INTO account \(
		id,
		username,
		password,
		plan
	\) VALUES \(\?, \?, \?, \?\);`

	mock.ExpectExec(stmt).WillReturnError(errors.New("error"))
	accInfo := &model.Account{}

	err := InsertAccount(accInfo)
	assert.Error(t, err)
}
