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
	"gopkg.in/mgo.v2"
	"time"
	"github.com/bluedevel/mosel/api"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type nodeDoc struct {
	Name     string `bson:"name"`
	Diagrams map[string]diagramDoc `bson:"diagrams"`
}

type diagramDoc struct {
	Name   string `bson:"name"`
	Graphs map[string]graphDoc `bson:"graphs"`
}

type graphDoc struct {
	Name       string `bson:"name"`
	DataPoints []dataPointDoc `bson:"dataPoints"`
}

type dataPointDoc struct {
	Time  time.Time `bson:"time"`
	Value string `bson:"value"`
}

type mongoDataPersistence struct {
	session  *mgo.Session
	database *mgo.Database
}

func NewMongoDataPersistence(session *mgo.Session) *mongoDataPersistence {
	return &mongoDataPersistence{
		session:session,
	}
}

func (pers *mongoDataPersistence) Init() error {
	pers.database = pers.session.DB("")
	return nil
}

func (pers *mongoDataPersistence) Add(nodeName string, t time.Time, info api.NodeInfo) {
	coll := pers.database.C("nodes")

	selector := bson.M{"name": nodeName}

	var node nodeDoc
	itr := coll.Find(selector).Iter()

	doUpdate := itr.Next(&node)

	if err := itr.Err(); err != nil {
		log.Fatal(err)
		return
	}

	initNodeIfNil(nodeName, &node)

	for diagramName, graphs := range info {
		var diagram diagramDoc
		initDiagramIfNil(diagramName, &node, &diagram)

		for graphName, value := range graphs {
			var graph graphDoc
			initGraphIfNil(graphName, &diagram, &graph)

			point := dataPointDoc{
				Time:t,
				Value:value,
			}

			graph.DataPoints = append(graph.DataPoints, point)
		}
	}

	if doUpdate {
		coll.Update(selector, node)
	} else {
		coll.Insert(node)
	}
}

func (pers *mongoDataPersistence) GetAll() (DataCacheStorage, error) {
	return nil, nil
}

func (pers *mongoDataPersistence) GetAllSince(since time.Duration) (DataCacheStorage, error) {
	return nil, nil
}

func initNodeIfNil(name string, node *nodeDoc) {
	if node.Diagrams == nil {
		node.Diagrams = make(map[string]diagramDoc)
	}
}

func initDiagramIfNil(name string, node *nodeDoc, dia *diagramDoc) {
	if _, ok := node.Diagrams[name]; !ok {
		dia.Name = name
		dia.Graphs = make(map[string]graphDoc)
		node.Diagrams[name] = *dia
	} else {
		(*dia) = node.Diagrams[name]
	}
}

func initGraphIfNil(name string, dia *diagramDoc, gr *graphDoc) {
	if _, ok := dia.Graphs[name]; !ok {
		gr.Name = name
		gr.DataPoints = make([]dataPointDoc, 0)
		dia.Graphs[name] = *gr
	} else {
		(*gr) = dia.Graphs[name]
	}
}