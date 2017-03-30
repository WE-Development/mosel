package context

import (
	"time"

	"github.com/bluedevel/mosel/api"
)

type DataPersistence interface {
	Init() error
	Add(node string, t time.Time, info api.NodeInfo)
	GetAll() (DataCacheStorage, error)
	GetAllSince(since time.Duration) (DataCacheStorage, error)
}
