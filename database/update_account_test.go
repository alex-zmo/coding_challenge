package database

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gmo-personal/coding_challenge/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdateUserSuccess(t *testing.T) {
	db, mock := NewMock(t)

	stmt := `UPDATE account 
		SET 
			plan \= \?
		WHERE id \= \?;`

	mock.ExpectExec(stmt).WithArgs(
		1,
		"testacct-0000-0000-0000-000000000000",
	).WillReturnResult(sqlmock.NewResult(1, 1))

	accInfo := &model.Account{
		ID:       "testacct-0000-0000-0000-000000000000",
		Username: "",
		Password: "",
		Plan:     1,
	}

	err := UpdateAccount(db, accInfo)
	assert.NoError(t, err)
}

// Unhappy Path
func TestUpdateUserFailure(t *testing.T) {
	db, mock := NewMock(t)

	stmt := `UPDATE account 
		SET plan \= \?
		WHERE id \= \?;`

	mock.ExpectExec(stmt).WillReturnError(errors.New("error"))

	accInfo := &model.Account{
		ID:       "",
		Username: "",
		Password: "",
		Plan:     1,
	}

	err := UpdateAccount(db, accInfo)
	assert.Error(t, err)
}
