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
	"github.com/bluedevel/mosel/api"
	"os"
	"log"
	"io/ioutil"
	"github.com/bluedevel/mosel/commons"
)

type collector struct {
	scriptFolder string
	scripts      []string
}

func NewCollector() *collector {
	return &collector{
		scriptFolder: "/tmp/mosel",
		scripts: make([]string, 0),
	}
}

func (collector *collector) AddScript(name string, src []byte) error {
	filePath := collector.scriptFolder + "/" + name

	if _, err := mkdirIfNotExist(collector.scriptFolder, 0764); err != nil {
		return err
	}

	err := ioutil.WriteFile(filePath, src, 0775)

	if err != nil {
		log.Println(err)
		return err
	}

	collector.scripts = append(collector.scripts, name)
	log.Printf("Added script %s", name)
	return nil
}

func (collector *collector) FillNodeInfo(info *api.NodeInfo) {
	for _, script := range collector.scripts {
		res, err := commons.ExecuteScript(collector.scriptFolder + "/" + script)

		if err != nil {
			log.Printf("Error executing script %s: %s", script, err.Error())
			continue
		}

		(*info)[script] = res
	}
}

func mkdirIfNotExist(path string, perm os.FileMode) (bool, error) {
	if ok, _ := exists(path); !ok {
		err := os.Mkdir(path, perm)
		return err != nil, err
	}
	return false, nil
}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}