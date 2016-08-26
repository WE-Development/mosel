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
	"fmt"
	"github.com/WE-Development/mosel/moseld/server/context"
	"log"
)

type debugHandler struct {
	ctxd *context.MoseldServerContext
}

func NewDebugHandler(ctxd *context.MoseldServerContext) debugHandler {
	return debugHandler{
		ctxd: ctxd,
	}
}

func (handler debugHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var flusher http.Flusher

	if f, ok := w.(http.Flusher); ok {
		flusher = f
	} else {
		return
	}

	for now := range time.Tick(1 * time.Second) {

		log.Println(handler.ctxd.Debug)
		handler.ctxd.Debug++

		fmt.Fprintln(w, now)
		flusher.Flush()
	}
}

func (handler debugHandler) GetPath() string {
	return "/debug"
}

func (handler debugHandler) Secure() bool {
	return false
}
