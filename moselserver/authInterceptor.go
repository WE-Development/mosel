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
	"github.com/WE-Development/mosel/config"
	"regexp"
)

// Wrap a http.HandleFunc such that it's authenticated before execution.
// It's build on basic auth or use with sessions as provided by moselserver.sessionCache.
func (server *MoselServer) secure(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if server.Config.Sessions.Enabled {
			authSession(server, fn, w, r)
		} else {
			authDirect(server, fn, w, r)
		}
	}
}

// Authenticate via sessions
func authSession(server *MoselServer, fn http.HandlerFunc, w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	if key == "" || !server.Context.Sessions.ValidateSession(key) {
		commons.HttpUnauthorized(w)
		return
	}
	fn(w, r)
}

// Authenticate directly via basic auth
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

	// if auth true is enabled, it's possible that a positiv authentication occurs without a user configuration
	if !server.Config.AuthTrue.Enabled {
		if ok, err := validateAccessRights(r.URL.Path, server.Config.Users[user], server.Config.Groups); !ok {
			if err != nil {
				log.Println(err)
			}

			commons.HttpUnauthorized(w)
			return
		}
	}

	fn(w, r)
}

func validateAccessRights(path string, userConfig *moselconfig.UserConfig, groupConfigs map[string]*moselconfig.GroupConfig) (bool, error) {
	allow := false;

	rights := make([]moselconfig.AccessRights, 0)

	for groupName, groupConf := range groupConfigs {
		for _, userGroup := range userConfig.Groups {
			if userGroup != groupName {
				continue
			}
			rights = append(rights, groupConf.AccessRights)
		}
	}
	rights = append(rights, userConfig.AccessRights)

	for _, rightConf := range rights {
		var err error
		var match bool
		for _, denyRegex := range rightConf.Deny {
			match, err = regexp.MatchString(denyRegex, path)

			if err != nil {
				return false, err
			}

			allow = allow && !match
		}
		for _, allowRegex := range rightConf.Allow {
			match, err = regexp.MatchString(allowRegex, path)

			if err != nil {
				return false, err
			}

			allow = allow || match
		}

	}

	return allow, nil
}