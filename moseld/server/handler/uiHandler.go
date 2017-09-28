package handler

import (
	"net/http"
	"github.com/bluedevel/mosel/ui"
	"github.com/gorilla/mux"
)

type uiHandler struct {
	handler http.Handler
}

func NewUiHandler() uiHandler {
	return uiHandler{
		handler: http.FileServer(ui.AssetFS()),
	}
}

func (uiHandler uiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uiHandler.handler.ServeHTTP(w, r)
}

func (uiHandler uiHandler) ConfigureRoute(r *mux.Router, h http.Handler) {
	r.PathPrefix("/ui/").Handler(http.StripPrefix("/ui/", h))
}