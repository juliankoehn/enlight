package database

import (
	"crypto/tls"
	"fmt"
)

type (
	// Manager manages Databases
	Manager struct {
		config      Config
		Factory     *factory
		Connections Connections
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
	factory := NewFactory()

	return &Manager{
		Connections: connections,
		Factory:     factory,
		config: Config{
			connections: conns,
		},
	}
}

// GetConnection gets a database connection instance
func (m *Manager) GetConnection(name string) (*Connection, error) {
	var err error
	name = m.parseConnectionName(name)

	conn := m.Connections[name]
	if conn == nil {
		conn, err = m.makeConnection(name)
		if err != nil {
			return nil, err
		}

		m.Connections[name] = conn
	}
	return conn, nil
}

// GetConnections returns Connections
func (m *Manager) GetConnections() Connections {
	return m.Connections
}

// Disconnect disconnects from the given database
func (m *Manager) Disconnect(name string) error {
	if conn, ok := m.Connections[name]; ok {
		if err := conn.Close(); err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("connection with name [%s] not found", name)
}

// Reconnect to given database
func (m *Manager) Reconnect(name string) (*Connection, error) {
	if name == "" {
		name = m.GetDefaultConnection()
	}
	// disconnect from database
	if err := m.Disconnect(name); err != nil {
		return nil, err
	}
	if _, ok := m.Connections[name]; !ok {
		return m.GetConnection(name)
	}
	// reconenct to database
	return m.makeConnection(name)
}

// AddConnection registers a connection with the Manager
// If it's the first connection it gets automatically set to default
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

func (m *Manager) makeConnection(name string) (*Connection, error) {
	config, err := m.getConfig(name)
	if err != nil {
		return nil, err
	}

	conn, err := m.Factory.Make(config, name)
	if err != nil {
		return nil, err
	}

	conn.config = *config

	return conn, nil
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
