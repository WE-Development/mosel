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
package moselserver

import (
	"net/http"
	"log"
	"fmt"
	"strconv"
	"github.com/gorilla/mux"

	"github.com/bluedevel/mosel/config"
)

// The abstract http-server type underlying the mosel servers.
type MoselServer struct {
	Config    moselconfig.MoselServerConfig
	Context   *MoselServerContext

	Handlers  []MoselHandler
	InitFuncs []func() error
}

// Boot up the server
// 1: Run init functions
// 2: Init request handlers and wrap them with gorilla/mux
// 3: Handle the gorilla/mux router
func (server *MoselServer) Run() error {

	//initializing server context
	server.InitFuncs = append([]func() error{
		server.initAuth,
		server.initDataSources,
		server.initSessionCache,
	}, server.InitFuncs...)

	err := server.initContext()

	if ! server.Context.IsInitialized {

		if err != nil {
			return err
		}

		return fmt.Errorf("Mosel Server - Run: Context wasn't initialized correctly")
	}

	//init router and handlers
	r := mux.NewRouter()
	server.initHandler(r)
	http.Handle("/", r)

	addr := server.Config.Http.BindAddress
	log.Printf("Binding http server to %s", addr)

	//do async jobs after initialization here
	errors := make(chan error)

	go func() {
		errors <- http.ListenAndServe(addr, nil)
	}()

	return <-errors
}

/*
 * Initialize Context
 */

// Initialize the server context.
// The configured init functions will be called and on success server.Context.IsInitialized will be set to true.
func (server *MoselServer) initContext() error {
	server.Context = &MoselServerContext{}

	for _, fn := range server.InitFuncs {
		err := fn()

		if (err != nil) {
			return err
		}
	}

	server.Context.IsInitialized = true
	return nil
}

// Initialize the configured authentication method
func (server *MoselServer) initAuth() error {
	config := server.Config

	var enabledCount int = 0

	if config.AuthStatic.Enabled {
		enabledCount++
		server.Context.Auth = &AuthStatic{
			Users: config.Users,
		}
	}

	if config.AuthSys.Enabled {
		enabledCount++
		server.Context.Auth = &AuthSys{
			AllowedUsers: config.AuthSys.AllowedUsers,
		}
	}

	if config.AuthMySQL.Enabled {
		enabledCount++
	}

	if config.AuthTrue.Enabled {
		enabledCount++
		log.Println("Using AuthTrue! This is for debug purposes only, make sure you don't deploy this in production")
		server.Context.Auth = &AuthTrue{}
	}

	if enabledCount > 1 {
		return fmt.Errorf("More then one auth services enabled")
	} else if enabledCount == 0 {
		return fmt.Errorf("No auth service configured")
	}

	return nil
}

// Initialize the configured data sources
func (server *MoselServer) initDataSources() error {
	server.Context.DataSources = make(map[string]dataSource)

	for name, config := range server.Config.DataSources {
		var ds dataSource
		var err error

		if config.Type == "mysql" {
			ds, err = server.initMySql(config.Type, config.Connection)
		} else if config.Type == "mongo" {
			ds, err = server.initMongo(config.Type, config.Connection)
		} else {
			return fmt.Errorf("Data source type '%s' not supported", config.Type)
		}

		if err != nil {
			return err
		}

		log.Printf("Register data source %s of type %s", name, ds.GetType())
		server.Context.DataSources[name] = ds
	}

	return nil
}

// Initialize the session cache
// This gets executed even if sessions are disabled by the config as the decision on weather to use session is
// taken in the authInterceptor
func (server *MoselServer) initSessionCache() error {
	c := NewSessionCache()
	server.Context.Sessions = *c
	return nil
}


/*
 * Initialize Handler
 */

// Initialize the configured http handlers and wrap them into a gorilla/mux router
func (server *MoselServer) initHandler(r *mux.Router) {

	for n, _ := range server.Handlers {

		h := server.Handlers[n]

		f := func(w http.ResponseWriter, r *http.Request) {
			//h.ServeHTTPContext(server.Context, w, r)
			h.ServeHTTP(w, r)
		}

		secure := h.Secure()

		if secure {
			f = server.secure(f)
		}

		log.Printf("Handling %s - secure=%s", h.GetPath(), strconv.FormatBool(secure))
		r.HandleFunc(h.GetPath(), f)
	}
}