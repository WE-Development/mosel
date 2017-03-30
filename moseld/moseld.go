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
package main

import (
	"log"
	"os"

	"gopkg.in/gcfg.v1"

	"github.com/bluedevel/mosel/moseld/server"
	"github.com/bluedevel/mosel/config"
)

func main() {

	config, err := loadConfig()

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	server := moseldserver.NewMoseldServer(*config)
	log.Fatal(server.Run())
}

func loadConfig() (*moselconfig.MoseldServerConfig, error) {
	config := new(moselconfig.MoseldServerConfig)
	err := gcfg.ReadFileInto(config, "/etc/mosel/moseld.conf")
	return config, err
}
