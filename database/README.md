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
| table.Char("name", 100) | CHAR equivalent column with a length. |
| table.Integer("votes", false, false) | INTEGER equivalent column. (autoIncrement, Unsigned) |
| table.LongText("description") | LONGTEXT equivalent column. |
| table.MediumInteger("votes", false, false) | MEDIUMINT equivalent column. (autoIncrement, Unsigned) |
| table.MediumText("description") | MEDIUMTEXT equivalent column. |
| table.SmallInteger("votes", false, false) | SMALLINT equivalent column. (autoIncrement, Unsigned) |
| table.String("name", 100) | VARCHAR equivalent column with a length. |
| table.Text("description") | TEXT equivalent column. |
| table.TinyInteger("votes", false, false) | TINYINT equivalent column. (autoIncrement, Unsigned) |
| table.UnsignedBigInteger("votes", false) | UNSIGNED BIGINT equivalent column. (autoIncrement) |
| table.UnsignedInteger("votes", false) | UNSIGNED INTEGER equivalent column. (autoIncrement) |
| table.UnsignedMediumInteger("votes", false) | UNSIGNED MEDIUMTEXT equivalent column. (autoIncrement) |
| table.UnsignedSmallInteger("votes", false) | UNSIGNED SMALLINT equivalent column. (autoIncrement) |
| table.UnsignedTinyInteger("votes", false) | UNSIGNED TINYINT equivalent column. (autoIncrement) |