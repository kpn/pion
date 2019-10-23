package request

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/authz"
	"github.com/kpn/pion/pkg/pion/shared"
	"github.com/kpn/pion/pkg/pion/shared/rbac"
	"github.com/kpn/pion/pkg/pion/shared/s3"
	"github.com/pkg/errors"
)

// Build creates rbac-request and authorization-request from a http request with authorization headers
func Build(httpRequest *http.Request) (*rbac.Request, *authz.Request, error) {
	userId, groups := GetUserAttributes(httpRequest)

	s3Action, err := GetS3Action(httpRequest)
	if err != nil {
		glog.Error(err.Error())
		return nil, nil, err
	}
	azAction, err := ToAuthzAction(s3Action)
	if err != nil {
		glog.Error(err.Error())
		return nil, nil, err
	}

	customer, err := GetCustomerName(httpRequest)
	if err != nil {
		glog.Errorf(err.Error())
		return nil, nil, err
	}

	bucketName, keyPath, err := GetResources(httpRequest)
	if err != nil {
		glog.Errorf(err.Error())
		return nil, nil, err
	}

	return &rbac.Request{
			UserID: userId,
			Groups: groups,
			Action: s3.ToRbacAction(s3Action),
			Target: getRBACTarget(keyPath),
		}, &authz.Request{
			Username: userId,
			Groups:   groups,
			Action:   azAction,
			Target:   bucketName,
			Customer: customer,
		}, nil
}

func getRBACTarget(keyPath string) rbac.Resource {
	if keyPath == "" {
		return rbac.BucketResource
	} else {
		return rbac.ObjectResource
	}
}

// GetResources returns the specific resourceIDs (bucket name and key path) in the http request
func GetResources(r *http.Request) (bucketName string, keyPath string, err error) {
	resourceId := r.Header.Get(shared.ResourceKey)
	return s3.GetResources(resourceId)
}

// GetS3Action returns the S3 action type specified in the header 'X-Action'
func GetS3Action(r *http.Request) (s3.ActionType, error) {
	str := r.Header.Get(shared.ActionKey)
	glog.V(2).Infof("action header: %v", str)

	return s3.ActionFromString(str)
}

// ToAuthzAction converts from S3Action type to ACL action type
func ToAuthzAction(s3Action s3.ActionType) (authz.ActionType, error) {
	var aclActionMap = map[s3.ActionType]authz.ActionType{
		s3.CreateBucket: authz.Write,
		s3.DeleteBucket: authz.Write,
		s3.ListBucket:   authz.Read,
		s3.GetObject:    authz.Read,
		s3.DeleteObject: authz.Write,
		s3.CreateObject: authz.Write,
	}
	authzAction, ok := aclActionMap[s3Action]
	if !ok {
		return "", fmt.Errorf("unknown s3 action: %s", s3Action)
	}
	return authzAction, nil
}

// GetCustomerName returns customer name set in the http request header
func GetCustomerName(r *http.Request) (string, error) {
	name := r.Header.Get(shared.CustomerKey)
	if name == "" {
		return "", errors.New("Customer identifier not found in the request")
	}
	return name, nil
}

// GetUserAttributes returns userID and userGroups set in the http request headers
func GetUserAttributes(r *http.Request) (id string, groups []string) {
	id = r.Header.Get(shared.UserIdKey)
	groupsStr := r.Header.Get(shared.UserGroupKey)
	groups = strings.Split(groupsStr, ",")
	return id, groups

}
