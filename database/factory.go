package database

import (
	"bytes"
	"crypto/tls"
	"database/sql"
	"fmt"
	"strings"
	"sync"
)

type (
	// Factory builds the connection to the database
	factory struct{}
)

var (
	tlsConfigLock     sync.RWMutex
	tlsConfigRegistry map[string]*tls.Config
)

// NewFactory returns a new factory instance
func NewFactory() *factory {
	return &factory{}
}

// Make establishes a new connection based on the configuration
func (f *factory) Make(config *ConnectionConfig, name string) (*Connection, error) {
	dsn, err := f.createConnector(config)
	if err != nil {
		return nil, err
	}

	conn, err := f.openConnection(name, dsn)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		conn.Close()
		return nil, err
	}

	if config.Database == "" {
		return nil, fmt.Errorf("config.Database is not set")
	}

	// set database
	_, err = conn.DB.Exec(fmt.Sprintf("use %s;", config.Database))
	if err != nil {
		if strings.Contains(err.Error(), "Unknown database") {
			_, err = conn.DB.Exec(fmt.Sprintf("CREATE DATABASE %s;", config.Database))
			if err != nil {
				return nil, fmt.Errorf("error creating Database [%s]", config.Database)
			}
		} else {
			return nil, err
		}
	}

	return conn, err
}

func (f *factory) createConnector(config *ConnectionConfig) (string, error) {
	if config.Driver == "" {
		return "", fmt.Errorf("a driver must be specified")
	}

	switch config.Driver {
	case "mysql":
		dsn, err := f.getMysqlDSN(config)
		if err != nil {
			return "", err
		}
		return dsn, nil
	case "pgsql":
	case "sqlite":
	case "sqlsrv":
	}

	return "", fmt.Errorf("unsupported driver [%s]", config.Driver)
}

func (f *factory) getMysqlDSN(config *ConnectionConfig) (string, error) {
	config, err := f.normalizeMySQL(config)
	if err != nil {
		return "", err
	}
	// Standard Connection String
	// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	// Server=myServerAddress;Database=myDataBase;Uid=myUsername;Pwd=myPassword;
	var buf bytes.Buffer

	// [username[:password]@]
	if len(config.Username) > 0 {
		buf.WriteString(config.Username)
		if len(config.Password) > 0 {
			buf.WriteByte(':')
			buf.WriteString(config.Password)
		}
		buf.WriteByte('@')
	}
	// [protocol[(address)]]
	if len(config.UnixSocket) > 0 {
		buf.WriteString(config.UnixSocket)
		if len(config.Host) > 0 {
			buf.WriteByte('(')
			buf.WriteString(config.Host)
			if len(config.Port) > 0 {
				buf.WriteByte(':')
				buf.WriteString(config.Port)
			}
			buf.WriteByte(')')
		}
	}

	// /dbname
	buf.WriteByte('/')

	return buf.String(), nil
}

func (f *factory) normalizeMySQL(config *ConnectionConfig) (*ConnectionConfig, error) {
	if config.UnixSocket == "" {
		config.UnixSocket = "tcp"
	}

	if config.Host == "" {
		switch config.UnixSocket {
		case "tcp":
			config.Host = "127.0.0.1"
			if config.Port == "" {
				config.Port = "3306"
			}
		case "unix":
			config.Host = "/tmp/mysql.sock"
		default:
			return nil, fmt.Errorf("default addr for network [%s] unknown", config.UnixSocket)
		}
	}

	switch config.TLSConfig {
	case "false", "":
		// nothing
	case "true":
		config.tls = &tls.Config{}
	case "skip-verify":
		config.tls = &tls.Config{
			InsecureSkipVerify: true,
		}
	default:
		config.tls = getTLSConfigClone(config.TLSConfig)
		if config.tls == nil {
			return nil, fmt.Errorf("invalid value / known config name: [%s]", config.TLSConfig)
		}
	}

	if config.tls != nil && config.tls.ServerName == "" && !config.tls.InsecureSkipVerify {
		config.tls.ServerName = config.Host
	}

	return config, nil
}

func (f *factory) openConnection(driver string, dataSourceName string) (*Connection, error) {
	db, err := sql.Open(driver, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &Connection{
		DB:     db,
		Driver: driver,
	}, nil
}

func getTLSConfigClone(key string) (config *tls.Config) {
	tlsConfigLock.RLock()
	if v, ok := tlsConfigRegistry[key]; ok {
		config = v.Clone()
	}
	tlsConfigLock.RUnlock()
	return
}
