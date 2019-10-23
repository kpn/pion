package watcher

import (
	"context"
	"encoding/json"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/shared/multi-tenancy/model"
	etcd "go.etcd.io/etcd/clientv3"
)

// BucketMap is the map type from bucket name to Bucket object
type BucketMap map[string]*model.Bucket

// loadBuckets loads Buckets  at given etcd key path
func loadBuckets(bucketsKeyPrefix string, etcdClient *etcd.Client) (bucketMap BucketMap, err error) {
	resp, err := etcdClient.Get(context.Background(), bucketsKeyPrefix, etcd.WithPrefix())
	if err != nil {
		glog.Errorf("Failed to list etcd keys with prefix '%s': %v", bucketsKeyPrefix, err)
		return nil, err
	}
	bucketMap = make(BucketMap)

	for _, kv := range resp.Kvs {
		var bkt model.Bucket
		err = json.Unmarshal(kv.Value, &bkt)
		if err != nil {
			glog.Errorf("Failed to unmarshal model.Bucket object: %v", err)
			continue
		}
		glog.V(3).Infof("Got model.Bucket '%s'", bkt.Name)
		bucketMap[bkt.Name] = &bkt
	}
	return bucketMap, nil
}
