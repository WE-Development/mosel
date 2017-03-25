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
	Diagrams []diagramDoc `bson:"diagrams"`
}

type diagramDoc struct {
	Name   string `bson:"name"`
	Graphs []graphDoc `bson:"graphs"`
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

	if node.Name == "" {
		node.Name = nodeName
	}

	if node.Diagrams == nil {
		node.Diagrams = make([]diagramDoc, len(info))
	}

	for diagramName, graphs := range info {
		diaIndex := findDiagramByName(diagramName, node.Diagrams)

		var dia diagramDoc
		if diaIndex == -1 {
			dia = diagramDoc{
				Name:diagramName,
				Graphs: make([]graphDoc, 0),
			}
		} else {
			dia = node.Diagrams[diaIndex]
		}

		for graphName, value := range graphs {
			graphIndex := findGraphByName(graphName, dia.Graphs)

			var graph graphDoc
			if graphIndex == -1 {
				graph = graphDoc{
					Name:graphName,
					DataPoints: make([]dataPointDoc, 0),
				}
			} else {
				graph = dia.Graphs[graphIndex]
			}

			point := dataPointDoc{
				Time:t,
				Value:value,
			}

			graph.DataPoints = append(graph.DataPoints, point)

			if graphIndex == -1 {
				dia.Graphs = append(dia.Graphs, graph)
			} else {
				dia.Graphs[graphIndex] = graph
			}
		}

		if diaIndex == -1 {
			node.Diagrams = append(node.Diagrams, dia)
		} else {
			node.Diagrams[diaIndex] = dia
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

func findDiagramByName(name string, diagrams []diagramDoc) int {
	for i, dia := range diagrams {
		if dia.Name == name {
			return i
		}
	}
	return -1
}

func findGraphByName(name string, graphs []graphDoc) int {
	for i, gr := range graphs {
		if gr.Name == name {
			return i
		}
	}
	return -1
}