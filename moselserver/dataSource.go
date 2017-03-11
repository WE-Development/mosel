package moselserver

import "database/sql"

type dataSource interface {
	GetType() string
}

type SqlDataSource interface {
	dataSource
	GetDb() *sql.DB
}

type sqlDataSourceImpl struct {
	typev string
	db    *sql.DB
}

func NewSqlDataSource(t string, db *sql.DB) *sqlDataSourceImpl {
	return &sqlDataSourceImpl{
		typev:t,
		db:db,
	}
}

func (ds *sqlDataSourceImpl) GetType() string {
	return ds.typev
}

func (ds *sqlDataSourceImpl) GetDb() *sql.DB {
	return ds.db
}