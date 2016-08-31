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
	"github.com/WE-Development/mosel/moselserver"
	"github.com/WE-Development/mosel/moseld/server/handler"
	"github.com/WE-Development/mosel/moseld/server/context"
	"github.com/WE-Development/mosel/config"
	"errors"
	"net/url"
	"log"
)

type moseldServer struct {
	moselserver.MoselServer

	config  moselconfig.MoseldServerConfig
	context *context.MoseldServerContext
}

func NewMoseldServer(config moselconfig.MoseldServerConfig) *moseldServer {
	server := moseldServer{
		config: config,
		context: new(context.MoseldServerContext),
	}

	server.MoselServer = moselserver.MoselServer{
		Config: config.MoselServerConfig,
	}

	server.InitFuncs = append(server.InitFuncs,
		server.initDebs,
		server.initNodeCache,
		server.initDataCache)

	server.Handlers = []moselserver.MoselHandler{
		handler.NewLoginHandler(server.context),
		handler.NewPingHandler(),
		handler.NewDebugHandler(server.context),
		handler.NewNodeInfoHandler(server.context),
		handler.NewInfoHandler(server.context),
	}

	return &server
}

/*
 * Initialize Context
 */

func (server *moseldServer) initDebs() error {
	ctx := server.context

	var err error
	if ctx.DataCache, err =
		context.NewDataCache(); err != nil {
		return err
	}
	if ctx.NodeHandler, err =
		context.NewNodeRespHandler(ctx.DataCache); err != nil {
		return err
	}

	if ctx.Scripts, err =
		context.NewScriptCache(server.config.Scripts); err != nil {
		return err
	}

	if ctx.ScriptsRunner, err =
		context.NewScriptsRunner(ctx.Scripts, ctx.DataCache); err != nil {
		return err
	}

	if ctx.Nodes, err =
		context.NewNodeCache(ctx.NodeHandler, ctx.Scripts); err != nil {
		return err
	}

	return nil
}

func (server *moseldServer) initNodeCache() error {
	for nodeName, nodeConf := range server.config.Node {

		var scripts []string

		//get configured
		if len(nodeConf.Scripts) > 0 {
			//manual config
			scripts = nodeConf.Scripts
		} else {
			//default
			scripts = server.context.Scripts.GetScripts()
		}

		//exclude certain scripts
		for _, exclude := range nodeConf.ScriptsExclude {
			for i, script := range scripts {
				if script == exclude {
					scripts = append(scripts[:i], scripts[i + 1:]...)
				}
			}
		}

		nodeScripts := make([]string, 0)
		localScripts := make([]string, 0)

		//sort scripts by scope
		for _, script := range scripts {
			if scriptConf, err := server.context.Scripts.GetScriptConfig(script); err == nil {
				switch scriptConf.Scope {
				case "node":
					nodeScripts = append(nodeScripts, script)
					break
				case "local":
					localScripts = append(localScripts, script)
					break
				default:
					return errors.New("Unknown script scope " + scriptConf.Scope + " for node " + nodeName)
				}

			} else {
				return err
			}
		}

		log.Printf("Scripts configured for node %s: node=%s local=%s", nodeName, nodeScripts, localScripts)

		//get base url
		var url *url.URL
		var err error
		if url, err = url.Parse(nodeConf.URL); err != nil {
			return err
		}

		//instantiate node
		node, err := context.NewNode(
			nodeName,
			*url,
			nodeConf.User,
			nodeConf.Password,
			nodeScripts,
			server.context.NodeHandler,
			server.context.Scripts)

		if err != nil {
			return err
		}

		server.context.Nodes.Add(node)
		server.context.ScriptsRunner.Run(localScripts, node)
	}

	return nil
}

func (server *moseldServer) initDataCache() error {
	//c := server.context.Cache
	return nil
}