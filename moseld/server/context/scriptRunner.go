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
	"strings"

	"github.com/bluedevel/mosel/commons"
	"github.com/bluedevel/mosel/api"
)

type scriptsRunner struct {
	scriptCache     *scriptCache
	dataCache       *dataCache
	dataPersistence *sqlDataPersistence
	closer          map[string]chan struct{}
}

func NewScriptsRunner(scriptCache *scriptCache, dataCache *dataCache, dataPersistence *sqlDataPersistence) (*scriptsRunner, error) {
	runner := &scriptsRunner{}
	runner.scriptCache = scriptCache
	runner.dataCache = dataCache
	runner.dataPersistence = dataPersistence
	runner.closer = make(map[string]chan struct{})
	return runner, runner.initialize()
}

func (runner *scriptsRunner) initialize() error {
	return nil
}

func (runner *scriptsRunner) Run(scripts []string, node *node) {
	close := make(chan struct{})
	runner.closer[node.Name] = close
	go func() {
		Run: for {
			select {
			case <-close:
				break Run
			default:
				runner.runScripts(scripts, node)
			//todo make reconnection timeout configurable by moseld.conf
				time.Sleep(1 * time.Second)
			}
		}
	}()
}

func (runner *scriptsRunner) runScripts(scripts []string, node *node) {
	info := api.NodeInfo{}
	for _, script := range scripts {
		conf, err := runner.scriptCache.GetScriptConfig(script)
		if err != nil {
			runner.logError(script, nil, node, err)
			continue
		}

		args := make([]string, len(conf.Arguments) + 1)
		args[0] = conf.Path

		//interpret arguments
		for i, argId := range conf.Arguments {
			var arg string
			switch argId {
			case "node":
				arg = node.Name
				break
			case "host":
				arg = strings.Split(node.URL.Host, ":")[0]
				break
			case "port":
				s := strings.Split(node.URL.Host, ":")
				if len(s) == 2 {
					arg = s[1]
				}
				break
			case "scheme":
				arg = node.URL.Scheme
				break
			case "path":
				arg = node.URL.Path
				break
			case "query":
				arg = node.URL.RawQuery
			default:
				arg = "nil"
				break
			}
			args[i + 1] = arg
		}

		res, err := commons.ExecuteScript(args...)
		if err != nil {
			runner.logError(script, args, node, err)
		}
		info[script] = res
	}
	runner.dataCache.Add(node.Name, time.Now(), info)
	if runner.dataPersistence != nil {
		runner.dataPersistence.Add(node.Name, time.Now(), info)
	}
}

func (runner *scriptsRunner) logError(script string, args []string, node *node, err error) {
	log.Printf("Error while executing local script %s on node %s: %s. Args: %s", script, node.Name, err, args)
}