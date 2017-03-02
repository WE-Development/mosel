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
	"sync"
	"errors"

	"github.com/bluedevel/mosel/api"
)

type DataCacheStorage map[string]map[time.Time]DataPoint

// Collects and caches NodeInfo objects by node name and time (unix).
// todo optimize locks
type dataCache struct {
	points    DataCacheStorage
	cacheLock sync.RWMutex
	CacheSize time.Duration
}

type DataPoint struct {
	Time time.Time
	Info api.NodeInfo
}

func NewDataCache() (*dataCache, error) {
	c := &dataCache{}
	c.points = make(DataCacheStorage)
	return c, nil
}

func (cache *dataCache) SetStorage(storage DataCacheStorage) {
	cache.cacheLock.Lock()
	defer cache.cacheLock.Unlock()

	cache.points = storage
}

// Cache a new node info object.
// If a node info is already cached for this node and time, the existing data gets extended or overwritten.
func (cache *dataCache) Add(node string, t time.Time, info api.NodeInfo) {
	var points map[time.Time]DataPoint

	cache.cacheLock.Lock()
	defer cache.cacheLock.Unlock()

	t = t.Round(time.Second)

	if _, ok := cache.points[node]; !ok {
		points = make(map[time.Time]DataPoint)
		cache.points[node] = points
	} else {
		points = cache.points[node]
	}

	if point, ok := points[t]; ok {
		for diag, graphs := range info {
			for graph, value := range graphs {
				if _, ok := point.Info[diag]; !ok {
					point.Info[diag] = make(map[string]string)
				}
				point.Info[diag][graph] = value
			}
		}
	} else {
		points[t] = DataPoint{
			Time: t,
			Info: info,
		}
	}
}

func (cache *dataCache) Clean() {
	cache.cacheLock.Lock()
	defer cache.cacheLock.Unlock()

	maxAge := time.Now().Add(-cache.CacheSize)
	for _, points := range cache.points {
		for stamp := range points {
			if stamp.Unix() < maxAge.Unix() {
				delete(points, stamp)
			}
		}
	}
}

// Get node info at a given time (rounded to seconds)
func (cache *dataCache) Get(node string, t time.Time) (api.NodeInfo, error) {
	cache.cacheLock.RLock()
	defer cache.cacheLock.RUnlock()

	points, err := cache.getAllByTime(node)

	if err != nil {
		return api.NodeInfo{}, err
	}

	t = t.Round(time.Second)
	if point, ok := points[t]; ok {
		return point.Info, nil
	}

	return api.NodeInfo{}, errors.New("No datapoint found for time " + t.String())
}

// Get all node infos for a node and not older than a given time (rounded to seconds)
func (cache *dataCache) GetSince(node string, t time.Time) ([]DataPoint, error) {
	cache.cacheLock.RLock()
	defer cache.cacheLock.RUnlock()

	points, err := cache.getAllByTime(node)

	if err != nil {
		return nil, err
	}

	result := make([]DataPoint, 0)
	t = t.Round(time.Second)
	for pt, p := range points {
		if pt.Unix() > t.Unix() {
			result = append(result, p)
		}
	}

	return result, nil
}

// Get all cached node infos
func (cache *dataCache) GetAll(node string) ([]DataPoint, error) {
	var points map[time.Time]DataPoint
	var err error

	cache.cacheLock.RLock()
	defer cache.cacheLock.RUnlock()

	if points, err = cache.getAllByTime(node); err != nil {
		return nil, err
	}

	res := make([]DataPoint, len(points))
	for _, point := range points {
		res = append(res, point)
	}
	return res, nil
}

// Get nodes for which data is cached
func (cache *dataCache) GetNodes() []string {
	nodes := make([]string, len(cache.points))

	i := 0
	for k := range cache.points {
		nodes[i] = k
		i++
	}

	return nodes
}

// Get all node infos for a given time (rounded to seconds)
func (cache *dataCache) getAllByTime(node string) (map[time.Time]DataPoint, error) {
	points, ok := cache.points[node]
	if !ok {
		return nil, errors.New("No node with name " + node)
	}
	return points, nil
}