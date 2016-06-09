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
	points map[string][]dataPoint
}

type dataPoint struct {
	time time.Time
	val  float64
}

func NewDataCache() *dataCache {
	c := &dataCache{}
	c.points = make(map[string][]dataPoint)
	return c
}

func (cache dataCache) Add(name string, t time.Time, val float64) {
	points, ok := cache.points[name]

	if !ok {
		log.Println("Alloc")
		points = make([]dataPoint, 0)
	}

	points = append(points, dataPoint{
		time: t.Round(time.Second),
		val: val,
	})

	log.Println(points)
	cache.points[name] = points
}