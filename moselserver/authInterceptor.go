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
)

func (server MoselServer) secure(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if server.Config.Sessions.Enabled {
			key := r.FormValue("key")
			if key == "" || !server.Context.Sessions.ValidateSession(key) {
				commons.HttpUnauthorized(w)
				return
			}
		} else {
			user, passwd, enabled := r.BasicAuth()

			if !server.Config.AuthTrue.Enabled &&
				(!enabled || !server.Context.Auth.Authenticate(user, passwd)) {
				commons.HttpUnauthorized(w)
				return
			}
		}

		fn(w, r)
	}
}