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
package handler

import (
	"net/http"
	"encoding/json"

	"github.com/bluedevel/mosel/api"
	"github.com/bluedevel/mosel/moseld/server/context"
)

// Handler for providing general information about the server instance. See api for details.
type infoHandler struct {
	ctxd *context.MoseldServerContext
}

func NewInfoHandler(ctxd *context.MoseldServerContext) *infoHandler {
	return &infoHandler{
		ctxd: ctxd,
	}
}

func (handler infoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := api.NewInfoResponse()
	resp.Nodes = handler.ctxd.DataCache.GetNodes()
	json.NewEncoder(w).Encode(resp)
}

func (handler infoHandler) GetPath() string {
	return "/info"
}

func (handler infoHandler) Secure() bool {
	return true
}
