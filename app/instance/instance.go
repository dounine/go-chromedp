package instance

import (
	"github.com/patrickmn/go-cache"
	"time"
)

type CacheFile struct {
	Path     *string
	Content  []byte
	FileName string
}

var (
	DownloadCaches = cache.New(1*time.Minute, 3*time.Minute)
)
