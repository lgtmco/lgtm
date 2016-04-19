package middleware

import (
	"time"

	"github.com/bradrydzewski/lgtm/cache"

	"github.com/gin-gonic/gin"
	"github.com/ianschenck/envflag"
)

var (
	ttl = envflag.Duration("CACHE_TTL", time.Minute*15, "")
)

func Cache() gin.HandlerFunc {
	cache_ := cache.NewTTL(*ttl)
	return func(c *gin.Context) {
		c.Set("cache", cache_)
		c.Next()
	}
}
