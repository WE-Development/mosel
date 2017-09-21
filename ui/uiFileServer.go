package ui

import (
	"net/http"
	"github.com/elazarl/go-bindata-assetfs"
)

func NewUiFileServer() http.Handler {

	/*names,_ := AssetDir("")

	log.Println(len(names))

	for _, name := range names {
		data,_ := Asset(name)
		log.Println(string(data[:]))
	}*/

	return http.FileServer(assetFS())
}

func AssetFS() *assetfs.AssetFS {
	return assetFS()
}