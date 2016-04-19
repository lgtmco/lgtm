package middleware

import (
	"github.com/lgtmco/lgtm/version"

	"github.com/gin-gonic/gin"
)

// Version is a middleware function that appends the LGTM version information
// to the HTTP response. This is intended for debugging and troubleshooting.
func Version(c *gin.Context) {
	c.Header("X-LGTM-VERSION", version.Version)
	c.Next()
}
