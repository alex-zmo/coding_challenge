package database

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gmo-personal/coding_challenge/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Happy Path
func TestUpdateUser(t *testing.T) {
	var mock sqlmock.Sqlmock
	db, mock = NewMock()

	stmt := `UPDATE account 
		SET 
			plan \= \?
		WHERE id \= \?;`

	mock.ExpectExec(stmt).WithArgs(
		1,
		"testacct-0000-0000-0000-000000000000",
	).WillReturnResult(sqlmock.NewResult(1, 1))

	accInfo := &model.Account{
		ID:      "testacct-0000-0000-0000-000000000000",
		Username: "",
		Password: "",
		Plan:     1,
	}

	err := UpdateAccount(accInfo)
	assert.NoError(t, err)
}

// Unhappy Path
func TestUpdateUser2(t *testing.T) {
	var mock sqlmock.Sqlmock
	db, mock = NewMock()

	stmt := `UPDATE account 
		SET plan \= \?
		WHERE id \= \?;`

	mock.ExpectExec(stmt).WillReturnError(errors.New("error"))

	accInfo := &model.Account{
		ID:      "",
		Username: "",
		Password: "",
		Plan:     1,
	}

	err := UpdateAccount(accInfo)
	assert.Error(t, err)
}
