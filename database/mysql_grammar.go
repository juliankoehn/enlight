package database

import (
	"bytes"
	"strconv"
)

type (
	MySQLGrammar struct {
		*grammar
		modifiers []string
		serials   []string
	}
)

var (
	mysqlDefaultModifiers = []string{
		"Unsigned", "Charset", "Collate", "VirtualAs", "StoredAs",
		"Nullable", "Srid", "Default", "Increment", "Comment", "After", "First",
	}
	mysqlDefaultSerials = []string{
		"bigInteger", "integer", "mediumInteger", "smallInteger", "tinyInteger",
	}
)

// NewMySQLGrammar returns a new mysql grammar instance
func NewMySQLGrammar() *MySQLGrammar {
	grammar := NewGrammar()

	return &MySQLGrammar{
		grammar:   grammar,
		modifiers: mysqlDefaultModifiers,
		serials:   mysqlDefaultSerials,
	}
}

// CompileTableExists Compile the query to determine the list of tables.
func (mysql *MySQLGrammar) CompileTableExists() string {
	// return "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ?  AND table_name = ? and table_type = 'BASE TABLE';"
	return "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ?  AND table_name = ? and table_type = 'BASE TABLE';"
}

// CompileColumnListing Compile the query to determine the list of columns.
// @args: databaseName, tableName
func (mysql *MySQLGrammar) CompileColumnListing() string {
	return "select column_name as `column_name` from information_schema.columns where table_schema = ? and table_name = ?"
}

// ProcessColumnListing Process the results of a column listing query.
func (mysql *MySQLGrammar) ProcessColumnListing() {

}

// CompileCreate Compile a create table command.
func (mysql *MySQLGrammar) CompileCreate(bp *Blueprint, conn *Connection) string {
	sql := mysql.compileCreateTable(bp)

	// once we have the primary SQL, we can add the encoding option to the SQL for
	// the table. Then we can check if a storage engine has been supplied for
	// the table. If so we will add the engine declration to the SQL query.
	sql = sql + mysql.compileCreateEncoding(bp, conn)

	return sql
}

func (mysql *MySQLGrammar) compileCreateEncoding(bp *Blueprint, conn *Connection) string {
	var buf bytes.Buffer

	// First we will set the character set if one has been set on either the create
	// blueprint itself or on the root configuration for the connection that the
	// table is being created on. We will add these to the create table query.
	if bp.GetCharset() != "" {
		buf.WriteString(" default character set ")
		buf.WriteString(bp.GetCharset())
	} else if conn.config.Charset != "" {
		buf.WriteString(" default character set ")
		buf.WriteString(conn.config.Charset)
	}

	// Next we will add the collation to the create table statement if one has been
	// added to either this create table blueprint or the configuration for this
	// connection that the query is targeting. We'll add it to this SQL query.
	if bp.GetCollation() != "" {
		buf.WriteString(" collate ")
		buf.WriteString(bp.GetCollation())
	} else if conn.config.Charset != "" {
		buf.WriteString(" collate ")
		buf.WriteString(conn.config.Charset)
	}

	return buf.String()
}

func (mysql *MySQLGrammar) compileCreateTable(bp *Blueprint) string {
	var buf bytes.Buffer
	command := "create"

	if bp.temporary {
		command = "create temporary"
	}

	buf.WriteString(command)
	buf.WriteString(" table ")
	buf.WriteString(bp.table)
	buf.WriteByte('(')

	columns := mysql.getColumns(bp)

	for key, col := range columns {
		buf.WriteString(col)
		if key != len(columns)-1 {
			buf.WriteString(",")
		}
	}

	buf.WriteByte(')')

	return buf.String()
}

func (mysql *MySQLGrammar) getColumns(bp *Blueprint) []string {
	columnsToAdd := bp.GetAddedColumns()
	columns := make([]string, len(columnsToAdd))

	for key, column := range columnsToAdd {
		var buf bytes.Buffer

		buf.WriteString(column.Name)
		buf.WriteByte(' ')
		switch column.Typ {
		case "char":
			buf.WriteString(mysql.typeChar(column))
		case "string":
			buf.WriteString(mysql.typeString(column))
		case "text":
			buf.WriteString(mysql.typeText(column))
		case "mediumText":
			buf.WriteString(mysql.typeText(column))
		case "longText":
			buf.WriteString(mysql.typeLongText(column))
		case "bigInteger":
			buf.WriteString(mysql.typeBigInteger(column))
		case "integer":
			buf.WriteString(mysql.typeInteger(column))
		case "mediumInteger":
			buf.WriteString(mysql.typeMediumInteger(column))
		case "tinyInteger":
			buf.WriteString(mysql.typeTinyInteger(column))
		case "smallInteger":
			buf.WriteString(mysql.typeSmallInteger(column))
		case "float":
			buf.WriteString(mysql.typeFloat(column))
		case "double":
			buf.WriteString(mysql.typeFloat(column))
		case "decimal":
			buf.WriteString(mysql.typeFloat(column))
		case "boolean":
			buf.WriteString(mysql.typeBoolean(column))
		case "enum":
			buf.WriteString(mysql.typeEnum(column))
		case "set":
			buf.WriteString(mysql.typeSet(column))
		case "json":
			buf.WriteString(mysql.typeJSON(column))
		case "jsonb":
			buf.WriteString(mysql.typeJSONB(column))
		case "date":
			buf.WriteString(mysql.typeDate(column))
		case "dateTime":
			buf.WriteString(mysql.typeDateTime(column))
		case "dateTimeTz":
			buf.WriteString(mysql.typeDateTimeTz(column))
		case "time":
			buf.WriteString(mysql.typeTime(column))
		case "timeTz":
			buf.WriteString(mysql.typeTime(column))
		case "timestamp":
			buf.WriteString(mysql.typeTimestamp(column))
		case "timestampTz":
			buf.WriteString(mysql.typeTimestampTz(column))
		case "year":
			buf.WriteString(mysql.typeYear(column))
		case "binary":
			buf.WriteString(mysql.typeBinary(column))
		case "uuid":
			buf.WriteString(mysql.typeUUID(column))
		case "ipAddress":
			buf.WriteString(mysql.typeIPAddress(column))
		case "macAddress":
			buf.WriteString(mysql.typeMacAddress(column))
		case "geometry":
			buf.WriteString(mysql.typeGeometry(column))
		case "point":
			buf.WriteString(mysql.typePoint(column))
		case "lineString":
			buf.WriteString(mysql.typeLineString(column))
		case "polygon":
			buf.WriteString(mysql.typePolygon(column))
		case "geometryCollection":
			buf.WriteString(mysql.typeGeometryCollection(column))
		case "multiPoint":
			buf.WriteString(mysql.typeMultiPoint(column))
		case "multiLineString":
			buf.WriteString(mysql.typeMultiLineString(column))
		case "multiPolygon":
			buf.WriteString(mysql.typeMultiPolygon(column))
		}

		columns[key] = buf.String()
	}

	return columns
}

// Get the SQL for a generated virtual column modifier.
func (mysql *MySQLGrammar) modifyVirtualAs(column *ColumnDefinition) string {
	if column.VirtualAs != "" {
		return " as (" + column.VirtualAs + ")"
	}

	return ""
}

// Create the column definition for a spatial MultiLineString type.
func (mysql *MySQLGrammar) typeMultiPolygon(column *ColumnDefinition) string {
	return "multipolygon"
}

// Create the column definition for a spatial MultiLineString type.
func (mysql *MySQLGrammar) typeMultiLineString(column *ColumnDefinition) string {
	return "multilinestring"
}

// Create the column definition for a spatial MultiPoint type.
func (mysql *MySQLGrammar) typeMultiPoint(column *ColumnDefinition) string {
	return "multipoint"
}

// Create the column definition for a spatial GeometryCollection type.
func (mysql *MySQLGrammar) typeGeometryCollection(column *ColumnDefinition) string {
	return "geometrycollection"
}

// Create the column definition for a spatial Polygon type.
func (mysql *MySQLGrammar) typePolygon(column *ColumnDefinition) string {
	return "polygon"
}

// Create the column definition for a spatial LineString type.
func (mysql *MySQLGrammar) typeLineString(column *ColumnDefinition) string {
	return "linestring"
}

// Create the column definition for a spatial Point type.
func (mysql *MySQLGrammar) typePoint(column *ColumnDefinition) string {
	return "point"
}

// Create the column definition for a spatial Geometry type.
func (mysql *MySQLGrammar) typeGeometry(column *ColumnDefinition) string {
	return "geometry"
}

// Create the column definition for a MAC address type.
func (mysql *MySQLGrammar) typeMacAddress(column *ColumnDefinition) string {
	return "varchar(17)"
}

// Create the column definition for an IP address type.
func (mysql *MySQLGrammar) typeIPAddress(column *ColumnDefinition) string {
	return "varchar(45)"
}

// Create the column definition for a uuid type.
func (mysql *MySQLGrammar) typeUUID(column *ColumnDefinition) string {
	return "char(36)"
}

// Create the column definition for a binary type.
func (mysql *MySQLGrammar) typeBinary(column *ColumnDefinition) string {
	return "blob"
}

// Create the column definition for a year type.
func (mysql *MySQLGrammar) typeYear(column *ColumnDefinition) string {
	return "year"
}

// Create the column definition for a time (with time zone) type.
func (mysql *MySQLGrammar) typeTimestampTz(column *ColumnDefinition) string {
	return mysql.typeTimestamp(column)
}

// Create the column definition for a timestamp type.
func (mysql *MySQLGrammar) typeTimestamp(column *ColumnDefinition) string {
	columnType := "timestamp"
	if column.Precision > 0 {
		columnType = "timestamp(" + strconv.Itoa(column.Precision) + ")"
	}

	if column.UseCurrent {
		useCurrent := columnType + " default CURRENT_TIMESTAMP"
		if column.Precision > 0 {
			useCurrent = columnType + " default CURRENT_TIMESTAMP(" + strconv.Itoa(column.Precision) + ")"
		}
		return useCurrent
	}
	return columnType
}

// Create the column definition for a time (with time zone) type.
func (mysql *MySQLGrammar) typeTimeTz(column *ColumnDefinition) string {
	return mysql.typeTime(column)
}

// Create the column definition for a time type.
func (mysql *MySQLGrammar) typeTime(column *ColumnDefinition) string {
	if column.Precision > 0 {
		return "time(" + strconv.Itoa(column.Precision) + ")"
	}
	return "time"
}

// Create the column definition for a date-time (with time zone) type.
func (mysql *MySQLGrammar) typeDateTimeTz(column *ColumnDefinition) string {
	return mysql.typeDateTime(column)
}

// Create the column definition for a date-time type.
func (mysql *MySQLGrammar) typeDateTime(column *ColumnDefinition) string {
	columnType := "datetime"
	if column.Precision > 0 {
		columnType = "datetime(" + strconv.Itoa(column.Precision) + ")"
	}

	if column.UseCurrent {
		return columnType + " default CURRENT_TIMESTAMP"
	}
	return columnType
}

// Create the column definition for a date type.
func (mysql *MySQLGrammar) typeDate(column *ColumnDefinition) string {
	return "date"
}

// Create the column definition for a jsonb type.
func (mysql *MySQLGrammar) typeJSONB(column *ColumnDefinition) string {
	return "json"
}

// Create the column definition for a json type.
func (mysql *MySQLGrammar) typeJSON(column *ColumnDefinition) string {
	return "json"
}

// Create the column definition for a set enumeration type.
func (mysql *MySQLGrammar) typeSet(column *ColumnDefinition) string {
	return "set(" + mysql.QuoteString(column.Allowed) + ")"
}

// Create the column definition for a enumeration type.
func (mysql *MySQLGrammar) typeEnum(column *ColumnDefinition) string {
	return "enum(" + mysql.QuoteString(column.Allowed) + ")"
}

// Create the column definition for a Boolean type.
func (mysql *MySQLGrammar) typeBoolean(column *ColumnDefinition) string {
	return "tinyint(1)"
}

// Create the column definition for a decimal type.
func (mysql *MySQLGrammar) typeDecimal(column *ColumnDefinition) string {
	return "decimal(" + strconv.Itoa(column.Total) + "," + strconv.Itoa(column.Places) + ")"
}

// Create the column definition for a double type.
// total int, places int
func (mysql *MySQLGrammar) typeDouble(column *ColumnDefinition) string {
	if column.Total > 0 && column.Places > 0 {
		return "double(" + strconv.Itoa(column.Total) + "," + strconv.Itoa(column.Places) + ")"
	}
	return "double"
}

// Create the column definition for a float type.
func (mysql *MySQLGrammar) typeFloat(column *ColumnDefinition) string {
	return mysql.typeDouble(column)
}

// Create the column definition for a smallint type.
func (mysql *MySQLGrammar) typeSmallInteger(column *ColumnDefinition) string {
	return "smallint"
}

// Create the column definition for a tinyint type.
func (mysql *MySQLGrammar) typeTinyInteger(column *ColumnDefinition) string {
	return "tinyint"
}

// Create the column definition for a mediumint type.
func (mysql *MySQLGrammar) typeMediumInteger(column *ColumnDefinition) string {
	return "mediumint"
}

// Create the column definition for a char type.
func (mysql *MySQLGrammar) typeChar(column *ColumnDefinition) string {
	return "char(" + column.Length + ")"
}

// Create the column definition for a string type.
func (mysql *MySQLGrammar) typeString(column *ColumnDefinition) string {
	return "varchar(" + column.Length + ")"
}

// Create the column definition for a text type.
func (mysql *MySQLGrammar) typeText(column *ColumnDefinition) string {
	return "text"
}

// Create the column definition for a medium text type.
func (mysql *MySQLGrammar) typeMediumText(column *ColumnDefinition) string {
	return "mediumtext"
}

// Create the column definition for a long text type.
func (mysql *MySQLGrammar) typeLongText(column *ColumnDefinition) string {
	return "longtext"
}

// Create the column definition for a big integer type.
func (mysql *MySQLGrammar) typeBigInteger(column *ColumnDefinition) string {
	return "bigint"
}

// Create the column definition for a integer type.
func (mysql *MySQLGrammar) typeInteger(column *ColumnDefinition) string {
	return "int"
}
