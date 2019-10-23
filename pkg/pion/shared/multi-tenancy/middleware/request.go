package middleware

import (
	"net/http"
	"strings"

	"github.com/kpn/pion/pkg/pion/shared"
)

// getUserAttributes returns the userId and userGroups in the HTTP request headers X-User-Id and X-User-Groups
func getUserAttributes(r *http.Request) (id string, groups []string) {
	id = r.Header.Get(shared.UserIdKey)
	groupsStr := r.Header.Get(shared.UserGroupKey)
	groups = strings.Split(groupsStr, ",")
	return id, groups

}
