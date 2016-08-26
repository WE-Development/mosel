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
	"io/ioutil"
	"strings"
	"log"
)

type scriptCache struct {
	path    string
	Scripts []string
}

func NewScriptCache(path string) (*scriptCache, error) {
	c := &scriptCache{}
	c.path = path
	return c, c.initialize()
}

func (cache *scriptCache) initialize() error {

	files, err := ioutil.ReadDir(cache.path)

	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if strings.HasSuffix(
			file.Name(), ".sh") {
			cache.Scripts = append(cache.Scripts, file.Name())
			log.Printf("Registerd script %s", file.Name())
		}
	}

	return nil
}

func (cache *scriptCache) getScript(name string) (string,error) {
	bytes, err := cache.getScriptBytes(name)
	return string(bytes), err
}

func (cache *scriptCache) getScriptBytes(name string) ([]byte,error) {
	return ioutil.ReadFile(cache.path + "/" + name)
}