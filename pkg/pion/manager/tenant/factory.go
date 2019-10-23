package tenant

import (
	"path"

	"github.com/kpn/pion/pkg/pion/shared"
	"github.com/kpn/pion/pkg/pion/shared/multi-tenancy/etcd-store"
	"github.com/kpn/pion/pkg/pion/shared/multi-tenancy/validator"
	"go.etcd.io/etcd/clientv3"
)

// NewBucketStore returns the bucket store containing only buckets of the given customer.
func NewBucketStore(etcdClient *clientv3.Client, customerName string) etcd_store.BucketStore {
	return etcd_store.NewBucketStore(path.Join(shared.DefaultKeyPrefix, "buckets", customerName), etcdClient)
}

// NewUniqueBucketValidator retunrs the validator checking cross-tenant if the bucket name is unique
func NewUniqueBucketValidator(etcdClient *clientv3.Client) validator.BucketValidator {
	return validator.NewUniqueBucketValidator(path.Join(shared.DefaultKeyPrefix, "buckets"), etcdClient)
}
