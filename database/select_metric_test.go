package database

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelectMetricSuccess(t *testing.T) {
	db, mock := NewMock(t)
	id := "testacct-0000-0000-0000-000000000000"

	query := `SELECT COUNT\(account_id\) FROM metric WHERE account_id \= \? FOR UPDATE;`

	rows := sqlmock.NewRows([]string{"count"}).AddRow(1)

	mock.ExpectQuery(query).WillReturnRows(rows)

	count, err := CountMetrics(nil, db, id)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	assert.Equal(t, 1, count)
}

func TestSelectMetricFailure(t *testing.T) {
	db, mock := NewMock(t)

	query := `SELECT COUNT\(account_id\) FROM metric WHERE account_id \= \? FOR UPDATE;`

	mock.ExpectQuery(query).WillReturnError(errors.New("no rows found"))

	count, err := CountMetrics(nil, db, "non-existent-account")
	assert.Equal(t, count, -1)
	assert.Error(t, err)
}
