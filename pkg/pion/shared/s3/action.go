package s3

import (
	"fmt"
	"net/http"

	"github.com/kpn/pion/pkg/pion/shared/rbac"
	"github.com/pkg/errors"
)

type ActionType string

const (
	CreateBucket ActionType = "create-bucket"
	DeleteBucket ActionType = "delete-bucket"
	ListBucket   ActionType = "list-bucket"
	GetObject    ActionType = "get-object"
	DeleteObject ActionType = "delete-object"
	CreateObject ActionType = "create-object"
)

func ActionFromString(str string) (ActionType, error) {
	var actionMap = map[string]ActionType{
		"create-bucket": CreateBucket,
		"delete-bucket": DeleteBucket,
		"list-bucket":   ListBucket,
		"get-object":    GetObject,
		"delete-object": DeleteObject,
		"create-object": CreateObject,
	}
	a := actionMap[str]
	if a == "" {
		return "", errors.New("Invalid s3 action")
	}
	return a, nil
}

func ToS3Action(httpMethod, keyPath string) (ActionType, error) {
	if keyPath == "" {
		// bucket actions
		switch httpMethod {
		case http.MethodPut:
			return CreateBucket, nil
		case http.MethodDelete:
			return DeleteBucket, nil
		case http.MethodGet:
			return ListBucket, nil
		default:
			return "", fmt.Errorf("unknow S3 action on bucket: httpMethod='%s' ", httpMethod)
		}
	}

	// object actions
	var actionMap = map[string]ActionType{
		http.MethodPut:    CreateObject,
		http.MethodPost:   CreateObject,
		http.MethodDelete: DeleteObject,
		http.MethodGet:    GetObject,
		http.MethodHead:   GetObject,
	}
	action := actionMap[httpMethod]
	if action == "" {
		return "", fmt.Errorf("unknown equivalent S3 action on object: httpMethod='%s'", httpMethod)
	}
	return action, nil
}

func ToRbacAction(s3Action ActionType) rbac.Action {
	var actionMap = map[ActionType]rbac.Action{
		CreateBucket: rbac.Create,
		DeleteBucket: rbac.Delete,
		ListBucket:   rbac.List,
		CreateObject: rbac.Update,
		DeleteObject: rbac.Update,
		GetObject:    rbac.Get,
	}
	return actionMap[s3Action]
}
