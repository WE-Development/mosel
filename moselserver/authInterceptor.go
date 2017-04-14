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
	"regexp"

	"github.com/bluedevel/mosel/commons"
	"github.com/bluedevel/mosel/config"
)

type authFilter struct {
	server *MoselServer
}

func (filter authFilter) Apply(w http.ResponseWriter, r *http.Request, next ApplyNext) {
	if filter.server.Config.Sessions.Enabled {
		filter.authSession(w, r, next)
	} else {
		filter.authDirect(w, r, next)
	}
}

func newAuthFilter(server *MoselServer) *authFilter {
	return &authFilter{
		server: server,
	}
}

// Authenticate via sessions
func (filter *authFilter) authSession(w http.ResponseWriter, r *http.Request, next ApplyNext) {
	key := r.FormValue("key")
	if key == "" || !filter.server.Context.Sessions.ValidateSession(key) {
		commons.HttpUnauthorized(w)
		return
	}

	next()
}

// Authenticate directly via basic auth
func (filter *authFilter) authDirect(w http.ResponseWriter, r *http.Request, next ApplyNext) {
	server := filter.server;
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
	if server.Config.AuthTrue.Enabled {
		next()
		return
	}

	if ok, err := filter.validateAccessRights(r.URL.Path, server.Config.Users[user], server.Config.Groups); !ok {
		if err != nil {
			log.Println(err)
		}

		commons.HttpUnauthorized(w)
		return
	}

	next()
}

func (filter *authFilter) validateAccessRights(path string, userConfig *moselconfig.UserConfig, groupConfigs map[string]*moselconfig.GroupConfig) (bool, error) {
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
