package handler

import (
	"github.com/WE-Development/mosel/moselnoded/server/context"
	"net/http"
	"io/ioutil"
	"github.com/gorilla/mux"
)

type scriptHandler struct {
	ctx *context.MoselnodedServerContext
}

func NewScriptHandler (ctx *context.MoselnodedServerContext) scriptHandler {
	return scriptHandler{
		ctx: ctx,
	}
}

func (handler scriptHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method != "PUT" {
		return
	}

	vars := mux.Vars(r)
	name := vars["script"]

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return
	}

	script := string(b)
	handler.ctx.Collector.AddScript(name, script)
}

func (handler scriptHandler) GetPath() string {
	return "/script/{script}"
}

func (handler scriptHandler) Secure() bool {
	return true
}