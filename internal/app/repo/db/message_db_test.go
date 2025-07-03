package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/pkg/db"
)

const (
	_level          = "1"
	_id             = "6"
	_defaultTestDSN = "test_user:test_password@tcp(127.0.0.1:3306)/test_db?parseTime=true&timeout=10s"
)

var _repo MessageRepoDB

func TestMain(m *testing.M) {
	// use TEST_DB_DSN env variable if it exists
	dsn, ok := os.LookupEnv("TEST_DB_DSN")
	if !ok {
		dsn = _defaultTestDSN
	}

	// create repo
	dbStorage, err := db.New(dsn,
		db.WithTranslateError(),
		db.WithDisableColorful())
	if err != nil {
		panic(err)
	}
	_repo = NewMessageRepoDB(dbStorage)
	// run tests
	os.Exit(m.Run())
}

func TestGetLevelsCount(t *testing.T) {
	t.Log("Get levels count")

	levelsCount, err := _repo.GetLevelsCount()
	require.NoError(t, err, "get levels count")

	t.Logf("levelsCount: %+v", levelsCount)
}

func TestGetWithLevel(t *testing.T) {
	t.Log("Get slice of messages with given level")

	msgSlice, err := _repo.GetWithLevel(_level)
	require.NoError(t, err, "get slice of messages")

	t.Logf("msgSlice: %+v", msgSlice)
}

func TestGetByID(t *testing.T) {
	t.Log("Get message by ID")

	msg := &entity.Message{ID: _id}

	err := _repo.GetByID(msg)
	require.NoError(t, err, "get message by id")

	t.Logf("msg: %+v", msg)
}

func TestDeleteByID(t *testing.T) {
	t.Log("Delete message by ID")

	msg := &entity.Message{ID: _id}

	err := _repo.DeleteByID(msg)
	require.NoError(t, err, "delete message by id")
}

func TestDeleteWithLevel(t *testing.T) {
	t.Log("Delete all messages with given level")

	err := _repo.DeleteAllWithLevel(_level)
	require.NoError(t, err, "delete messages with level")
}
