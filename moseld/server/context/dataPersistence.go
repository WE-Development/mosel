package context

import (
	"github.com/WE-Development/mosel/commons"
	"time"
	"github.com/WE-Development/mosel/api"
	"database/sql"
	"fmt"
	"log"
)

type dataPersistence interface {
	Init() error
	Add(node string, t time.Time, info api.NodeInfo)
}

type sqlDataPersistence struct {
	db *sql.DB
	q  commons.SqlQueries
}

func NewSqlDataPersistence(db *sql.DB, queries commons.SqlQueries) dataPersistence {
	return sqlDataPersistence{
		db:db,
		q:queries,
	}
}

func (pers sqlDataPersistence) query(name string, args ...interface{}) (*sql.Rows, error) {
	query, exists := pers.q[name]

	if !exists {
		return nil, fmt.Errorf("Quers %s is not registered", name)
	}

	return pers.db.Query(query, args...)
}

func (pers sqlDataPersistence) Init() error {
	tables := make(map[string]string)
	tables["Nodes"] = "createNodes"
	tables["Diagrams"] = "createDiagrams"
	tables["Graphs"] = "createGraphs"
	tables["DataPoints"] = "createDataPoints"

	for table, createQuery := range tables {
		if exists, err := pers.tableExists(table); err != nil {
			return err
		} else if !exists {
			log.Printf("Create table %s ", table)
			_, err := pers.query(createQuery)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (pers sqlDataPersistence) tableExists(name string) (bool, error) {
	//todo be clever bout this
	rows, err := pers.db.Query(pers.q["tableExists"] + " '" + name + "'")

	if err != nil {
		return false, err
	}

	return rows.Next(), nil
}

func (pers sqlDataPersistence) Add(node string, t time.Time, info api.NodeInfo) {

}