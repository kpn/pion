package watchers

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/shared/multi-tenancy/model"
	etcd "go.etcd.io/etcd/clientv3"
)

const (
	publicPathPrefix = "/pion/oss/public-objects"
)

// TODO This implementation does not following micro-services design pattern. Might be better to use event-driven paradigm:
// TODO i.e. pion-manager publishes event upon adding/removing paths, pion-proxy subscribes event to update changes
// PublicPaths object is used to get latest list of public paths in the Etcd by monitoring changes of the given Etcd key.
type PublicPathsWatcher struct {
	publicPaths []string
	client      *etcd.Client
}

func NewPublicPathsWatcher(client *etcd.Client) *PublicPathsWatcher {
	return &PublicPathsWatcher{
		publicPaths: nil,
		client:      client,
	}
}

func (p *PublicPathsWatcher) Init() error {
	// Read all existing public paths
	resp, err := p.client.Get(context.Background(), publicPathPrefix, etcd.WithPrefix())
	if err != nil {
		glog.Errorf("reading etcd key '%s' failed: %v", publicPathPrefix, err)
		return err
	}
	for _, kv := range resp.Kvs {
		obj, err := newFileObject(kv.Value)
		if err != nil {
			continue
		}
		p.publicPaths = append(p.publicPaths, obj.Path)
	}
	glog.V(2).Infof("Initial public paths: %v", p.publicPaths)
	return nil
}

func (p PublicPathsWatcher) GetPublicPaths() []string {
	return p.publicPaths
}

// Watch monitors latest changes in the list of public paths
func (p *PublicPathsWatcher) Watch() {
	watchChan := p.client.Watch(context.Background(), publicPathPrefix, etcd.WithPrefix(), etcd.WithPrevKV())
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
					obj, err := newFileObject((*evt.Kv).Value)
					if err != nil {
						continue
					}
					glog.V(2).Infof("Add new item: %v", obj)
					p.addPublicPath(obj.Path)
				case mvccpb.DELETE:
					glog.V(3).Infof("deleted key: %s", string((*evt.PrevKv).Key))
					glog.V(3).Infof("deleted value: %s", string((*evt.PrevKv).Value))
					obj, err := newFileObject((*evt.PrevKv).Value)
					if err != nil {
						continue
					}
					glog.V(2).Infof("Delete item: %v", obj)
					err = p.removePublicPath(obj.Path)
					if err != nil {
						glog.Errorf("error to remove public path '%s': %v", obj.Path, err)
					}
				default:
					glog.V(2).Infof("Unknown Etcd event: %v", *evt)
				}
			}

			glog.V(2).Infof("Public paths: %v", p.publicPaths)
		}
	}()
}

func (p *PublicPathsWatcher) removePublicPath(deletedPath string) error {
	var deletedIndex = -1
	for i, item := range p.publicPaths {
		if item == deletedPath {
			deletedIndex = i
			break
		}
	}
	if deletedIndex < 0 {
		return errors.New("item not found")
	}
	// remove item deletedIndex-th
	p.publicPaths = append(p.publicPaths[:deletedIndex], p.publicPaths[deletedIndex+1:]...)
	return nil
}

func (p *PublicPathsWatcher) addPublicPath(path string) {
	p.publicPaths = append(p.publicPaths, path)
}

func newFileObject(data []byte) (*model.FileObject, error) {
	var obj model.FileObject
	err := json.Unmarshal(data, &obj)
	if err != nil {
		glog.Errorf("Failed to unmarshal object: '%v'", string(data))
		return nil, err
	}
	return &obj, nil
}
