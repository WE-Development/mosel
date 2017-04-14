package moselserver

import "net/http"

type Filter interface {
	Apply(http.ResponseWriter, *http.Request, ApplyNext)
}

type ApplyNext func()
