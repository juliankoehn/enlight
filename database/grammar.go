package database

import (
	"bytes"
	"strings"
)

type (
	// Grammar manages translation drivers for sql syntax
	Grammar interface {
		CompileTableExists() string
		CompileColumnListing() string
		SetTablePrefix(prefix string)
		CompileCreate(bp *Blueprint, conn *Connection) string
	}
	//
	grammar struct {
		tablePrefix string
		dateFormat  string
		expression  *Expression
	}
)

// NewGrammar returns a new Grammar instance
func NewGrammar() *grammar {
	return &grammar{
		dateFormat: "Y-m-d H:i:s",
	}
}

// GetDateFormat returns the date Format
func (g *grammar) GetDateFormat() string {
	return g.dateFormat
}

// GetTablePrefix returns the table prefix
func (g *grammar) GetTablePrefix() string {
	return g.tablePrefix
}

// SetTablePrefix sets the table Prefix
func (g *grammar) SetTablePrefix(prefix string) {
	g.tablePrefix = prefix
}

// QuoteString quote the given string literal
func (g *grammar) QuoteString(value []string) string {
	return strings.Join(value, ", ")
}

func (g *grammar) getColumns(bp *Blueprint) []string {
	columnsToAdd := bp.GetAddedColumns()
	columns := make([]string, len(columnsToAdd))

	for key, value := range columnsToAdd {
		var buf bytes.Buffer

		buf.WriteString(value.Name)
		buf.WriteByte(' ')
		buf.WriteString(value.Typ)
		buf.WriteByte('(')
		buf.WriteString(value.Length)
		buf.WriteByte(')')

		columns[key] = buf.String()
	}

	return columns
}
