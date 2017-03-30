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
	"errors"
)

type nodeCache struct {
	handler *nodeRespHandler
	scripts *scriptCache

	nodes map[string]*node
}

func NewNodeCache(handler *nodeRespHandler, scripts *scriptCache) (*nodeCache, error) {
	c := &nodeCache{}
	c.handler = handler
	c.scripts = scripts
	c.nodes = make(map[string]*node)
	return c, nil
}

func (cache *nodeCache) Add(node *node) {
	cache.nodes[node.Name] = node

	//cache.ProvisionScripts(node.Name, cache.scripts.Scripts)
	go func() {
		node.ListenStream()
	}()
}

func (cache *nodeCache) Get(name string) (*node, error) {
	if val, ok := cache.nodes[name]; ok {
		return val, nil
	}

	return nil, errors.New("No node with name " + name + " registered")
}

func (cache *nodeCache) CloseNode(name string) error {
	node, err := cache.Get(name)

	if err != nil {
		return err
	}

	node.close <- struct{}{}
	return nil
}
