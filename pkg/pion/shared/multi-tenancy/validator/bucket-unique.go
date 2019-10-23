package validator

import (
	"context"
	"fmt"
	"path"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	etcd "go.etcd.io/etcd/clientv3"
)

var (
	ErrInternal      = errors.New("internal error")
	ErrBucketExisted = errors.New("BucketAlreadyExists")
)

// BucketValidator is the generic validator for bucket
type BucketValidator interface {
	Validate(bucketName string) error
}

type uniqueBucketValidator struct {
	keyPrefix string
	client    *etcd.Client
}

// NewUniqueBucketValidator is the validator checking if the bucket name is unique
func NewUniqueBucketValidator(keyPrefix string, client *etcd.Client) BucketValidator {
	return &uniqueBucketValidator{
		keyPrefix: keyPrefix,
		client:    client,
	}
}

// Validate checks if the bucket name has already existed in the cross-customers scope
func (validator uniqueBucketValidator) Validate(bucketName string) error {
	// Improve with list keys using pagination, i.e. etcd.WithLimit(n)
	resp, err := validator.client.Get(context.Background(), validator.keyPrefix, etcd.WithPrefix(), etcd.WithKeysOnly())
	if err != nil {
		glog.Errorf("Failed to list etcd keys with prefix '%s': %v", validator.keyPrefix, err)
		return ErrInternal
	}
	for _, kv := range resp.Kvs {
		keyString := string(kv.Key)
		currentBucketName, err := bucketNameFromKey(keyString)
		if err != nil {
			glog.Error(err.Error())
			continue
		}
		if currentBucketName == bucketName {
			return ErrBucketExisted
		}
	}
	return nil
}

func bucketNameFromKey(etcdKey string) (string, error) {
	// return the last element of the etcd-key
	name := path.Base(etcdKey)
	if name == "" {
		return "", fmt.Errorf("no bucket name found in the etcd key '%s'", etcdKey)
	}
	return name, nil
}
