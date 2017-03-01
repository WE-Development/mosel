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
	"github.com/bluedevel/mosel/moselnoded/server/context"
	"net/http"
	"io/ioutil"
	"github.com/gorilla/mux"
	"log"
)

type scriptHandler struct {
	ctx *context.MoselnodedServerContext
}

func NewScriptHandler(ctx *context.MoselnodedServerContext) scriptHandler {
	return scriptHandler{
		ctx: ctx,
	}
}

func (handler scriptHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	name := vars["script"]

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return
	}

	err = handler.ctx.Collector.AddScript(name, b)

	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (handler scriptHandler) GetPath() string {
	return "/script/{script}"
}

func (handler scriptHandler) Secure() bool {
	return true
}