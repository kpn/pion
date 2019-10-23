package etcd_store

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"path"
	"time"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/shared/multi-tenancy/model"
	"github.com/kpn/pion/pkg/sts/cache"
	"go.etcd.io/etcd/clientv3"
)

// FileObjectStore defines API to manage file-object store
type FileObjectStore interface {
	Add(path string) (*model.FileObject, error)

	GetAllPaths() ([]string, error)

	Delete(path string) error
}
type fileObjectStore struct {
	keyPrefix string
	client    *clientv3.Client
}

// NewCustomerStore creates new store managing file-objects in Etcd
func NewFileObjectStore(keyPrefix string, etcdAddress string) (FileObjectStore, error) {
	etcdClient, err := cache.NewEtcdClient(etcdAddress)
	if err != nil {
		return nil, err
	}
	return &fileObjectStore{
		keyPrefix: keyPrefix,
		client:    etcdClient,
	}, nil
}

// Close shuts down the connection to Etcd cluster.
func (s fileObjectStore) Close() {
	err := s.client.Close()
	if err != nil {
		glog.Errorf("Closing Etcd client failed: %v", err)
	}
}

func (fos fileObjectStore) Add(path string) (*model.FileObject, error) {
	obj := model.FileObject{
		Path:      path,
		CreatedAt: time.Now().UTC(),
	}
	indexKey := fos.generateIndexKey(path)
	objStr, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	_, err = fos.client.Put(context.Background(), indexKey, string(objStr))
	return &obj, err
}

func (fos fileObjectStore) GetAllPaths() (paths []string, err error) {
	indexKeyPrefix := fos.generateIndexKey("")
	resp, err := fos.client.Get(context.Background(), indexKeyPrefix, clientv3.WithPrefix())
	if err != nil {
		glog.Errorf("Failed to list etcd keys with prefix '%s': %v", indexKeyPrefix, err)
		return nil, err
	}

	var obj model.FileObject
	for _, kv := range resp.Kvs {
		err = json.Unmarshal(kv.Value, &obj)
		if err != nil {
			glog.Errorf("Failed to unmarshal model.FileObject: %v", err)
			continue
		}
		paths = append(paths, obj.Path)
	}
	return paths, nil
}

func (fos fileObjectStore) Delete(path string) error {
	indexKey := fos.generateIndexKey(path)
	_, err := fos.client.Delete(context.Background(), indexKey)
	if err != nil {
		glog.Errorf("Failed to delete key '%s': %v", indexKey, err)
		return err
	}
	glog.V(1).Infof("Deleted file object path '%s'", indexKey)
	return nil
}

func (fos fileObjectStore) generateIndexKey(filePath string) string {
	if filePath == "" {
		return fos.keyPrefix + "/"
	}

	// hash the path to the unique id
	h := sha1.New()
	h.Write([]byte(filePath))
	bs := fmt.Sprintf("%x", h.Sum(nil))
	return path.Join(fos.keyPrefix, bs)
}
