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
package moselnodedserver

import (
	"github.com/WE-Development/mosel/moselserver"
	"github.com/WE-Development/mosel/moselnoded/server/handler"
	"github.com/WE-Development/mosel/moselnoded/server/context"
)

type moselnodedServer struct {
	moselserver.MoselServer

	config MoselNodedServerConfig
	context *context.MoselnodedServerContext
}

func NewMoselNodedServer(config MoselNodedServerConfig) *moselnodedServer {
	server := moselnodedServer{
		config: config,
		context: new(context.MoselnodedServerContext),
	}

	server.MoselServer = moselserver.MoselServer{
		Config: config.MoselServerConfig,
	}

	server.InitFuncs = append(server.InitFuncs,
		server.initCollector,)

	server.Handlers = []moselserver.MoselHandler{
		handler.NewPingHandler(),
		handler.NewStreamHandler(server.context),
		handler.NewScriptHandler(server.context),
	}

	return &server
}

func (server *moselnodedServer) initCollector() error {
	ctx := server.context
	ctx.Collector = context.NewCollector()
	return nil
}