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
	"github.com/bluedevel/mosel/moselserver"
	"net/http"
	"time"
	"log"
	"github.com/bluedevel/mosel/api"
	"encoding/json"
	"fmt"
)

type streamHandler struct {
	ctx *moselserver.MoselServerContext
}

func NewStreamHandler(ctx *moselserver.MoselServerContext) streamHandler {
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
		data, err := json.Marshal(handler.createResponse(r, now))

		if err != nil {
			log.Println(err)
			ticker.Stop()
		}

		_, err = w.Write(data)

		if err != nil {
			log.Println(err)
			ticker.Stop()
		}

		fmt.Fprint(w, "\n")

		if err == nil {
			flusher.Flush()
		} else {
			log.Println(err)
			ticker.Stop()
		}

	}

}

func (handler streamHandler) createResponse(r *http.Request, now time.Time) interface{} {
	resp := api.NewNodeResponse()

	return resp
}

func (handler streamHandler) GetPath() string {
	return "/stream"
}

func (handler streamHandler) Secure() bool {
	return false
}