package moselserver

import "database/sql"

type dataSource struct {
	Type string
	Db   *sql.DB
}