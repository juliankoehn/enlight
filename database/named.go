package database

// NamedStmt is a prepared statement that executes named queries.  Prepare it
// how you would execute a NamedQuery, but pass in a struct or map when executing.
type NamedStmt struct {
	Params      []string
	QueryString string
	Stmt        *Stmt
}
