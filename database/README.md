# enlight/database

Database is a multi-tenancy database connector using `database/sql`

## Package is under heavy-developement. There is only `mysql` implemented yet.

## Supported drivers

| Database | driver |
|---|---|
| MSSQL | mssql, go-mssqldb |
| MySQL (mysql) | my,mariadb,maria,percona,aurora |
| PostgreSQL (postgres) | pg, postgresql, pgsql |
| SQLite3 (sqlite3) | sq, sqlite, file |

See https://github.com/golang/go/wiki/SQLDrivers for a list of 3rd Party drivers

## Available Column Types

| Command | Description |
|---|---|
| table.BigInteger("votes", false, false) | BIGINT equivalent column. (autoIncrement, Unsigned) |
| table.Binary("data") | BLOB equivalent column. |
| table.Boolean("confirmed") | BOOLEAN equivalent column. |
| table.Char("name", 100) | CHAR equivalent column with a length. |
| table.Date("created_at") | DATE equivalent column. |
| table.DateTime("created_at", 0) | DATETIME equivalent column with precision (total digits). |
| table.DateTimeTz("created_at", 0) | DATETIME (with timezone) equivalent column with precision (total digits). |
| table.Decimal("amount", 8, 2) | DECIMAL equivalent column with precision (total digits) and scale (decimal digits).|
| table.Double("amount", 8, 2) | DOUBLE equivalent column with precision (total digits) and scale (decimal digits). |
| table.Enum("level", ["easy", "hard"]) | ENUM equivalent column. |
| table.Float("amount", 8, 2) | FLOAT equivalent column with a precision (total digits) and scale (decimal digits). |
| table.Geometry("positions") | GEOMETRY equivalent column. |
| table.GeometryCollection("positions") | GEOMETRYCOLLECTION equivalent column. |
| table.Integer("votes", false, false) | INTEGER equivalent column. (autoIncrement, Unsigned) |
| table.IPAddress("visitor") | IP address equivalent column. |
| table.JSON("options") | JSON equivalent column. |
| table.JSONB("options") | JSONB equivalent column. |
| table.LineString("positions") | LINESTRING equivalent column. |
| table.LongText("description") | LONGTEXT equivalent column. |
| table.MacAddress("device") | MAC address equivalent column. |
| table.MediumInteger("votes", false, false) | MEDIUMINT equivalent column. (autoIncrement, Unsigned) |
| table.MediumText("description") | MEDIUMTEXT equivalent column. |
| table.Morphs("taggable") | Adds `taggable_id` UNSIGNED BIGINT and `taggable_type` VARCHAR equivalent columns. |
| table.NullableMorphs("taggable") | Adds nullable versions of `Morphs()` columns. |
| table.MultiLineString("positions") | MULTILINESTRING equivalent column. |
| table.MultiPoint("positions") | MULTIPOINT equivalent column. |
| table.MultiPolygon("positions") | MULTIPOLYGON equivalent column. |
| table.Point("position") | POINT equivalent column. |
| table.Polygon("positions") | POLYGON equivalent column. |
| table.RememberToken() | Adds a nullable `remember_token` VARCHAR(100) equivalent column. |
| table.Set("flavors", ["strawberry", "vanilla"]) | SET equivalent column. |
| table.SmallInteger("votes", false, false) | SMALLINT equivalent column. (autoIncrement, Unsigned) |
| table.SoftDeletes(0) | Adds a nullable `deleted_at` TIMESTAMP equivalent column for soft deletes with precision (total digits). |
| table.SoftDeletesTz(0) | 	Adds a nullable `deleted_at` TIMESTAMP (with timezone) equivalent column for soft deletes with precision (total digits). |
| table.String("name", 100) | VARCHAR equivalent column with a length. |
| table.Text("description") | TEXT equivalent column. |
| table.Time("sunrise", 0) | TIME equivalent column with precision (total digits). |
| table.TimeTz("sunrise", 0) | TIME (with timezone) equivalent column with precision (total digits). |
| table.Timestamp("added_on", 0) | TIMESTAMP equivalent column with precision (total digits). |
| table.TimestampTz("added_on", 0) | TIMESTAMP (with timezone) equivalent column with precision (total digits). |
| table.Timestamps(0) | Adds nullable `created_at` and `updated_at` TIMESTAMP equivalent columns with precision (total digits). |
| table.TimestampsTz(0) | Adds nullable `created_at` and `updated_at` TIMESTAMP (with timezone) equivalent columns with precision (total digits). |
| table.TinyInteger("votes", false, false) | TINYINT equivalent column. (autoIncrement, Unsigned) |
| table.UnsignedBigInteger("votes", false) | UNSIGNED BIGINT equivalent column. (autoIncrement) |
| table.UnsignedDecimal("votes", 8, 2) | UNSIGNED DECIMAL equivalent column with a precision (total digits) and scale (decimal digits). |
| table.UnsignedDouble("votes", 8, 2) | UNSIGNED DOUBLE equivalent column with a precision (total digits) and scale (decimal digits). |
| table.UnsignedFloat("votes", 8, 2) | UNSIGNED FLOAT equivalent column with a precision (total digits) and scale (decimal digits). |
| table.UnsignedInteger("votes", false) | UNSIGNED INTEGER equivalent column. (autoIncrement) |
| table.UnsignedMediumInteger("votes", false) | UNSIGNED MEDIUMTEXT equivalent column. (autoIncrement) |
| table.UnsignedSmallInteger("votes", false) | UNSIGNED SMALLINT equivalent column. (autoIncrement) |
| table.UnsignedTinyInteger("votes", false) | UNSIGNED TINYINT equivalent column. (autoIncrement) |
| table.UUID("id") | UUID equivalent column. |
| table.Year("birth_year") | YEAR equivalent column. |