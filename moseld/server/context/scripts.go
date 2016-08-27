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
	"github.com/WE-Development/mosel/config"
)

type scriptCache struct {
	path    string
	Scripts map[string]moselconfig.ScriptConfig
}

func NewScriptCache(configs map[string]*moselconfig.ScriptConfig) (*scriptCache, error) {
	cache := &scriptCache{}
	return cache, cache.initialize(configs)
}

func (cache *scriptCache) initialize(configs map[string]*moselconfig.ScriptConfig) error {
	scripts := make(map[string]moselconfig.ScriptConfig)
	for script, conf := range configs {
		scripts[script] = *conf
	}
	cache.Scripts = scripts
	return nil
}

func (cache *scriptCache) GetScripts() []string {
	scripts := make([]string, 0, len(cache.Scripts))
	for script := range cache.Scripts {
		scripts = append(scripts, script)
	}
	return scripts
}

func (cache *scriptCache) getScript(name string) (string, error) {
	bytes, err := cache.getScriptBytes(name)
	return string(bytes), err
}

func (cache *scriptCache) getScriptBytes(name string) ([]byte, error) {
	return ioutil.ReadFile(cache.path + "/" + name)
}