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
	"errors"

	"github.com/bluedevel/mosel/config"
)

type scriptCache struct {
	scripts map[string]moselconfig.ScriptConfig
}

func NewScriptCache(configs map[string]*moselconfig.ScriptConfig) (*scriptCache, error) {
	cache := &scriptCache{}
	return cache, cache.initialize(configs)
}

func (cache *scriptCache) initialize(configs map[string]*moselconfig.ScriptConfig) error {
	scripts := make(map[string]moselconfig.ScriptConfig)
	for script, conf := range configs {
		if conf.Scope == "" {
			//set default scope
			conf.Scope = "node"
		}

		scripts[script] = *conf
	}
	cache.scripts = scripts
	return nil
}

func (cache *scriptCache) GetScripts() []string {
	scripts := make([]string, 0, len(cache.scripts))
	for script := range cache.scripts {
		scripts = append(scripts, script)
	}
	return scripts
}

func (cache *scriptCache) GetScriptsByScope(scope string) []string {
	scripts := make([]string, 0)
	for name, script := range cache.scripts {
		if script.Scope == scope {
			scripts = append(scripts, name)
		}
	}
	return scripts
}

func (cache *scriptCache) GetScriptConfig(name string) (moselconfig.ScriptConfig, error) {
	if conf, ok := cache.scripts[name]; ok {
		return conf, nil
	}

	return moselconfig.ScriptConfig{}, errors.New("No script with name " + name + " loaded")
}

func (cache *scriptCache) getScript(name string) (string, error) {
	bytes, err := cache.getScriptBytes(name)
	return string(bytes), err
}

func (cache *scriptCache) getScriptBytes(name string) ([]byte, error) {
	if _, ok := cache.scripts[name]; !ok {
		return nil, errors.New("No script with name " + name)
	}
	return ioutil.ReadFile(cache.scripts[name].Path)
}
