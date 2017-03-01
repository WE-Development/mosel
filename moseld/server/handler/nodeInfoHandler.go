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
	"strconv"
	"time"
	"github.com/bluedevel/mosel/commons"
)

// Handler for providing time coded information on a node.
// It also supports queries for entries not older than a given unix-timestamp (Seconds since 01.01.1970).
// Use: ?since=<timestamp>
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

	var points []context.DataPoint
	var err error

	since := r.URL.Query().Get("since")
	if since == "" {
		points, err = handler.ctxd.DataCache.GetAll(node)
	} else {
		//test this with stamp=$(($(date +%s)-10)); curl http://localhost:8282/nodeInfo/self\?since\=${stamp}
		var i int64
		i, err = strconv.ParseInt(since, 10, 64)
		if err != nil {
			commons.HttpBadRequest(w)
			return
		}

		points, err = handler.ctxd.DataCache.GetSince(node, time.Unix(i, 0))
	}
	commons.HttpCheckError(err, http.StatusInternalServerError, w)

	resp := api.NewNodeInfoResponse()

	for _, point := range points {
		var stamp string = strconv.FormatInt(point.Time.Unix(), 10)
		resp.Data[stamp] = point.Info
	}

	json.NewEncoder(w).Encode(resp)
}

func (handler nodeInfoHandler) GetPath() string {
	return "/nodeInfo/{node}"
}

func (handler nodeInfoHandler) Secure() bool {
	return true
}
