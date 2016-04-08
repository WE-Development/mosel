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
package server

import (
	"reflect"
	"log"
	"fmt"
)

type dataSaver struct {

}

func NewDataSaver(config MoselServerConfig) *dataSaver {
	saver := new(dataSaver)
	return saver
}

func (dataSaver *dataSaver) SaveEntry(i interface{}) error {
	log.Printf("TODO save to %s", getFilename(i))
	return nil
}

func getFilename(i interface{}) (string, error) {
	name := reflect.TypeOf(i).Name()

	if name == "" {
		return "", fmt.Errorf("Could not evaluate filename to save data to")
	}

	return name + ".msl", nil
}