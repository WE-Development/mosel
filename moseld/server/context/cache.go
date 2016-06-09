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

import "time"

type dataCache struct {
	times  []*time.Time
	points map[string][]dataPoint
}

type dataPoint struct {
	time *time.Time
	val  float64
}

func NewDataCache() *dataCache {
	c := &dataCache{}
	c.points = make(map[string][]dataPoint)
	return c
}

func (cache dataCache) Add(name string, val float64) {
	points, ok := cache.points[name]

	if !ok {
		points = make([]dataPoint, 0)
	}

	now := time.Now()
	points = append(points, dataPoint{
		time: &now,
		val: val,
	})

}