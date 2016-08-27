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
package moselconfig

type Optional struct {
	Enabled bool `gcfg:"enabled"`
}

type MoselServerConfig struct {
	Http      struct {
			  BindAddress string `gcfg:"bind"`
		  } `gcfg:"http"`

	//Auth stuff
	Sessions  struct {
			  Optional
		  } `gcfg:"sessions"`

	AuthSys   struct {
			  Optional
		  } `gcfg:"auth-sys"`
	AuthMySQL struct {
			  Optional
		  } `gcfg:"auth-mysql"`
	AuthTrue  struct {
			  Optional
		  } `gcfg:"auth-true"`
}

type MoseldServerConfig struct {
	MoselServerConfig

	Script struct {
		       Path string `gcfg:"path"`
	       } `gcfg:"script"`

	Node   map[string]*NodeConfig `gcfg:"node"`
}

type NodeConfig struct {
	URL            string `gcfg:"url"`
	Scripts        []string `gcfg:"script"`
	ScriptsExclude []string `gcfg:"exclude-script"`
}

type MoselNodedServerConfig struct {
	MoselServerConfig
}