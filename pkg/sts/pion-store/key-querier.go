package pion_store

import (
	"encoding/json"

	"github.com/kpn/pion/pkg/sts/token"
	"go.etcd.io/etcd/clientv3"
)

// KeyQuerier interface is to query secret data from given access key
type KeyQuerier interface {
	// Query returns the stored secret data binding to the access key
	Query(accessKey string) (data *SecretData, err error)
}

type keyQuerier struct {
	tokenStore token.Store
}

func NewKeyQuerier(payloadKeyPrefix string, etcdClient *clientv3.Client) (KeyQuerier, error) {
	ts, err := token.NewEtcdStore("", payloadKeyPrefix, etcdClient)
	if err != nil {
		return nil, err
	}
	return &keyQuerier{
		tokenStore: ts,
	}, nil
}

// Query returns the stored secret data binding to the access key
func (kq keyQuerier) Query(accessKey string) (*SecretData, error) {
	dataStr, err := kq.tokenStore.GetData(accessKey)
	if err != nil {
		return nil, err
	}
	var dataObj SecretData
	if err = json.Unmarshal([]byte(dataStr), &dataObj); err != nil {
		return nil, err
	}

	return &dataObj, nil
}
