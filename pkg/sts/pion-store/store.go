package pion_store

import (
	"encoding/json"
	"time"

	"github.com/kpn/pion/pkg/sts/token"
	etcd "go.etcd.io/etcd/clientv3"
)

type SecretData struct {
	SecretKey  string                 `json:"secretKey"`
	UserId     string                 `json:"userId"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
	CreatedAt  time.Time              `json:"createdAt,omitempty"` // default is IS8601 (RFC3339) date format
}

type KeyStore interface {
	SaveKey(username, accessKey string, lifetime time.Duration, data SecretData) error

	ListAccessKeys(username string) (keys []string, err error)

	DeleteAccessKey(accessKey string) error

	KeyQuerier
}

type etcdKeyStore struct {
	keyStore   token.Store
	keyQuerier KeyQuerier
}

// NewEtcdKeyStore create a key store for AccessKey/SecretKey. It is the decorator of the generic token.EtcdStore
func NewEtcdKeyStore(indexKeyPrefix, payloadKeyPrefix string, etcdClient *etcd.Client) (KeyStore, error) {
	ks, err := token.NewEtcdStore(indexKeyPrefix, payloadKeyPrefix, etcdClient)
	if err != nil {
		return nil, err
	}
	kq, err := NewKeyQuerier(payloadKeyPrefix, etcdClient)
	if err != nil {
		return nil, err
	}
	return &etcdKeyStore{
		keyStore:   ks,
		keyQuerier: kq,
	}, nil
}

func (s etcdKeyStore) Query(accessKey string) (data *SecretData, err error) {
	return s.keyQuerier.Query(accessKey)
}

func (s etcdKeyStore) SaveKey(username, accessKey string, lifetime time.Duration, data SecretData) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return s.keyStore.StoreToken(username, accessKey, string(dataBytes), lifetime)
}

func (s etcdKeyStore) ListAccessKeys(username string) (keys []string, err error) {
	return s.keyStore.GetUserTokens(username)
}

func (s etcdKeyStore) DeleteAccessKey(accessKey string) error {
	return s.keyStore.DeleteToken(accessKey)
}
