package database

import (
	"database/sql"

	"github.com/jmoiron/sqlx/reflectx"
)

type (
	// Stmt is an sqlx wrapper around sql.Stmt with extra functionality
	Stmt struct {
		*sql.Stmt
		unsafe bool
		Mapper *reflectx.Mapper
	}

	// qStmt is an unexposed wrapper which lets you use a Stmt as a Queryer & Execer by
	// implementing those interfaces and ignoring the `query` argument.
	qStmt struct{ *Stmt }
)

// Unsafe returns a version of Stmt which will silently succeed to scan when
// columns in the SQL result have no fields in the destination struct.
func (s *Stmt) Unsafe() *Stmt {
	return &Stmt{Stmt: s.Stmt, unsafe: true, Mapper: s.Mapper}
}

// Select using the prepared statement.
// Any placeholder parameters are replaced with supplied args.
func (s *Stmt) Select(dest interface{}, args ...interface{}) error {
	return Select(&qStmt{s}, dest, "", args...)
}

func (q *qStmt) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return q.Stmt.Query(args...)
}
