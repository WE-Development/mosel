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
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"fmt"
	"strconv"
)

type MoselServer struct {
	Config  MoselServerConfig
	Context MoselServerContext
}

func (server *MoselServer) Run() error {

	//initializing server context
	err := server.initContext()

	if ! server.Context.IsInitialized {

		if err != nil {
			return err
		}

		return fmt.Errorf("Mosel Server - Run: Context wasn't initialized correctly")
	}

	//init router and handlers
	r := mux.NewRouter()
	//server.initHandler(r)
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

func (server *MoselServer) initContext() error {

	initFns := []func() error{
		server.initAuth,
		server.initSessionCache,
	}

	for _, fn := range initFns {
		err := fn()

		if (err != nil) {
			return err
		}
	}

	server.Context.IsInitialized = true
	return nil
}

func (server *MoselServer) initAuth() error {
	config := server.Config

	var enabledCount int = 0

	if config.AuthSys.Enabled {
		enabledCount++
	}

	if config.AuthMySQL.Enabled {
		enabledCount++
	}

	if config.AuthTrue.Enabled {
		enabledCount++
		log.Println("Using AuthTrue! This is for debug purposes only, make sure you don't deploy this in production")
		server.Context.Auth = AuthTrue{}
	}

	if enabledCount > 1 {
		return fmt.Errorf("More then one auth services enabled")
	} else if enabledCount == 0 {
		return fmt.Errorf("No auth service configured")
	}

	return nil
}

func (server *MoselServer) initSessionCache() error {
	c := NewSessionCache()
	server.Context.Sessions = *c
	return nil
}


/*
 * Initialize Handler
 */

func (server *MoselServer) initHandler(r *mux.Router) {

	var handlers = []MoselHandler{
	}

	for n, _ := range handlers {

		h := handlers[n]

		f := func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTPContext(server.Context, w, r)
		}

		secure := h.Secure()

		if secure {
			f = server.secure(f)
		}

		log.Printf("Handling %s - secure=%s", h.GetPath(), strconv.FormatBool(secure))
		r.HandleFunc(h.GetPath(), f)
	}
}