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
	"github.com/WE-Development/mosel/commons"
	"log"
)

func (server *MoselServer) secure(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if server.Config.Sessions.Enabled {
			authSession(server, fn, w, r)
		} else {
			authDirect(server, fn, w, r)
		}
	}
}

func authSession(server *MoselServer, fn http.HandlerFunc, w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	if key == "" || !server.Context.Sessions.ValidateSession(key) {
		commons.HttpUnauthorized(w)
		return
	}
	fn(w, r)
}

func authDirect(server *MoselServer, fn http.HandlerFunc, w http.ResponseWriter, r *http.Request) {
	user, passwd, _ := r.BasicAuth()

	if server.Context.Auth == nil {
		log.Println("No autheticator configured! Denying all requests")
		commons.HttpUnauthorized(w)
		return
	}

	if !server.Context.Auth.Authenticate(user, passwd) {
		commons.HttpUnauthorized(w)
		return
	}

	fn(w, r)
}