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
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"fmt"
	"strconv"
)

type moselServer struct {
	config  MoselServerConfig
	context MoselServerContext
}

func NewMoselServer(config MoselServerConfig) *moselServer {
	server := new(moselServer)
	server.config = config

	return server
}

func (server *moselServer) Run() error {

	err := server.initContext()

	if ! server.context.isInitialized {

		if err != nil {
			return err
		}

		return fmt.Errorf("Mosel Server - Run: Context wasn't initialized correctly")
	}

	r := mux.NewRouter()
	server.initHandler(r)
	http.Handle("/", r)

	addr := server.config.Http.BindAddress
	log.Printf("Binding http server to %s", addr)
	return http.ListenAndServe(addr, nil)
}

func (server *moselServer) initContext() error {

	initFns := []func() error{
		server.initAuth,
	}

	for _, fn := range initFns {
		err := fn()

		if (err != nil) {
			return err
		}
	}

	server.context.isInitialized = true
	return nil
}

func (server *moselServer) initAuth() error {
	config := server.config

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
		server.context.auth = authTrue{}
	}

	if enabledCount > 1 {
		return fmt.Errorf("More then one auth services enabled")
	} else if enabledCount == 0 {
		return fmt.Errorf("No auth service configured")
	}

	return nil
}

func (server *moselServer) initHandler(r *mux.Router) {

	var handlers = []MoselHandler{pingHandler{}, loginHandler{},
	}

	for _, h := range handlers {

		f := func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTPContext(server.context, w, r)
		}

		secure := h.Secure()

		if secure {
			f = server.secure(f)
		}

		log.Printf("Handling %s - secure=%s", h.getPath(), strconv.FormatBool(secure))
		r.HandleFunc(h.getPath(), f)
	}
}