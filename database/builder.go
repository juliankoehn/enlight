package database

import (
	"strings"

	"github.com/juliankoehn/enlight/support/str"
)

type (
	// Builder is a database Schema Manager
	Builder struct {
		Connection          Connection // The database connection instance
		Grammar             string     // The schema grammar instance
		defaultStringLength int        // the default string length
	}
)

// NewBuilder returns a new database schema manager
func NewBuilder(conn Connection) *Builder {
	return &Builder{
		Connection:          conn,
		defaultStringLength: 255,
	}
}

// SetdefaultStringLength sets the defaultStringLength
func (b *Builder) SetdefaultStringLength(length int) {
	b.defaultStringLength = length
}

// HasTable checks if given table already exists
func (b *Builder) HasTable(table string) (bool, error) {
	dbName := b.Connection.config.Database
	tName := b.Connection.GetTablePrefix() + table
	var count int

	query := b.Connection.GetQueryGrammar().CompileTableExists()
	if err := b.Connection.QueryRow(query, dbName, tName).Scan(&count); err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

// HasColumn Determine if the given table has a given column.
func (b *Builder) HasColumn(table string, column string) (bool, error) {
	column = strings.ToLower(column)
	table = b.Connection.GetTablePrefix() + strings.ToLower(table)

	query := b.Connection.GetQueryGrammar().CompileColumnListing()
	var columnName string
	if err := b.Connection.QueryRow(query, b.Connection.config.Database, table).Scan(&columnName); err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			return false, nil
		}
		return false, err
	}

	if columnName == column {
		return true, nil
	}

	return false, nil
}

// HasColumns Determine if the given table has given columns.
func (b *Builder) HasColumns(table string, columns []string) (bool, error) {
	// converting potential uppercases to lowercases
	for i := range columns {
		columns[i] = strings.ToLower(columns[i])
	}
	table = b.Connection.GetTablePrefix() + strings.ToLower(table)

	query := b.Connection.GetQueryGrammar().CompileColumnListing()
	var tableColumns []string
	if err := b.Connection.Run(query, &tableColumns, b.Connection.config.Database, table); err != nil {
		return false, err
	}

	for _, col := range columns {
		if !str.Contains(tableColumns, strings.ToLower(col)) {
			return false, nil
		}
	}
	return true, nil
}

// GetColumnListing gets the column listing for a given table
func (b *Builder) GetColumnListing(table string) ([]string, error) {
	table = b.Connection.GetTablePrefix() + strings.ToLower(table)
	query := b.Connection.GetQueryGrammar().CompileColumnListing()

	var tableColumns []string
	if err := b.Connection.Run(query, &tableColumns, b.Connection.config.Database, table); err != nil {
		return nil, err
	}
	return tableColumns, nil
}

// Table Modify a table on the schema.
func (b *Builder) Table(table string, callback func(...interface{})) {
	b.Build(b.CreateBlueprint(table))
}

// Create a new table on the schema.
func (b *Builder) Create(table string, callback func(*Blueprint)) []string {
	bp := b.CreateBlueprint(table)
	bp.Create()
	callback(bp)

	return b.Build(bp)
}

// CreateBlueprint Create a new command set with a Closure.
func (b *Builder) CreateBlueprint(table string) *Blueprint {
	table = strings.ToLower(table)
	prefix := b.Connection.GetTablePrefix()

	return NewBlueprint(table, prefix)
}

// Build execute the blueprint to build / modify the table
func (b *Builder) Build(blueprint *Blueprint) []string {
	return blueprint.Execute(b.Connection)
}
