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
	"github.com/bluedevel/mosel/commons"
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
	Id   bson.ObjectId `bson:"_id"`
	Name string `bson:"name"`
	//DataPoints []dataPointDoc `bson:"dataPoints"`
}

type dataPointDoc struct {
	Time  time.Time `bson:"time"`
	Value string `bson:"value"`
	Graph bson.ObjectId `bson:"graph"`
}

type mongoDataPersistence struct {
	session  *mgo.Session
	database *mgo.Database
}

func NewMongoDataPersistence(session *mgo.Session) *mongoDataPersistence {
	return &mongoDataPersistence{
		session: session,
	}
}

func (pers *mongoDataPersistence) Init() error {
	pers.database = pers.session.DB("")
	return nil
}

func (pers *mongoDataPersistence) getCollections() (*mgo.Collection, *mgo.Collection) {
	return pers.database.C("nodes"),
		pers.database.C("datapoints")
}

func (pers *mongoDataPersistence) Add(nodeName string, t time.Time, info api.NodeInfo) {
	collNodes, collData := pers.getCollections()

	selector := bson.M{"name": nodeName}

	var node nodeDoc
	itr := collNodes.Find(selector).Iter()

	doUpdateNodes := itr.Next(&node)
	commons.LogFatal(itr.Err())

	if node.Name == "" {
		node.Name = nodeName
	}

	if node.Diagrams == nil {
		node.Diagrams = make([]diagramDoc, 0)
	}

	points := make([]interface{}, 0)

	for diagramName, graphs := range info {
		diaIndex := findDiagramByName(diagramName, node.Diagrams)

		var dia diagramDoc
		if diaIndex == -1 {
			dia = diagramDoc{
				Name:   diagramName,
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
					Id:   bson.NewObjectId(),
					Name: graphName,
				}
			} else {
				graph = dia.Graphs[graphIndex]
			}

			point := dataPointDoc{
				Time:  t,
				Value: value,
				Graph: graph.Id,
			}
			points = append(points, point)

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

	var err error
	// persist meta data
	if doUpdateNodes {
		err = collNodes.Update(selector, node)
	} else {
		err = collNodes.Insert(node)
	}
	commons.LogFatal(err)

	// persist points
	err = collData.Insert(points...)
	commons.LogFatal(err)
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

func (pers *mongoDataPersistence) GetAll() (DataCacheStorage, error) {
	return pers.get(bson.M{}), nil
}

func (pers *mongoDataPersistence) GetAllSince(since time.Duration) (DataCacheStorage, error) {
	return pers.get(bson.M{}), nil
}

func (pers *mongoDataPersistence) get(query interface{}) DataCacheStorage {
	res := make(DataCacheStorage)

	graphs := pers.getGraphCache()

	_, collData := pers.getCollections()
	points := collData.Find(bson.M{}).Iter()

	var point dataPointDoc
	for points.Next(&point) {
		info := graphs[point.Graph]

		node, ok := res[info.Node.Name]
		if !ok {
			node = make(map[time.Time]DataPoint)
			res[info.Node.Name] = node
		}

		t := point.Time
		p, ok := node[t]
		if !ok {
			p = DataPoint{
				Time: t,
				Info: make(api.NodeInfo),
			}
			node[t] = p
		}

		diagram, ok := p.Info[info.Diagram.Name]
		if !ok {
			diagram = make(map[string]string)
			p.Info[info.Diagram.Name] = diagram
		}

		diagram[info.Graph.Name] = point.Value
	}

	return res
}

type graphInfo struct {
	Node    nodeDoc
	Diagram diagramDoc
	Graph   graphDoc
}

func (pers *mongoDataPersistence) getGraphCache() map[bson.ObjectId]graphInfo {
	collNodes, _ := pers.getCollections()

	graphs := make(map[bson.ObjectId]graphInfo)

	nodeItr := collNodes.Find(bson.M{}).Iter()
	var node nodeDoc
	for nodeItr.Next(&node) {
		for _, dia := range node.Diagrams {
			for _, gr := range dia.Graphs {
				graphs[gr.Id] = graphInfo{
					Node:    node,
					Diagram: dia,
					Graph:   gr,
				}
			}
		}
	}

	return graphs
}

func (pers *mongoDataPersistence) getNodeDataCaches() (map[string]nodeDoc, map[bson.ObjectId]graphDoc) {
	collNodes, _ := pers.getCollections()

	// build metadata cashes
	nodes := make(map[string]nodeDoc)
	graphs := make(map[bson.ObjectId]graphDoc)

	nodeItr := collNodes.Find(bson.M{}).Iter()
	var node nodeDoc
	for nodeItr.Next(&node) {
		nodes[node.Name] = node

		for _, dia := range node.Diagrams {
			for _, gr := range dia.Graphs {
				graphs[gr.Id] = gr
			}
		}
	}

	return nodes, graphs
}
