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
	"time"
	"encoding/json"
	"log"

	"github.com/bluedevel/mosel/api"
	"github.com/bluedevel/mosel/moselnoded/server/context"
)

type streamHandler struct {
	ctx  *context.MoselnodedServerContext

	test int
}

func NewStreamHandler(ctx *context.MoselnodedServerContext) streamHandler {
	return streamHandler{
		ctx: ctx,
	}
}

func (handler streamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var flusher http.Flusher

	if f, ok := w.(http.Flusher); ok {
		flusher = f
	} else {
		return
	}

	ticker := time.NewTicker(1 * time.Second)
	for now := range ticker.C {
		err := json.NewEncoder(w).Encode(handler.createResponse(r, now))

		if err != nil {
			log.Println(err)
			ticker.Stop()
			break
		}

		flusher.Flush()
	}

}

func (handler *streamHandler) createResponse(r *http.Request, now time.Time) interface{} {
	resp := api.NewNodeResponse()
	handler.ctx.Collector.FillNodeInfo(&resp.NodeInfo)
	return resp
}

func (handler streamHandler) GetPath() string {
	return "/stream"
}

func (handler streamHandler) Secure() bool {
	return true
}