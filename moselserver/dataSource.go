package moselserver

import (
	"database/sql"
	"gopkg.in/mgo.v2"
)

type dataSource interface {
	GetType() string
}

type dataSourceImpl struct {
	typev string
}

func (ds *dataSourceImpl) GetType() string {
	return ds.typev
}

// SQL data source

type SqlDataSource interface {
	dataSource
	GetDb() *sql.DB
}

type sqlDataSourceImpl struct {
	dataSourceImpl
	db *sql.DB
}

func NewSqlDataSource(t string, db *sql.DB) SqlDataSource {
	ds := sqlDataSourceImpl{db: db}
	ds.typev = t
	return &ds
}

func (ds *sqlDataSourceImpl) GetDb() *sql.DB {
	return ds.db
}

// mongo data source

type MongoDataSource interface {
	dataSource
	GetSession() *mgo.Session
}

type mongoDataSourceImpl struct {
	dataSourceImpl
	session *mgo.Session
}

func NewMongoDataSource(t string, session *mgo.Session) MongoDataSource {
	ds := mongoDataSourceImpl{session: session}
	ds.typev = t
	return &ds
}

func (ds *mongoDataSourceImpl) GetSession() *mgo.Session {
	return ds.session
}
