package database

type (
	// ColumnDefinition holds definitions of columns
	ColumnDefinition struct {
		Typ        string // kind of column, string / int
		Name       string // name of column
		UseCurrent bool   // CURRENT TIMESTAMP
		ColumnOptions
	}
	// ColumnOptions are given to addColumn
	ColumnOptions struct {
		Length        int      // length of column
		storedAs      string   // stored column
		virtualAs     string   // virtual column
		Allowed       []string // allowed values
		Precision     int
		autoIncrement bool
		unsigned      bool
		Total         int
		Places        int
		first         bool
		nullable      bool   // determines if a column is nullable
		charset       string // custom charset for column
		collation     string // collation of column
		def           string // column default
		srid          int    // SRID
		comment       string // holds a column comment
		after         string // col after given
	}
)

// VirtualAs Create a virtual generated column (MySQL)
func (c *ColumnDefinition) VirtualAs(as string) *ColumnDefinition {
	c.virtualAs = as
	return c
}

// StoredAs Create a stored generated column (MySQL)
func (c *ColumnDefinition) StoredAs(as string) *ColumnDefinition {
	c.storedAs = as
	return c
}

// Unsigned Set INTEGER columns as UNSIGNED (MySQL)
func (c *ColumnDefinition) Unsigned() *ColumnDefinition {
	c.unsigned = true
	return c
}

// First Place the column "first" in the table (MySQL)
func (c *ColumnDefinition) First() *ColumnDefinition {
	c.first = true
	return c
}

// Default Specify a "default" value for the column
func (c *ColumnDefinition) Default(def string) *ColumnDefinition {
	c.def = def
	return c
}

// Comment Add a comment to a column (MySQL/PostgreSQL)
func (c *ColumnDefinition) Comment(comm string) *ColumnDefinition {
	c.comment = comm
	return c
}

// Collaction Specify a collation for the column (MySQL/PostgreSQL/SQL Server)
func (c *ColumnDefinition) Collaction(coll string) *ColumnDefinition {
	c.collation = coll
	return c
}

// Charset Specify a character set for the column (MySQL)
func (c *ColumnDefinition) Charset(chars string) *ColumnDefinition {
	c.charset = chars
	return c
}

// AutoIncrement set INTEGER columns as auto-increment (primary key)
func (c *ColumnDefinition) AutoIncrement() *ColumnDefinition {
	c.autoIncrement = true
	return c
}

// After place the column "after" another column (MySQL)
func (c *ColumnDefinition) After(column string) *ColumnDefinition {
	c.after = column
	return c
}

// Nullable makes a column nullable
func (c *ColumnDefinition) Nullable() *ColumnDefinition {
	c.nullable = true
	return c
}
