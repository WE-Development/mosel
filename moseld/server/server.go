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
package moseldserver

import (
	"github.com/bluedevel/mosel/moselserver"
	"github.com/bluedevel/mosel/moseld/server/handler"
	"github.com/bluedevel/mosel/moseld/server/context"
	"net/url"
)

type moseldServer struct {
	moselserver.MoselServer

	config  MoseldServerConfig
	context context.MoseldServerContext
}

func NewMoseldServer(config MoseldServerConfig) *moseldServer {
	server := moseldServer{
		config: config,
	}

	server.MoselServer = moselserver.MoselServer{
		Config: config.MoselServerConfig,
	}

	server.InitFuncs = append(server.InitFuncs,
		server.initNodeCache,
		server.initDataCache)

	server.Handlers = []moselserver.MoselHandler{
		handler.NewLoginHandler(&server.Context),
		handler.NewPingHandler(),
		handler.NewDebugHandler(&server.context),
	}

	return &server
}

/*
 * Initialize Context
 */

func (server *moseldServer) initNodeCache() error {
	c, err := context.NewNodeCache()
	server.context.Nodes = *c

	url, _ := url.Parse("http://localhost:8181/stream")
	c.Add(&context.Node{
		Name: "self",
		URL: *url,
	})

	return err
}

func (server *moseldServer) initDataCache() error {
	c, err := context.NewDataCache()
	server.context.Cache = *c
	return err
}