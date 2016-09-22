package context

import (
	"github.com/WE-Development/mosel/commons"
	"time"
	"github.com/WE-Development/mosel/api"
	"database/sql"
	"fmt"
	"log"
)

type table struct {
	name        string
	createQuery string
}

type result map[string]map[string]map[string]string

type dbState map[string]map[string][]string

type dbResult struct {
	value     string
	timestamp []uint8
	graph     string
	diagram   string
	node      string
	url       string
}

type dataPersistence interface {
	Init() error
	Add(node string, t time.Time, info api.NodeInfo)
	GetAll() (result, error)
}

type sqlDataPersistence struct {
	db      *sql.DB
	q       commons.SqlQueries

	dbState dbState
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

func (pers sqlDataPersistence) queryResultNotEmpty(name string, args ...interface{}) (bool, error) {
	rows, err := pers.query(name, args...)
	defer rows.Close()

	if err != nil {
		return false, err
	}

	return !rows.Next(), nil
}

func (pers sqlDataPersistence) tableExists(name string) (bool, error) {
	//todo be clever bout this
	rows, err := pers.db.Query(pers.q["tableExists"] + " '" + name + "'")
	defer rows.Close()
	if err != nil {
		return false, err
	}

	return rows.Next(), nil
}

func (pers sqlDataPersistence) Init() error {
	tables := make([]table, 4)
	tables[0] = table{name:"Nodes", createQuery:"createNodes", }
	tables[1] = table{name:"Diagrams", createQuery:"createDiagrams", }
	tables[2] = table{name:"Graphs", createQuery:"createGraphs", }
	tables[3] = table{name:"DataPoints", createQuery:"createDataPoints", }

	for _, table := range tables {
		if exists, err := pers.tableExists(table.name); err != nil {
			return err
		} else if !exists {
			log.Printf("Create table %s ", table)
			_, err := pers.query(table.createQuery)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (pers sqlDataPersistence) Add(node string, t time.Time, info api.NodeInfo) {



}

func (pers sqlDataPersistence) GetAll() (result, error) {
	res := make(result)
	rows, err := pers.query("all")
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var dbRes dbResult
		err := rows.Scan(&dbRes.value, &dbRes.timestamp, &dbRes.graph, &dbRes.diagram, &dbRes.node, &dbRes.url)

		if err != nil {
			return nil, err
		}
		pers.updateDbState(dbRes)
	}
	log.Println(pers.dbState)

	return res, nil
}

func (pers *sqlDataPersistence) updateDbState(dbRes dbResult) {
	//log.Println(value, timestamp, graph, diagram, node, url)
	if dbRes.node == "" {
		return
	}

	if pers.dbState == nil {
		pers.dbState = make(dbState)
	}

	diagrams, ok := pers.dbState[dbRes.node]
	if !ok {
		diagrams = make(map[string][]string)
		pers.dbState[dbRes.node] = diagrams
	}

	if dbRes.diagram == "" {
		return
	}

	graphs, ok := pers.dbState[dbRes.node][dbRes.diagram]
	if !ok {
		graphs = make([]string, 0)
		pers.dbState[dbRes.node][dbRes.diagram] = graphs
	}

	if dbRes.graph == "" {
		return
	}

	graphs = append(graphs, dbRes.graph)
}