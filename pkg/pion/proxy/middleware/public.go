package middleware

import (
	"net/http"
	"strings"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/proxy/watchers"
	"github.com/labstack/echo"
)

const publicContextKey = "is-public-url"

// PublicPathsWatcher monitors changes in the public paths stored in etcd keys
var PublicPathsWatcher *watchers.PublicPathsWatcher

// isPublic checks if the request targets to the public URLs of files in buckets
func isPublic(c echo.Context) bool {
	public, ok := c.Get(publicContextKey).(bool)
	if ok && public {
		return true
	}

	// TODO Allow select methods on public urls. For now only support GET to these URLs
	if c.Request().Method != http.MethodGet {
		return false
	}

	requestPath := c.Request().URL.Path
	glog.V(2).Infof("Checking if prefix '%s' is public", requestPath)

	if PublicPathsWatcher == nil {
		glog.Errorf("Failed to initialize public-path manager")
		return false
	}

	publicPathPrefixes := PublicPathsWatcher.GetPublicPaths()
	for _, prefix := range publicPathPrefixes {
		if strings.HasPrefix(requestPath, prefix) {
			c.Set(publicContextKey, true)
			return true
		}
	}

	return false
}
