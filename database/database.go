package database

import (
	"crypto/tls"
	"database/sql"
	"fmt"
)

type (
	// Manager manages Databases
	Manager struct {
		config      Config
		Connections map[string]*Connection
	}

	// Connection holds the actual database Connection
	Connection struct {
		*sql.DB
		Driver string
	}

	// ConnectionConfig holds Connection Options for a Database
	ConnectionConfig struct {
		Driver                string      `json:"driver"`
		URL                   string      `json:"url"`
		Host                  string      // Host and Port
		Port                  string      // Port
		Net                   string      // Network type
		Addr                  string      // Network address (requires Net)
		Database              string      `json:"database"`
		Username              string      `json:"username"`
		Password              string      `json:"-"`
		UnixSocket            string      `json:"unix_socket"`
		Charset               string      `json:"charset"`
		Collaction            string      `json:"collation"`
		Prefix                string      `json:"prefix"`
		PrefixIndexes         bool        `json:"prefix_indexes"`
		Strict                bool        `json:"strict"`
		Engine                string      `json:"engine"`
		Schema                string      `json:"schema"`
		TLSConfig             string      `json:"sslmode"`
		tls                   *tls.Config // TLS configuration
		ForeignKeyConstraints bool        `json:"foreign_key_constraints"`
	}

	// Config holds the configuration of the Databae Manager
	Config struct {
		def         string
		connections map[string]*ConnectionConfig
	}
)

// New creates a new database manager instance.
func New() *Manager {
	conns := make(map[string]*ConnectionConfig)
	connections := make(map[string]*Connection)

	return &Manager{
		Connections: connections,
		config: Config{
			connections: conns,
		},
	}
}

// GetConnection gets a database connection instance
func (m *Manager) GetConnection(name string) (*Connection, error) {
	name = m.parseConnectionName(name)

	conn := m.Connections[name]
	if conn == nil {
		err := m.makeConnection(name)
		if err != nil {
			return nil, err
		}
	}

	return conn, nil
}

// AddConnection registers a connection with the Manager
func (m *Manager) AddConnection(conn *ConnectionConfig, name string) {
	m.config.connections[name] = conn

	if m.config.def == "" {
		m.SetDefaultConnection(name)
	}
}

// GetDefaultConnection gets the default Connection Name
func (m *Manager) GetDefaultConnection() string {
	return m.config.def
}

// SetDefaultConnection set the default connection name
func (m *Manager) SetDefaultConnection(name string) {
	m.config.def = name
}

func (m *Manager) parseConnectionName(name string) string {
	if name == "" {
		name = m.GetDefaultConnection()
	}

	return name
}

func (m *Manager) makeConnection(name string) error {
	config, err := m.getConfig(name)
	if err != nil {
		return err
	}

	_ = config

	return nil
}

// getConfig gets the configuration for a connection
func (m *Manager) getConfig(name string) (*ConnectionConfig, error) {
	if name == "" {
		name = m.GetDefaultConnection()
	}

	config := m.config.connections[name]
	if config == nil {
		return nil, fmt.Errorf("database [%s] not configured", name)
	}

	return config, nil
}
