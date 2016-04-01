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
	"net/http"
	"github.com/bluedevel/mosel/api"
	"encoding/json"
)

type loginHandler struct {
}

func (handler loginHandler) ServeHTTPContext(ctx MoselServerContext, w http.ResponseWriter, r *http.Request) {
	resp := api.NewLoginResponse()

	user, passwd, enabled := r.BasicAuth()

	if !enabled || !ctx.auth.Authenticate(user, passwd) {
		resp.Successful = false
		//resp.ValidTo
	} else {
		resp.Successful = true
	}

	json.NewEncoder(w).Encode(resp)
}

func (handler loginHandler) getPath() string {
	return "/login"
}

func (handler loginHandler) Secure() bool {
	return false
}