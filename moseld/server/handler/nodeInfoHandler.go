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
	"github.com/bluedevel/mosel/moseld/server/context"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/bluedevel/mosel/api"
	"encoding/json"
	"log"
	"strconv"
)

type nodeInfoHandler struct {
	ctxd *context.MoseldServerContext
}

func NewNodeInfoHandler(ctxd *context.MoseldServerContext) *nodeInfoHandler {
	return &nodeInfoHandler{
		ctxd:ctxd,
	}
}

func (handler nodeInfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	node := vars["node"]

	points, _ := handler.ctxd.Cache.GetAll(node)
	resp := api.NewNodeInfoResponse()

	for _, point := range points {
		var stamp string = strconv.FormatInt(point.Time.Unix(), 10)
		resp.Data[stamp] = point.Info
	}

	log.Println(resp)
	json.NewEncoder(w).Encode(resp)
}

func (handler nodeInfoHandler) GetPath() string {
	return "/nodeInfo/{node}"
}

func (handler nodeInfoHandler) Secure() bool {
	return false
}
