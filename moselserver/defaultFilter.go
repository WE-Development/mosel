package moselserver

import (
	"net/http"
	"github.com/bluedevel/mosel/commons"
)

type DefaultOptionsFilter struct {
}

func (filter DefaultOptionsFilter) Apply(w http.ResponseWriter, r *http.Request, next ApplyNext) {
	if r.Method == http.MethodOptions {
		commons.HttpNoContent(w)
	} else {
		next()
	}
}