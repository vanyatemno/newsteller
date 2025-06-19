package cache

import (
	"github.com/puzpuzpuz/xsync/v4"
	"go.uber.org/zap"
	"strings"
)

/**
A cache for basic pages, which implements reverse caching logic.
*/

type Event string

const (
	PostsUpdated Event = "post"
)

type PagesCache struct {
	pages *xsync.Map[string, string]
}

func NewPagesCache() *PagesCache {
	return &PagesCache{
		pages: xsync.NewMap[string, string](),
	}
}

func (c *PagesCache) Set(key string, value string) {
	c.pages.Store(key, value)
}

func (c *PagesCache) Get(key string) (string, bool) {
	return c.pages.Load(key)
}

func (c *PagesCache) Invalidate(event Event) {
	switch event {
	case PostsUpdated:
		c.invalidateForPosts()
	default:
		zap.L().Warn("invalid event", zap.String("event", string(event)))
	}
}

func (c *PagesCache) invalidateForPosts() {
	c.pages.Range(func(k, v string) bool {
		if strings.Contains(k, "/home") || strings.Contains(k, "/posts") {
			zap.L().Info("invalidated page with key:", zap.String("key", k))
			c.pages.Delete(k)
		}
		return true
	})
}
