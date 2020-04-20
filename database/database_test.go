package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManager(t *testing.T) {
	manager := New()

	sql := ConnectionConfig{
		Driver: "mysql",
	}
	manager.AddConnection(&sql, "mysql")

	defaultConnection := manager.GetDefaultConnection()
	assert.Equal(t, "mysql", defaultConnection)

	connection, err := manager.GetConnection("")
	if err != nil {
		t.Error(err)
	}

	_ = connection

	conn, err := manager.GetConnection("mysql")
	if err != nil {
		t.Error(err)
	}

	t.Log(conn)
}
