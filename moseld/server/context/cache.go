/*
 * Copyright 2016 Robin Engel
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
	"log"
)

type dataCache struct {
	points map[string]map[string][]dataPoint
}

type dataPoint struct {
	time time.Time
	val  float64
}

func NewDataCache() *dataCache {
	c := &dataCache{}
	c.points = make(map[string]map[string][]dataPoint)
	return c
}

func (cache dataCache) Add(node string, name string, t time.Time, val float64) {

	if _, ok := cache.points[node]; !ok {
		log.Println("Alloc")
		cache.points[node] = make(map[string][]dataPoint)
	}

	if _, ok := cache.points[node][name]; !ok {
		log.Println("Alloc2")
		cache.points[node][name] = make([]dataPoint, 0)
	}

	cache.points[node][name] =
	append(cache.points[node][name], dataPoint{
		time: t.Round(time.Second),
		val: val,
	})

	log.Println(cache.points[node][name])
}