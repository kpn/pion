package etcd_store

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"time"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/authz"
	"github.com/kpn/pion/pkg/pion/shared/multi-tenancy/model"
	"github.com/pkg/errors"
	etcd "go.etcd.io/etcd/clientv3"
)

// BucketStore defines API to manage bucket store
type BucketStore interface {
	Add(b model.Bucket) (*model.Bucket, error)

	Delete(name string) error

	Get(name string) (*model.Bucket, error)

	List() ([]model.Bucket, error)

	// TODO implement bucket metadata and store ACLs there
	UpdateACLs(bucketName string, acls authz.ACLList) error
}

type bucketStore struct {
	keyPrefix string
	client    *etcd.Client
}

// NewBucketStore creates new store object managing buckets in Etcd
func NewBucketStore(keyPrefix string, client *etcd.Client) BucketStore {
	return bucketStore{
		keyPrefix: keyPrefix,
		client:    client,
	}
}

func (store bucketStore) Add(b model.Bucket) (*model.Bucket, error) {
	indexKey := store.generateIndexKey(b.Name)
	if store.exists(indexKey) {
		return nil, fmt.Errorf("bucket '%s' already existed", b.Name)
	}

	b.CreatedAt = time.Now().UTC()

	bytes, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}

	_, err = store.client.Put(context.Background(), indexKey, string(bytes))
	if err != nil {
		glog.Errorf("Failed to save new model.Bucket '%s': %v", b.Name, err)
		return nil, err
	}
	glog.V(1).Infof("Added new model.Bucket '%s'", b.Name)
	return &b, nil
}

func (store bucketStore) UpdateACLs(bucketName string, acls authz.ACLList) error {
	bucket, err := store.Get(bucketName)
	if err != nil {
		return err
	}
	bucket.ACLs = acls
	bucket.ModifiedAt = time.Now().UTC()

	bytes, err := json.Marshal(bucket)
	if err != nil {
		return err
	}

	indexKey := store.generateIndexKey(bucketName)
	_, err = store.client.Put(context.Background(), indexKey, string(bytes))
	if err != nil {
		glog.Errorf("Failed to update bucket ACLs '%s': %v", bucket.Name, err)
		return err
	}
	glog.V(1).Infof("Update ACLs of bucket '%s'", bucket.Name)
	return nil
}

func (store bucketStore) Delete(name string) error {
	indexKey := store.generateIndexKey(name)
	resp, err := store.client.Delete(context.Background(), indexKey)
	if err != nil {
		glog.Errorf("Failed to delete model.Bucket '%s': %v", name, err)
		return err
	}
	if resp.Deleted != 1 {
		err = fmt.Errorf("key '%s' not found", indexKey)
		glog.Error(err.Error())
		return err
	}
	glog.V(1).Infof("Deleted model.Bucket '%s'", name)
	return nil
}

func (store bucketStore) Get(name string) (*model.Bucket, error) {
	indexKey := store.generateIndexKey(name)
	resp, err := store.client.Get(context.Background(), indexKey)
	if err != nil {
		glog.Errorf("Failed to get etcd key '%s': %v", indexKey, err)
		return nil, err
	}

	if len(resp.Kvs) < 1 {
		glog.Errorf("Reading etcd K-V error at '%s': should have 1 key-value", indexKey)
		return nil, errors.New("Reading error: K-V does not exist")
	}
	var obj model.Bucket
	err = json.Unmarshal(resp.Kvs[0].Value, &obj)
	if err != nil {
		glog.Errorf("Failed to unmarshal model.Bucket object: %v", err)
		return nil, err
	}
	glog.V(2).Infof("Get model.Bucket '%s: %v'", name, obj)
	return &obj, nil
}

func (store bucketStore) List() (buckets []model.Bucket, err error) {
	indexKeyPrefix := store.generateIndexKey("")
	resp, err := store.client.Get(context.Background(), indexKeyPrefix, etcd.WithPrefix())
	if err != nil {
		glog.Errorf("Failed to list etcd keys with prefix '%s': %v", indexKeyPrefix, err)
		return nil, err
	}
	for _, kv := range resp.Kvs {
		var obj model.Bucket
		err = json.Unmarshal(kv.Value, &obj)
		if err != nil {
			glog.Errorf("Failed to unmarshal model.Bucket object: %v", err)
			continue
		}
		buckets = append(buckets, obj)
	}
	glog.V(3).Infof("Got buckets '%v'", buckets)
	return
}

func (store bucketStore) generateIndexKey(name string) string {
	if name == "" {
		return store.keyPrefix + "/"
	}

	return path.Join(store.keyPrefix, name)
}

func (store bucketStore) exists(indexKey string) bool {
	resp, err := store.client.Get(context.Background(), indexKey, etcd.WithCountOnly())
	if err != nil {
		glog.Errorf("Failed to check key existence of '%s': %v", indexKey, err)
		return false
	}
	glog.V(2).Infof("Returned %d keys at '%s", resp.Count, indexKey)
	return resp.Count == 1
}
