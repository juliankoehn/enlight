package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManager(t *testing.T) {
	manager := New()

	sql := ConnectionConfig{
		Driver:   "mysql",
		Host:     "127.0.0.1",
		Username: "testUser",
		Password: "yfLpFsBG2uMRhMaG",
		Database: "test",
	}
	manager.AddConnection(&sql, "mysql")

	defaultConnection := manager.GetDefaultConnection()
	assert.Equal(t, "mysql", defaultConnection)

	_, err := manager.GetConnection("")
	if err != nil {
		t.Error(err)
	}

	if _, err := manager.GetConnection("mysql"); err != nil {
		t.Error(err)
	}

	// testing disconnect
	if err := manager.Disconnect("mysql"); err != nil {
		t.Error(err)
	}

	if err := manager.Disconnect("unknown"); err == nil {
		t.Error(err)
	}

	// reconnect to database
	if _, err := manager.Reconnect("mysql"); err != nil {
		t.Error(err)
	}

	if _, err := manager.Reconnect("unknown"); err == nil {
		t.Error(err)
	}

	// getConfig
	if _, err := manager.getConfig("unknown"); err == nil {
		t.Error(err)
	}

	if config, err := manager.getConfig("mysql"); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, "mysql", config.Driver)
	}

	// we have just 1 connection setup so far
	if config, err := manager.getConfig(""); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, "mysql", config.Driver)
	}
}
