package s3

import (
	"errors"
	"strings"

	"github.com/golang/glog"
)

// GetResources split the requestURI to the s3 resources (bucket name, key-path).
func GetResources(requestURI string) (bucketName string, keyPath string, err error) {
	if !strings.HasPrefix(requestURI, "/") {
		glog.Errorf("Invalid URI, URI must start with '/': %s", requestURI)
		return bucketName, keyPath, errors.New("invalid URI, URI must start with '/'")
	}

	paths := strings.Split(requestURI, "/")
	if len(paths) < 2 {
		glog.Errorf("Invalid URI: %s", requestURI)
		return bucketName, keyPath, errors.New("invalid URI, bucket not found")
	}

	bucketName = paths[1]
	keyPath = strings.Join(paths[2:], "/")
	return bucketName, keyPath, nil
}
