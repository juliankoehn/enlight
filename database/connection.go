package database

import (
	"database/sql"
	"reflect"
	"strings"
	"sync"

	"github.com/jmoiron/sqlx/reflectx"
)

type (
	// Connection holds the actual database Connection
	Connection struct {
		*sql.DB
		Driver       string
		queryGrammar Grammar
		tablePrefix  string
		config       ConnectionConfig
		unsafe       bool
		Mapper       *reflectx.Mapper
	}

	// Connections holds connections Key:Value
	Connections map[string]*Connection

	// Rows is a wrapper around sql.Rows which caches costly reflect operations
	// during a looped StructScan
	Rows struct {
		*sql.Rows
		unsafe  bool
		Mapper  *reflectx.Mapper
		started bool
		fields  [][]int
		values  []interface{}
	}
)

// NewConnection returns a new Connection instance
func NewConnection(db *sql.DB, driver string) *Connection {
	m := mapper()
	return &Connection{
		DB:     db,
		Driver: driver,
		Mapper: m,
	}
}

// SetConfig sets the config
func (c Connection) SetConfig(cfg ConnectionConfig) {
	c.config = cfg
}

// GetTablePrefix returns the table Prefix
func (c Connection) GetTablePrefix() string {
	return c.tablePrefix
}

// SetTablePrefix set the table prefix in use by the connection.
func (c Connection) SetTablePrefix(prefix string) {
	c.tablePrefix = prefix
	c.GetQueryGrammar().SetTablePrefix(prefix)
}

// UseDefaultQueryGrammar set the query grammar to the default implementation
func (c Connection) UseDefaultQueryGrammar() Grammar {
	g := c.getDefaultQueryGrammar()
	c.queryGrammar = g
	return g
}

// getDefaultQueryGrammar gets the default query grammar instance.
func (c Connection) getDefaultQueryGrammar() Grammar {
	var g Grammar

	switch c.config.Driver {
	case "mysql":
		g = NewMySQLGrammar()
	}
	return g
}

// GetQueryGrammar gets the query grammar used by connection
func (c Connection) GetQueryGrammar() Grammar {
	if c.queryGrammar == nil {
		return c.UseDefaultQueryGrammar()
	}

	return c.queryGrammar
}

// SetQueryGrammar set the query grammar used by the connection
func (c Connection) SetQueryGrammar(grammar Grammar) {
	c.queryGrammar = grammar
}

// Select using this DB
func (c *Connection) Select(dest interface{}, query string, args ...interface{}) error {
	return Select(c, dest, query, args...)
}

// QueryRow selets one row from database
func (c *Connection) QueryRow(query string, args ...interface{}) *Row {
	rows, err := c.DB.Query(query, args...)
	return &Row{rows: rows, err: err, unsafe: c.unsafe, Mapper: c.Mapper}
}

func (c *Connection) Run(query string, dest interface{}, args ...interface{}) error {
	// reconnectIfMissingConnection
	rows, err := c.DB.Query(query, args...)
	if err != nil {
		return err
	}

	return scanAll(rows, dest, false)
}

// Although the NameMapper is convenient, in practice it should not
// be relied on except for application code.  If you are writing a library
// that uses sqlx, you should be aware that the name mappings you expect
// can be overridden by your user's application.

// NameMapper is used to map column names to struct field names.  By default,
// it uses strings.ToLower to lowercase struct field names.  It can be set
// to whatever you want, but it is encouraged to be set before sqlx is used
// as name-to-field mappings are cached after first use on a type.
var NameMapper = strings.ToLower
var origMapper = reflect.ValueOf(NameMapper)

// Rather than creating on init, this is created when necessary so that
// importers have time to customize the NameMapper.
var mpr *reflectx.Mapper

// mprMu protects mpr.
var mprMu sync.Mutex

// mapper returns a valid mapper using the configured NameMapper func.
func mapper() *reflectx.Mapper {
	mprMu.Lock()
	defer mprMu.Unlock()

	if mpr == nil {
		mpr = reflectx.NewMapperFunc("db", NameMapper)
	} else if origMapper != reflect.ValueOf(NameMapper) {
		// if NameMapper has changed, create a new mapper
		mpr = reflectx.NewMapperFunc("db", NameMapper)
		origMapper = reflect.ValueOf(NameMapper)
	}
	return mpr
}

// Queryx queries the database and returns an *database.Rows.
// Any placeholder parameters are replaced with supplied args.
func (c *Connection) Queryx(query string, args ...interface{}) (*Rows, error) {
	r, err := c.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return &Rows{Rows: r, unsafe: c.unsafe, Mapper: c.Mapper}, nil
}
