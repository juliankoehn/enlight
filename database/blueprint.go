package database

import (
	"fmt"
	"strconv"
)

type (
	// ColumnDefinition holds definitions of columns
	ColumnDefinition struct {
		Typ           string // kind of column, string / int
		Name          string // name of column
		Length        string // length of column
		Total         int
		Places        int
		Allowed       []string // allowed values
		Precision     int
		UseCurrent    bool   // CURRENT TIMESTAMP
		StoredAs      string // stored column
		VirtualAs     string // virtual column
		AutoIncrement bool
		Unsigned      bool
	}
	// ColumnOptions are given to addColumn
	ColumnOptions struct {
		Length        int
		AutoIncrement bool
		Unsigned      bool
	}
	// Blueprint is a
	Blueprint struct {
		table     string              // the table the blueprint describes.
		prefix    string              // the prefix of the table
		columns   []*ColumnDefinition // columns that should be added to the table
		commands  []string            //
		temporary bool                // Whether to make the table temporary.
		charset   string              // The default character set that should be used for the table.
		collation string              // The collation that should be used for the table.
	}
)

// NewBlueprint creates a new blueprint
func NewBlueprint(table string, prefix string) *Blueprint {
	bp := &Blueprint{
		table:  table,
		prefix: prefix,
	}
	return bp
}

// GetAddedColumns get the columns on the blueprint that should be added.
func (b *Blueprint) GetAddedColumns() []*ColumnDefinition {
	return b.columns
}

// GetCharset returns the charset of blueprint
func (b *Blueprint) GetCharset() string {
	return b.charset
}

// GetCollation returns the collation of blueprint
func (b *Blueprint) GetCollation() string {
	return b.collation
}

func (b *Blueprint) addCommand(name string) {
	b.commands = append(b.commands, name)
}

func (b *Blueprint) addColumn(typ, name string, options *ColumnOptions) *ColumnDefinition {
	definition := &ColumnDefinition{
		Typ:  typ,
		Name: name,
	}

	if options != nil {
		if options.Length > 0 {
			definition.Length = strconv.Itoa(options.Length)
		}
	}

	b.columns = append(b.columns, definition)

	return definition
}

// Create indicate that the table needs to be created.
func (b *Blueprint) Create() {
	b.addCommand("create")
}

// Temporary indicate that the table needs to be temporary.
func (b *Blueprint) Temporary() {
	b.temporary = true
}

// Drop indicate that the table should be dropped.
func (b *Blueprint) Drop() {
	b.addCommand("drop")
}

// Char create a new char column on the table
func (b *Blueprint) Char(column string, length int) *ColumnDefinition {
	if length == 0 {
		length = 255
	}

	return b.addColumn("char", column, &ColumnOptions{
		Length: length,
	})
}

// String create a new string column on the table
func (b *Blueprint) String(column string, length int) *ColumnDefinition {
	if length == 0 {
		length = 255
	}
	return b.addColumn("string", column, &ColumnOptions{
		Length: length,
	})
}

// Text create a new text column on the table
func (b *Blueprint) Text(column string) *ColumnDefinition {
	return b.addColumn("text", column, nil)
}

// MediumText create a new mediumText column on the table
func (b *Blueprint) MediumText(column string) *ColumnDefinition {
	return b.addColumn("mediumText", column, nil)
}

// LongText create a new longText column on the table
func (b *Blueprint) LongText(column string) *ColumnDefinition {
	return b.addColumn("longText", column, nil)
}

// Integer create a new integer column on the table
func (b *Blueprint) Integer(column string, autoIncrement, unsigned bool) *ColumnDefinition {
	return b.addColumn("integer", column, &ColumnOptions{
		AutoIncrement: autoIncrement,
		Unsigned:      unsigned,
	})
}

// TinyInteger create a new tiny integer (1-byte) column on the table
func (b *Blueprint) TinyInteger(column string, autoIncrement, unsigned bool) *ColumnDefinition {
	return b.addColumn("tinyInteger", column, &ColumnOptions{
		AutoIncrement: autoIncrement,
		Unsigned:      unsigned,
	})
}

// SmallInteger create a new small integer (2-byte) column on the table
func (b *Blueprint) SmallInteger(column string, autoIncrement, unsigned bool) *ColumnDefinition {
	return b.addColumn("smallInteger", column, &ColumnOptions{
		AutoIncrement: autoIncrement,
		Unsigned:      unsigned,
	})
}

// MediumInteger create a new mediumInteger (3-byte) column on the table
func (b *Blueprint) MediumInteger(column string, autoIncrement, unsigned bool) *ColumnDefinition {
	return b.addColumn("mediumInteger", column, &ColumnOptions{
		AutoIncrement: autoIncrement,
		Unsigned:      unsigned,
	})
}

// BigInteger create a new bigInteger (8-byte) column on the table
func (b *Blueprint) BigInteger(column string, autoIncrement, unsigned bool) *ColumnDefinition {
	return b.addColumn("bigInteger", column, &ColumnOptions{
		AutoIncrement: autoIncrement,
		Unsigned:      unsigned,
	})
}

// UnsignedInteger create a new unsigned integer (4-byte) column on the table
func (b *Blueprint) UnsignedInteger(column string, autoIncrement bool) *ColumnDefinition {
	return b.Integer(column, autoIncrement, true)
}

// UnsignedTinyInteger create a new unsigned tiny integer (1-byte) column on the table
func (b *Blueprint) UnsignedTinyInteger(column string, autoIncrement bool) *ColumnDefinition {
	return b.TinyInteger(column, autoIncrement, true)
}

// UnsignedSmallInteger create a new unsigned small integer (2-byte) column on the table
func (b *Blueprint) UnsignedSmallInteger(column string, autoIncrement bool) *ColumnDefinition {
	return b.SmallInteger(column, autoIncrement, true)
}

// UnsignedMediumInteger create a new medium integer (3-byte) column on the table
func (b *Blueprint) UnsignedMediumInteger(column string, autoIncrement bool) *ColumnDefinition {
	return b.MediumInteger(column, autoIncrement, true)
}

// UnsignedBigInteger create a new big integer (8-byte) column on the table
func (b *Blueprint) UnsignedBigInteger(column string, autoIncrement bool) *ColumnDefinition {
	return b.BigInteger(column, autoIncrement, true)
}

func (b *Blueprint) toSQL(conn *Connection, grammar Grammar) []string {
	var statements []string

	for _, cmd := range b.commands {
		switch cmd {
		case "create":
			statements = append(statements, grammar.CompileCreate(b, conn))
		}
	}

	return statements
}

// Execute the blueprint against the database.
func (b *Blueprint) Execute(conn Connection) []string {
	grammar := conn.GetQueryGrammar()
	statements := b.toSQL(&conn, grammar)

	for _, statement := range statements {
		fmt.Println(statement)
		if _, err := conn.Exec(statement); err != nil {
			panic(err)
		}
	}
	return statements
}
