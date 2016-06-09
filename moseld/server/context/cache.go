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
	"github.com/bluedevel/mosel/api"
	"sync"
)

type dataCache struct {
	points map[string][]dataPoint

	m sync.Mutex
}

type dataPoint struct {
	time time.Time
	info api.NodeInfo
}

func NewDataCache() *dataCache {
	c := &dataCache{}
	c.points = make(map[string][]dataPoint)
	return c
}

func (cache *dataCache) Add(node string, t time.Time, info api.NodeInfo) {

	var arr []dataPoint

	cache.m.Lock()

	if _, ok := cache.points[node]; !ok {
		arr = make([]dataPoint, 0)
	} else {
		arr = cache.points[node]
	}

	arr = append(arr, dataPoint{
		time: t.Round(time.Second),
		info: info,
	})

	cache.points[node] = arr
	log.Println(cache.points[node])

	cache.m.Unlock()
}