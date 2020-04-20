package database

import (
	"crypto/tls"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestMake(t *testing.T) {
	factory := NewFactory()

	config := &ConnectionConfig{
		Driver:   "mysql",
		Host:     "127.0.0.1",
		Username: "testUser",
		Password: "yfLpFsBG2uMRhMaG",
		Database: "test",
	}

	conn, err := factory.Make(config, "mysql")
	if err != nil {
		t.Error(err)
	}
	if conn.DB == nil {
		t.Error("Database Connection is empty")
	}
}

func TestCreateConnector(t *testing.T) {
	factory := NewFactory()
	// empty config
	config := &ConnectionConfig{}

	_, err := factory.createConnector(config)
	if err == nil {
		t.Error("createConnector did mess up with the driver check")
	}
	config.Driver = "mysql"
	dsn, err := factory.createConnector(config)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "tcp(127.0.0.1:3306)/", dsn)

	config.Driver = "unsupported"
	_, err = factory.createConnector(config)
	if err == nil {
		t.Error("createConnector did mess up with the driver switch.")
	}
}

func TestMysqlDSN(t *testing.T) {
	factory := NewFactory()
	// empty config
	config := &ConnectionConfig{
		Driver:   "mysql",
		Username: "root",
		Password: "root",
	}

	dsn, err := factory.getMysqlDSN(config)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "root:root@tcp(127.0.0.1:3306)/", dsn)
}

func TestNormalizeMySQL(t *testing.T) {
	factory := NewFactory()
	// empty config
	config := &ConnectionConfig{
		Driver: "mysql",
	}

	// invalid UnixSocket
	config.UnixSocket = "0000"
	config.Host = ""
	_, err := factory.normalizeMySQL(config)
	if err == nil {
		t.Error("invalid UnixSocket was parsed without error")
	}

	// unix socket
	config.UnixSocket = "unix"
	config.Host = ""
	config, err = factory.normalizeMySQL(config)
	if err != nil {
		t.Error("invalid UnixSocket was parsed without error")
	}

	assert.Equal(t, "/tmp/mysql.sock", config.Host)

	// empty socket
	config.UnixSocket = ""
	config.Host = ""
	config, err = factory.normalizeMySQL(config)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "tcp", config.UnixSocket)
	assert.Equal(t, "127.0.0.1", config.Host)
	assert.Equal(t, "3306", config.Port)

	// test tls
	config.TLSConfig = "true"
	config, err = factory.normalizeMySQL(config)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, &tls.Config{ServerName: "127.0.0.1"}, config.tls)
	config.TLSConfig = "skip-verify"
	config, err = factory.normalizeMySQL(config)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, &tls.Config{InsecureSkipVerify: true}, config.tls)

	// tls error
	config.TLSConfig = "0000"
	config, err = factory.normalizeMySQL(config)
	if err == nil {
		t.Error("invalid TLSConfig was parsed without error")
	}
	if !strings.Contains(err.Error(), "invalid value / known config name:") {
		t.Error("TLS returned an incorrect error message")
	}
}
