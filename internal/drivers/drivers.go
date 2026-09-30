// Package drivers registers the database/sql drivers this application opens.
//
// The platform speaks to database/sql through sqlx and registers nothing, so
// the driver set is the choice of whoever builds the binary. Blank import this
// package from a main or a test that opens a DSN.
//
// Only the schemes platform-app configures are here: mysql:// and sqlite://.
// A postgres:// DSN needs github.com/jackc/pgx/v5/stdlib added below.
package drivers

import (
	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
)
