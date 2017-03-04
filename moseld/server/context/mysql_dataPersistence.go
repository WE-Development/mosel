/*
 * Copyright 2017 Robin Engel
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package context

import (
	"time"
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/bluedevel/mosel/api"
	"github.com/bluedevel/mosel/commons"
)

// implements DataPersistence
type sqlDataPersistence struct {
	db            *sql.DB
	q             commons.SqlQueries
	serverContext *MoseldServerContext

	dbLock        sync.RWMutex
	dbState       *dbState
}

type table struct {
	name        string
	createQuery string
}

type dbState struct {
	nodes map[string]*dbNodeState
}

type dbNodeState struct {
	id       int
	name     string
	diagrams map[string]*dbDiagramState
}

type dbDiagramState struct {
	id     int
	name   string
	graphs map[string]*dbGraphState
}

type dbGraphState struct {
	id   int
	name string
}

type dbResult struct {
	value     string
	timestamp int64
	graphId   int
	graph     string
	diagramId int
	diagram   string
	nodeId    int
	node      string
	url       string
}

func NewSqlDataPersistence(db *sql.DB, queries commons.SqlQueries) *sqlDataPersistence {
	return &sqlDataPersistence{
		db:db,
		q:queries,
	}
}

func (pers *sqlDataPersistence) query(name string, args ...interface{}) (*sql.Rows, error) {
	query, exists := pers.q[name]

	if !exists {
		return nil, fmt.Errorf("Quers %s is not registered", name)
	}

	rows, err := pers.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	if len(cols) == 0 {
		// nothing to scan here
		rows.Close()
	}

	return rows, err
}

func (pers *sqlDataPersistence) queryResultNotEmpty(name string, args ...interface{}) (bool, error) {
	rows, err := pers.query(name, args...)
	defer rows.Close()

	if err != nil {
		return false, err
	}

	return !rows.Next(), nil
}

func (pers *sqlDataPersistence) tableExists(name string) (bool, error) {
	//todo be clever bout this
	rows, err := pers.db.Query(pers.q["tableExists"] + " '" + name + "'")
	defer rows.Close()
	if err != nil {
		return false, err
	}

	return rows.Next(), nil
}

func (pers *sqlDataPersistence) Init() error {
	tables := make([]table, 4)
	tables[0] = table{name:"Nodes", createQuery:"createNodes", }
	tables[1] = table{name:"Diagrams", createQuery:"createDiagrams", }
	tables[2] = table{name:"Graphs", createQuery:"createGraphs", }
	tables[3] = table{name:"DataPoints", createQuery:"createDataPoints", }

	for _, table := range tables {
		if exists, err := pers.tableExists(table.name); err != nil {
			return err
		} else if !exists {
			log.Printf("Create table %s ", table.name)
			_, err := pers.query(table.createQuery)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (pers *sqlDataPersistence) Add(node string, t time.Time, info api.NodeInfo) {
	pers.dbLock.Lock()
	defer pers.dbLock.Unlock()

	if pers.dbState == nil {
		log.Println("No database state initialized")
		return
	}

	nodeState, ok := pers.dbState.nodes[node]

	if !ok {
		if _, err := pers.query("insertNode", node, ""); err != nil {
			log.Println(err)
			return
		}

		// update the state. dirty but what the heck
		pers.getAll()
		nodeState = pers.dbState.nodes[node]
	}

	for diagram, graphs := range info {
		diagramState, ok := nodeState.diagrams[diagram]

		if !ok {
			if _, err := pers.query("insertDiagram", diagram, nodeState.id); err != nil {
				log.Println(err)
				return
			}

			// update the state. dirty but what the heck
			pers.getAll()
			diagramState = pers.dbState.nodes[node].diagrams[diagram]
		}

		for graph, value := range graphs {
			graphState, ok := diagramState.graphs[graph]

			if !ok {
				if _, err := pers.query("insertGraph", graph, diagramState.id); err != nil {
					log.Println(err)
					return
				}

				// update the state. dirty but what the heck
				pers.getAll()
				graphState = pers.dbState.nodes[node].diagrams[diagram].graphs[graph]
			}

			if _, err := pers.query("insertDataPoint", value, t.Round(time.Second).Unix(), graphState.id); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func (pers *sqlDataPersistence) GetAll() (DataCacheStorage, error) {
	pers.dbLock.RLock()
	defer pers.dbLock.RUnlock()
	return pers.getAll()
}

func (pers *sqlDataPersistence) getAll() (DataCacheStorage, error) {
	res := make(DataCacheStorage)

	rows, err := pers.query("all")

	if err != nil {
		return res, err
	}

	defer rows.Close()
	return pers.get(rows)
}

func (pers *sqlDataPersistence) GetAllSince(since time.Duration) (DataCacheStorage, error) {
	res := make(DataCacheStorage)

	pers.dbLock.RLock()
	defer pers.dbLock.RUnlock()

	t := time.Now().Add(-since)
	rows, err := pers.query("allSince", t)

	if err != nil {
		return res, err
	}

	defer rows.Close()
	return pers.get(rows)
}

func (pers *sqlDataPersistence) get(rows *sql.Rows) (DataCacheStorage, error) {
	res := make(DataCacheStorage)

	// dirty force reset
	pers.dbState = &dbState{
		nodes:make(map[string]*dbNodeState),
	}

	for rows.Next() {
		var dbRes dbResult
		err := rows.Scan(
			&dbRes.value,
			&dbRes.timestamp,
			&dbRes.graphId,
			&dbRes.graph,
			&dbRes.diagramId,
			&dbRes.diagram,
			&dbRes.nodeId,
			&dbRes.node,
			&dbRes.url)

		if err != nil {
			return nil, err
		}
		pers.updateDbState(dbRes)

		//parse db result
		node, ok := res[dbRes.node]
		if !ok {
			node = make(map[time.Time]DataPoint)
			res[dbRes.node] = node
		}

		t := time.Unix(dbRes.timestamp, 0)
		point, ok := node[t]
		if !ok {
			point = DataPoint{
				Time:t,
				Info:make(api.NodeInfo),
			}
			node[t] = point
		}

		diagram, ok := point.Info[dbRes.diagram]
		if !ok {
			diagram = make(map[string]string)
			point.Info[dbRes.diagram] = diagram
		}

		diagram[dbRes.graph] = dbRes.value
	}

	return res, nil
}

func (pers *sqlDataPersistence) updateDbState(dbRes dbResult) {
	if dbRes.nodeId == -1 {
		return
	}

	node, ok := pers.dbState.nodes[dbRes.node]
	if !ok {
		node = &dbNodeState{
			id:dbRes.nodeId,
			name:dbRes.node,
			diagrams:make(map[string]*dbDiagramState),
		}
		pers.dbState.nodes[dbRes.node] = node
	}

	if dbRes.diagramId == -1 {
		return
	}

	diagram, ok := node.diagrams[dbRes.diagram]
	if !ok {
		diagram = &dbDiagramState{
			id:dbRes.diagramId,
			name:dbRes.diagram,
			graphs:make(map[string]*dbGraphState),
		}
		node.diagrams[dbRes.diagram] = diagram
	}

	if dbRes.graphId == -1 {
		return
	}

	graph, ok := diagram.graphs[dbRes.graph]
	if !ok {
		graph = &dbGraphState{
			id:dbRes.graphId,
			name:dbRes.graph,
		}
		diagram.graphs[dbRes.graph] = graph
	}
}
