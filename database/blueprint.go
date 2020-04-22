package database

import (
	"fmt"
	"strings"
)

type (
	// Blueprint is a
	Blueprint struct {
		table     string              // the table the blueprint describes.
		prefix    string              // the prefix of the table
		columns   []*ColumnDefinition // columns that should be added to the table
		commands  []*Command          //
		temporary bool                // Whether to make the table temporary.
		charset   string              // The default character set that should be used for the table.
		collation string              // The collation that should be used for the table.
	}

	Command struct {
		Typ string
		CommandOptions
	}

	CommandOptions struct {
		Index     string
		Columns   []string
		Algorithm string
		To        string
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

func (b *Blueprint) addCommand(name string, options *CommandOptions) *Blueprint {
	command := &Command{
		Typ: name,
	}
	if options != nil {
		command.Index = options.Index
		command.Columns = options.Columns
		command.Algorithm = options.Algorithm
	}

	b.commands = append(b.commands, command)

	return b
}

// Add a new index command to the blueprint
func (b *Blueprint) indexCommand(typ string, columns []string, index string, algorithm string) *Blueprint {
	// if no name was specified for this index, we will create one using a bsaic
	// convention of the table name, followed by the columns, followd by an
	// index type, such as primary or index, which makes the index unique.
	if index == "" {
		index = b.createIndexName(typ, columns)
	}

	return b.addCommand(typ, &CommandOptions{
		Index:     index,
		Columns:   columns,
		Algorithm: algorithm,
	})
}

func (b *Blueprint) createIndexName(typ string, columns []string) string {
	index := strings.ToLower(b.prefix + b.table + "_" + strings.Join(columns, "_") + "_" + typ)
	index = strings.Replace(index, "-", "_", -1)
	index = strings.Replace(index, ".", "_", -1)

	return index
}

func (b *Blueprint) addColumn(typ, name string, options *ColumnOptions) *ColumnDefinition {
	definition := &ColumnDefinition{
		Typ:  typ,
		Name: name,
	}

	if options != nil {
		if options.Length > 0 {
			definition.Length = options.Length
		}
		if options.autoIncrement {
			definition.autoIncrement = true
		}
		if options.unsigned {
			definition.unsigned = true
		}
		if options.Total > 0 {
			definition.Total = options.Total
		}
		if options.Places > 0 {
			definition.Places = options.Places
		}

		if len(options.Allowed) > 0 {
			definition.Allowed = options.Allowed
		}

		if options.Precision > 0 {
			definition.Precision = options.Precision
		}
	}

	b.columns = append(b.columns, definition)

	return definition
}

// Create indicate that the table needs to be created.
func (b *Blueprint) Create() {
	b.addCommand("create", nil)
}

// Temporary indicate that the table needs to be temporary.
func (b *Blueprint) Temporary() {
	b.temporary = true
}

// Drop indicate that the table should be dropped.
func (b *Blueprint) Drop() {
	b.addCommand("drop", nil)
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
		autoIncrement: autoIncrement,
		unsigned:      unsigned,
	})
}

// TinyInteger create a new tiny integer (1-byte) column on the table
func (b *Blueprint) TinyInteger(column string, autoIncrement, unsigned bool) *ColumnDefinition {
	return b.addColumn("tinyInteger", column, &ColumnOptions{
		autoIncrement: autoIncrement,
		unsigned:      unsigned,
	})
}

// SmallInteger create a new small integer (2-byte) column on the table
func (b *Blueprint) SmallInteger(column string, autoIncrement, unsigned bool) *ColumnDefinition {
	return b.addColumn("smallInteger", column, &ColumnOptions{
		autoIncrement: autoIncrement,
		unsigned:      unsigned,
	})
}

// MediumInteger create a new mediumInteger (3-byte) column on the table
func (b *Blueprint) MediumInteger(column string, autoIncrement, unsigned bool) *ColumnDefinition {
	return b.addColumn("mediumInteger", column, &ColumnOptions{
		autoIncrement: autoIncrement,
		unsigned:      unsigned,
	})
}

// BigInteger create a new bigInteger (8-byte) column on the table
func (b *Blueprint) BigInteger(column string, autoIncrement, unsigned bool) *ColumnDefinition {
	return b.addColumn("bigInteger", column, &ColumnOptions{
		autoIncrement: autoIncrement,
		unsigned:      unsigned,
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

// Float create a new float column on the table.
func (b *Blueprint) Float(column string, total, places int, unsigned bool) *ColumnDefinition {
	if total == 0 {
		total = 8
	}
	if places == 0 {
		places = 2
	}

	return b.addColumn("float", column, &ColumnOptions{
		Total:    total,
		Places:   places,
		unsigned: unsigned,
	})
}

// Double create a new Double column on the table.
func (b *Blueprint) Double(column string, total, places int, unsigned bool) *ColumnDefinition {
	return b.addColumn("double", column, &ColumnOptions{
		Total:    total,
		Places:   places,
		unsigned: unsigned,
	})
}

// Decimal create a new Double column on the table.
func (b *Blueprint) Decimal(column string, total, places int, unsigned bool) *ColumnDefinition {
	if total == 0 {
		total = 8
	}
	if places == 0 {
		places = 2
	}
	return b.addColumn("decimal", column, &ColumnOptions{
		Total:    total,
		Places:   places,
		unsigned: unsigned,
	})
}

// UnsignedFloat create a new funsigned loat column on the table.
func (b *Blueprint) UnsignedFloat(column string, total, places int, unsigned bool) *ColumnDefinition {
	return b.Float(column, total, places, true)
}

// UnsignedDouble create a new unsigned double column on the table.
func (b *Blueprint) UnsignedDouble(column string, total, places int, unsigned bool) *ColumnDefinition {
	return b.Double(column, total, places, true)
}

// UnsignedDecimal create a new unsigned decimal column on the table.
func (b *Blueprint) UnsignedDecimal(column string, total, places int, unsigned bool) *ColumnDefinition {
	return b.Decimal(column, total, places, true)
}

// Boolean create a new boolean column on the table
func (b *Blueprint) Boolean(column string, autoIncrement, unsigned bool) *ColumnDefinition {
	return b.addColumn("boolean", column, nil)
}

// Enum create a new enum column on the table
func (b *Blueprint) Enum(column string, allowed []string) *ColumnDefinition {
	return b.addColumn("enum", column, &ColumnOptions{
		Allowed: allowed,
	})
}

// Set create a new set column on the table
func (b *Blueprint) Set(column string, allowed []string) *ColumnDefinition {
	return b.addColumn("set", column, &ColumnOptions{
		Allowed: allowed,
	})
}

// JSON create a new json column on the table
func (b *Blueprint) JSON(column string) *ColumnDefinition {
	return b.addColumn("json", column, nil)
}

// JSONB create a new jsonb column on the table
func (b *Blueprint) JSONB(column string) *ColumnDefinition {
	return b.addColumn("jsonb", column, nil)
}

// Date create a new Date column on the table
func (b *Blueprint) Date(column string) *ColumnDefinition {
	return b.addColumn("date", column, nil)
}

// DateTime create a new dateTime column on the table
func (b *Blueprint) DateTime(column string, precision int) *ColumnDefinition {
	return b.addColumn("dateTime", column, &ColumnOptions{
		Precision: precision,
	})
}

// DateTimeTz create a new date-time column (with time zone) on the table
func (b *Blueprint) DateTimeTz(column string, precision int) *ColumnDefinition {
	return b.addColumn("dateTimeTz", column, &ColumnOptions{
		Precision: precision,
	})
}

// Time create a new time column on the table.
func (b *Blueprint) Time(column string, precision int) *ColumnDefinition {
	return b.addColumn("time", column, &ColumnOptions{
		Precision: precision,
	})
}

// TimeTz create a new time column (with time zone) on the table.
func (b *Blueprint) TimeTz(column string, precision int) *ColumnDefinition {
	return b.addColumn("timeTz", column, &ColumnOptions{
		Precision: precision,
	})
}

// Timestamp create a new timestamp column on the table.
func (b *Blueprint) Timestamp(column string, precision int) *ColumnDefinition {
	return b.addColumn("timestamp", column, &ColumnOptions{
		Precision: precision,
	})
}

// TimestampTz create a new timestamp (with time zone) column on the table.
func (b *Blueprint) TimestampTz(column string, precision int) *ColumnDefinition {
	return b.addColumn("timestampTz", column, &ColumnOptions{
		Precision: precision,
	})
}

// Timestamps add nullable creation and update timestamps to the table.
func (b *Blueprint) Timestamps(precision int) {
	b.Timestamp("created_at", precision).Nullable()
	b.Timestamp("updated_at", precision).Nullable()
}

// TimestampsTz add nullable creation and update timestampTz columns to the table.
func (b *Blueprint) TimestampsTz(precision int) {
	b.TimestampTz("created_at", precision).Nullable()
	b.TimestampTz("updated_at", precision).Nullable()
}

// SoftDeletes add a "deleted at" timestamp for the table.
func (b *Blueprint) SoftDeletes(precision int) *ColumnDefinition {
	return b.Timestamp("deleted_at", precision).Nullable()
}

// SoftDeletesTz add a "deleted at" timestamp (with time zone) for the table.
func (b *Blueprint) SoftDeletesTz(precision int) *ColumnDefinition {
	return b.TimestampTz("deleted_at", precision).Nullable()
}

// Year create a new year column on the table.
func (b *Blueprint) Year(column string) *ColumnDefinition {
	return b.addColumn("year", column, nil)
}

// Binary create a new binary column on the table.
func (b *Blueprint) Binary(column string) *ColumnDefinition {
	return b.addColumn("binary", column, nil)
}

// UUID create a new UUID column on the table.
func (b *Blueprint) UUID(column string) *ColumnDefinition {
	return b.addColumn("uuid", column, nil)
}

// IPAddress create a new IPAddress column on the table.
func (b *Blueprint) IPAddress(column string) *ColumnDefinition {
	return b.addColumn("ipAddress", column, nil)
}

// MacAddress create a new MacAddress column on the table.
func (b *Blueprint) MacAddress(column string) *ColumnDefinition {
	return b.addColumn("macAddress", column, nil)
}

// Geometry create a new Geometry column on the table.
func (b *Blueprint) Geometry(column string) *ColumnDefinition {
	return b.addColumn("geometry", column, nil)
}

// Point create a new point column on the table.
func (b *Blueprint) Point(column string) *ColumnDefinition {
	return b.addColumn("point", column, nil)
}

// LineString create a new linestring column on the table.
func (b *Blueprint) LineString(column string) *ColumnDefinition {
	return b.addColumn("linestring", column, nil)
}

// Polygon create a new polygon column on the table.
func (b *Blueprint) Polygon(column string) *ColumnDefinition {
	return b.addColumn("polygon", column, nil)
}

// GeometryCollection create a new geometrycollection column on the table.
func (b *Blueprint) GeometryCollection(column string) *ColumnDefinition {
	return b.addColumn("geometrycollection", column, nil)
}

// MultiPoint create a new multiPoint column on the table.
func (b *Blueprint) MultiPoint(column string) *ColumnDefinition {
	return b.addColumn("multiPoint", column, nil)
}

// MultiLineString create a new multilinestring column on the table.
func (b *Blueprint) MultiLineString(column string) *ColumnDefinition {
	return b.addColumn("multilinestring", column, nil)
}

// MultiPolygon create a new multipolygon column on the table.
func (b *Blueprint) MultiPolygon(column string) *ColumnDefinition {
	return b.addColumn("multipolygon", column, nil)
}

// MultiPolygonZ create a new multipolygonz column on the table.
func (b *Blueprint) MultiPolygonZ(column string) *ColumnDefinition {
	return b.addColumn("multipolygonz", column, nil)
}

// Morphs add the proper columns for a polymorphic table.
func (b *Blueprint) Morphs(column string, indexName string) {
	b.String(column+"_type", 0)
	b.UnsignedBigInteger(column+"_id", false)
	b.Index([]string{column + "_type", column + "_id"}, indexName, "")
}

// NullableMorphs add nullable columns for a polymorphic table.
func (b *Blueprint) NullableMorphs(column string, indexName string) {
	b.String(column+"_type", 0).Nullable()
	b.UnsignedBigInteger(column+"_id", false).Nullable()
	b.Index([]string{column + "_type", column + "_id"}, indexName, "")
}

// UUIDMorphs add the proper columns for a polymorphic table using UUIDs.
func (b *Blueprint) UUIDMorphs(column string, indexName string) {
	b.String(column+"_type", 0)
	b.UUID(column + "_id")
	b.Index([]string{column + "_type", column + "_id"}, indexName, "")
}

// NullableUUIDMorphs add nullable uuid columns for a polymorphic table.
func (b *Blueprint) NullableUUIDMorphs(column string, indexName string) {
	b.String(column+"_type", 0).Nullable()
	b.UUID(column + "_id").Nullable()
	b.Index([]string{column + "_type", column + "_id"}, indexName, "")
}

// RememberToken adds the `remember_token` column to the table.
func (b *Blueprint) RememberToken() *ColumnDefinition {
	return b.String("remember_token", 100).Nullable()
}

// ID Create a new auto-incrementing big integer (8-byte) column on the table.
func (b *Blueprint) ID(column string) *ColumnDefinition {
	if column == "" {
		column = "id"
	}
	return b.BigIncrements(column)
}

// BigIncrements Create a new auto-incrementing big integer (8-byte) column on the table.
func (b *Blueprint) BigIncrements(column string) *ColumnDefinition {
	return b.UnsignedBigInteger(column, true)
}

// ForeignID Create a new unsigned big integer (8-byte) column on the table.
func (b *Blueprint) ForeignID(column string) *ColumnDefinition {
	return b.addColumn("bigInteger", column, &ColumnOptions{
		autoIncrement: true,
		unsigned:      true,
	})
}

// foreignID

// Index Specify an index for the table.
func (b *Blueprint) Index(columns []string, name string, algorithm string) *Blueprint {
	return b.indexCommand("index", columns, name, algorithm)
}

// SpatialIndex Specify a spatial index for the table.
func (b *Blueprint) SpatialIndex(columns []string, name string) *Blueprint {
	return b.indexCommand("spatialIndex", columns, name, "")
}

// Rename the table to a given name.
func (b *Blueprint) Rename(to string) *Blueprint {
	return b.addCommand("rename", &CommandOptions{
		To: to,
	})
}

// Primary Specify the primary key(s) for the table.
func (b *Blueprint) Primary(columns []string, name string, algorithm string) *Blueprint {
	return b.indexCommand("primary", columns, name, algorithm)
}

// Unique Specify a unique index for the table.
func (b *Blueprint) Unique(columns []string, name string, algorithm string) *Blueprint {
	return b.indexCommand("unique", columns, name, algorithm)
}

// Foreign Specify a foreign key for the table.
func (b *Blueprint) Foreign(columns []string, name string) *Blueprint {
	return b.indexCommand("foreign", columns, name, "")
}

func (b *Blueprint) toSQL(conn *Connection, grammar Grammar) []string {
	var statements []string

	for _, cmd := range b.commands {
		switch cmd.Typ {
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
