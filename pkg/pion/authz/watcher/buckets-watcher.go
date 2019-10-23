package watcher

import (
	"context"

	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/shared/multi-tenancy/model"
	etcd "go.etcd.io/etcd/clientv3"
)

// BucketsWatcher defines API for watching bucket changes in Etcd
type BucketsWatcher interface {
	GetBucket(name string) *model.Bucket

	Watch()
}

type bucketsWatcher struct {
	bucketPathPrefix string
	bucketMap        BucketMap // map from bucket name to bucket object
	client           *etcd.Client
}

// NewBucketsWatcher creates a watcher that load existing ACLs from buckets
func NewBucketsWatcher(bucketPathPrefix string, client *etcd.Client) (BucketsWatcher, error) {
	mgr := bucketsWatcher{
		bucketPathPrefix: bucketPathPrefix,
		client:           client,
	}
	err := mgr.init()
	if err != nil {
		return nil, err
	}
	return &mgr, nil
}

func (mgr *bucketsWatcher) init() error {
	bm, err := loadBuckets(mgr.bucketPathPrefix, mgr.client)
	if err != nil {
		glog.Errorf("Loading ACLs from '%s' failed: %v", mgr.bucketPathPrefix, err)
		return err
	}
	mgr.bucketMap = bm
	glog.V(2).Infof("Initial buckets: %v", mgr.bucketMap)
	return nil
}

func (mgr bucketsWatcher) GetBucket(name string) *model.Bucket {
	return mgr.bucketMap[name]
}

// Watch monitors latest ACLs in buckets
func (mgr bucketsWatcher) Watch() {
	watchChan := mgr.client.Watch(context.Background(), mgr.bucketPathPrefix, etcd.WithPrefix(), etcd.WithPrevKV())
	go func() {
		for {
			resp, ok := <-watchChan
			if !ok || resp.Err() != nil {
				if ok {
					glog.Errorf("Watching channel returns: %v", resp.Err())
					return
				}
				glog.Warningf("Watching channel closed")
				return
			}
			for _, evt := range resp.Events {
				switch evt.Type {
				case mvccpb.PUT:
					bkt, err := model.NewBucketFromBytes((*evt.Kv).Value)
					if err != nil {
						continue
					}
					mgr.onBucketChanged(*bkt)
				case mvccpb.DELETE:
					glog.V(3).Infof("deleted key: %s", string((*evt.PrevKv).Key))
					glog.V(3).Infof("deleted value: %s", string((*evt.PrevKv).Value))
					bkt, err := model.NewBucketFromBytes((*evt.PrevKv).Value)
					if err != nil {
						continue
					}
					glog.V(2).Infof("Delete ACLs of bucket: %v", bkt)
					mgr.onBucketDeleted(*bkt)

				default:
					glog.V(2).Infof("Unknown Etcd event: %v", *evt)
				}
			}

			if glog.V(3) {
				glog.Info("Bucket store: ")
				for name, bkt := range mgr.bucketMap {
					glog.Infof("%s: %v", name, bkt)
				}
			}
		}
	}()
}

func (mgr *bucketsWatcher) onBucketChanged(bkt model.Bucket) {
	glog.V(2).Infof("Updating info of bucket '%s'", bkt.Name)
	mgr.bucketMap[bkt.Name] = &bkt
}

func (mgr *bucketsWatcher) onBucketDeleted(bkt model.Bucket) {
	glog.V(2).Infof("Removing bucket '%s'", bkt.Name)
	delete(mgr.bucketMap, bkt.Name)
}
