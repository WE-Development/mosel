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
	"fmt"
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
	fmt.Fprintf(w, "Hallo Welt")
}

func (handler nodeInfoHandler) GetPath() string {
	return "/nodeInfo"
}

func (handler nodeInfoHandler) Secure() bool {
	return false
}
