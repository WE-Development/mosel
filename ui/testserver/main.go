package main

import (
	"net/http"
	"log"
	"github.com/bluedevel/mosel/ui"
)

func main() {
	http.Handle("/", http.FileServer(ui.AssetFS()))
	log.Panic(http.ListenAndServe(":8080", nil))
}
