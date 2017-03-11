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
	"errors"
	"net/url"
	"log"
	"fmt"
	"time"

	"github.com/bluedevel/mosel/commons"
	"github.com/bluedevel/mosel/moselserver"
	"github.com/bluedevel/mosel/moseld/server/handler"
	"github.com/bluedevel/mosel/moseld/server/context"
	"github.com/bluedevel/mosel/config"
	"reflect"
)

// The server started by moseld.
// It collects data from the configured nodes and from locally running scripts and provides it over a rest service.
type moseldServer struct {
	moselserver.MoselServer

	config  moselconfig.MoseldServerConfig
	context *context.MoseldServerContext
}

// Construct a new Instance of a MoselServer for a given configuration.
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
		server.initDataPersistence,
		server.initDataCache,
		server.initNodeCache,
	)

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

// Initializes dependencies within the server context.
func (server *moseldServer) initDebs() error {
	ctx := server.context
	ctx.MoselServerContext = server.Context

	var err error

	if ctx.DataCache, err =
		context.NewDataCache(); err != nil {
		return err
	}

	persistenceConfig := server.config.PersistenceConfig
	if persistenceConfig.Enabled {
		dataSource, ok := ctx.DataSources[persistenceConfig.DataSource]

		if !ok {
			return fmt.Errorf("Datasource %s not configured", persistenceConfig.DataSource)
		}

		if dataSource.GetType() == "mysql" {
			var sqlDs moselserver.SqlDataSource
			var correctDsType bool
			if sqlDs, correctDsType = dataSource.(moselserver.SqlDataSource); !correctDsType {
				return fmt.Errorf("Expected datasource of type SqlDataSource but got %s",
					reflect.TypeOf(dataSource).Name())
			}

			var queries commons.SqlQueries
			var err error
			if queries, err = commons.GetQueries(sqlDs.GetType()); err != nil {
				return fmt.Errorf("No queries configured for sql dialect %s", sqlDs.GetType())
			}

			ctx.DataPersistence = context.NewSqlDataPersistence(sqlDs.GetDb(), queries)
		} else if dataSource.GetType() == "mongo" {

		} else {
			return fmt.Errorf("Unable to build persistence for datasource type %s", dataSource.GetType())
		}
	}

	if ctx.NodeHandler, err =
		context.NewNodeRespHandler(ctx.DataCache, ctx.DataPersistence); err != nil {
		return err
	}

	if ctx.Scripts, err =
		context.NewScriptCache(server.config.Scripts); err != nil {
		return err
	}

	if ctx.ScriptsRunner, err =
		context.NewScriptsRunner(ctx.Scripts, ctx.DataCache, ctx.DataPersistence); err != nil {
		return err
	}

	if ctx.Nodes, err =
		context.NewNodeCache(ctx.NodeHandler, ctx.Scripts); err != nil {
		return err
	}

	return nil
}

// Initialize the node cache withe the configured nodes
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

		scripts = commons.ExcludeStr(scripts, nodeConf.ScriptsExclude)
		nodeScripts, localScripts, err := server.sortScriptsByScope(scripts)

		if err != nil {
			return err
		}

		log.Printf("Scripts configured for node %s: node=%s local=%s", nodeName, nodeScripts, localScripts)

		//get base url
		var baseUrl *url.URL
		if baseUrl, err = baseUrl.Parse(nodeConf.URL); err != nil {
			return err
		}

		//instantiate node
		node, err := context.NewNode(
			nodeName,
			*baseUrl,
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

func (server *moseldServer) sortScriptsByScope(scripts []string) ([]string, []string, error) {
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
				return nil, nil, errors.New("Unknown script scope " + scriptConf.Scope)
			}

		} else {
			return nil, nil, err
		}
	}

	return nodeScripts, localScripts, nil
}

// Initialize the data persistence
func (server *moseldServer) initDataPersistence() error {
	if server.context.DataPersistence == nil {
		return nil
	}
	log.Println("Init data persistence")
	defer log.Println("Finished initializing data persistence")
	return server.context.DataPersistence.Init()
}

// Initialize the data persistence
func (server *moseldServer) initDataCache() error {
	var storage context.DataCacheStorage
	var err error

	cacheSize := server.config.DataCache.CacheSize
	log.Printf("Init data cache. cache size: %s", cacheSize)

	if cacheSize == "" {
		storage, err = server.context.DataPersistence.GetAll()
		if err != nil {
			return err
		}
	} else {
		dur, err := time.ParseDuration(cacheSize)
		if err != nil {
			return err
		}

		server.context.DataCache.CacheSize = dur
		storage, err = server.context.DataPersistence.GetAllSince(dur)
		if err != nil {
			return err
		}
	}

	// count points
	pointCount := 0
	for _, points := range storage {
		pointCount += len(points)
	}

	log.Printf("Loaded %d data points. Putting them into cache", pointCount)
	server.context.DataCache.SetStorage(storage)

	// clean up the cache
	go func() {
		for {
			server.context.DataCache.Clean()
			time.Sleep(10 * time.Second)
		}
	}()

	log.Printf("Finished initializing data cache with %d data points", pointCount)
	return nil
}