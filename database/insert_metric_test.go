package database

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gmo-personal/coding_challenge/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInsertMetricSuccess(t *testing.T) {
	db, mock := NewMock(t)

	stmt := `INSERT INTO metric \(
		account_id,
		user_id,
		time_stamp
	\) VALUES \(\?, \?, \?\) ON DUPLICATE KEY UPDATE time_stamp = \?;`

	mock.ExpectExec(stmt).WithArgs(
		"testacct-0000-0000-0000-000000000000",
		"testuser-0000-0000-0000-000000000000",
		"2021-02-07 17:25:39",
		"2021-02-07 17:25:39",
	).WillReturnResult(sqlmock.NewResult(1, 1))

	metricsInfo := &model.Metric{
		AccountID: "testacct-0000-0000-0000-000000000000",
		UserID:    "testuser-0000-0000-0000-000000000000",
		Timestamp: "2021-02-07 17:25:39",
	}

	err := InsertMetric(nil, db, metricsInfo)
	assert.NoError(t, err)
}

func TestInsertMetricFailure(t *testing.T) {
	db, mock := NewMock(t)

	stmt := `INSERT INTO metric \(
		account_id,
		user_id,
		time_stamp
	\) VALUES \(\?, \?, \?\) ON DUPLICATE KEY UPDATE time_stamp = \?;`

	mock.ExpectExec(stmt).WillReturnError(errors.New("error"))
	metricsInfo := &model.Metric{}
	err := InsertMetric(nil, db, metricsInfo)
	assert.Error(t, err)
}
