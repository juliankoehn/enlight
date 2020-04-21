package database

type (
	// ColumnDefinition holds definitions of columns
	ColumnDefinition struct {
		Typ        string // kind of column, string / int
		Name       string // name of column
		Length     string // length of column
		UseCurrent bool   // CURRENT TIMESTAMP
		StoredAs   string // stored column
		VirtualAs  string // virtual column
		ColumnOptions
	}
	// ColumnOptions are given to addColumn
	ColumnOptions struct {
		Allowed       []string // allowed values
		Length        int      // length of column
		Precision     int
		AutoIncrement bool
		Unsigned      bool
		Total         int
		Places        int
	}
)

// Nullable makes a column nullable
func (c *ColumnDefinition) Nullable() *ColumnDefinition {
	return c
}
