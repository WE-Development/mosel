package context

import (
	"github.com/WE-Development/mosel/commons"
	"time"
	"github.com/WE-Development/mosel/api"
	"database/sql"
)

type dataPersistence interface {
	Add(node string, t time.Time, info api.NodeInfo)
}

type sqlDataPersistence struct {
	db      *sql.DB
	queries commons.SqlQueries
}

func NewSqlDataPersistence(db *sql.DB, queries commons.SqlQueries) dataPersistence {
	return sqlDataPersistence{
		db:db,
		queries:queries,
	}
}

func (pers sqlDataPersistence) Add(node string, t time.Time, info api.NodeInfo) {

}