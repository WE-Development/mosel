package context

import "github.com/WE-Development/mosel/moselserver"

type MoselnodedServerContext struct {
	moselserver.MoselServerContext

	Collector *collector
}
