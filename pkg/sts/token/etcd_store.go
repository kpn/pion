package token

import (
	"context"
	"errors"
	"path"
	"time"

	"github.com/golang/glog"
	rand "github.com/kpn/pion/pkg/sts/secure_rand"
	etcd "go.etcd.io/etcd/clientv3"
)

type etcdStore struct {
	indexKeyPrefix   string
	payloadKeyPrefix string
	client           *etcd.Client
}

// NewEtcdStore creates a Store managing access tokens. A token `x` is stored at both place:
//
// - $(indexKeyPrefix)/users/${u}/x: this key is used for indexing and listing tokens of user `u`.
//
// - $(payloadKeyPrefix)/x: this key contains actual payload data of the token. It's used for fast lookup tokens
func NewEtcdStore(indexKeyPrefix, payloadKeyPrefix string, etcdClient *etcd.Client) (Store, error) {
	return &etcdStore{
		indexKeyPrefix:   indexKeyPrefix,
		payloadKeyPrefix: payloadKeyPrefix,
		client:           etcdClient,
	}, nil
}

func (s etcdStore) Generate(username string, data string, lifetime time.Duration) (token string, err error) {
	token, err = rand.SecureRandomString(KeyLength)
	if err != nil {
		return "", err
	}

	err = s.StoreToken(username, token, data, lifetime)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s etcdStore) DeleteToken(token string) error {
	payload, err := s.getPayload(token)
	if err != nil {
		return err
	}
	if payload.Username == "" || token == "" {
		return errors.New("empty username")
	}
	indexKey := s.generateIndexKey(payload.Username, token)

	payloadKey := s.generatePayloadKey(token)
	_, err = s.client.Delete(context.Background(), payloadKey)
	if err != nil {
		glog.Errorf("Failed to delete key '%s': %v", payloadKey, err)
		return err
	}

	_, err = s.client.Delete(context.Background(), indexKey)
	if err != nil {
		glog.Errorf("Failed to delete key '%s': %v", indexKey, err)
		return err
	}
	glog.V(2).Infof("Deleted keys '%s' and '%s'", payloadKey, indexKey)
	return nil
}

func (s etcdStore) GetUserTokens(username string) (tokens []string, err error) {
	if username == "" {
		return nil, errors.New("empty username")
	}
	indexKeyPrefix := s.generateIndexKey(username, "")
	resp, err := s.client.Get(context.Background(), indexKeyPrefix, etcd.WithPrefix())
	if err != nil {
		glog.Errorf("Failed to list etcd keys with prefix '%s': %v", indexKeyPrefix, err)
		return nil, err
	}

	for _, kv := range resp.Kvs {
		tokens = append(tokens, string(kv.Value))
	}
	return tokens, nil
}

func (s etcdStore) DeleteUserTokens(username string) error {
	return errors.New("not implemented")
}

func (s etcdStore) GetData(token string) (data string, err error) {
	p, err := s.getPayload(token)
	if err != nil {
		return "", err
	}
	return p.Data, nil
}

// getPayload returns the whole payload value stored in DB
func (s etcdStore) getPayload(token string) (*Payload, error) {
	tokenKey := s.generatePayloadKey(token)

	resp, err := s.client.Get(context.Background(), tokenKey)
	if err != nil {
		glog.Errorf("Failed to get key '%s': %v", tokenKey, err)
		return nil, err
	}
	if len(resp.Kvs) < 1 {
		return nil, errors.New("token refers to empty value")
	}
	kv := resp.Kvs[0]
	glog.V(2).Infof("Get key '%s', value=%s", kv.Key, kv.Value)

	p, err := NewPayload(string(kv.Value))
	if err != nil {
		glog.Errorf("Failed to parse token payload: %v", err)
		return nil, err
	}
	return &p, nil
}

// StoreToken adds a new key to $(prefix)/$(token) with lifetime.
func (s etcdStore) StoreToken(username string, token string, data string, lifetime time.Duration) (err error) {
	if username == "" || token == "" {
		return errors.New("empty username or token")
	}
	tokenKey := s.generatePayloadKey(token)

	// indexKey is used for looking up from username to a list of tokens
	indexKey := s.generateIndexKey(username, token)

	leaseGrantResp, err := s.client.Grant(context.Background(), int64(lifetime.Seconds()))
	if err != nil {
		return err
	}

	resp, err := s.client.Put(context.Background(), indexKey, token, etcd.WithLease(leaseGrantResp.ID))
	if err != nil {
		glog.Errorf("Failed to store key '%s', error=%v, resp=%v", indexKey, err, resp)
		return err
	}

	p := Payload{
		Username: username,
		Data:     data,
	}

	resp, err = s.client.Put(context.Background(), tokenKey, p.String(), etcd.WithLease(leaseGrantResp.ID))
	if err != nil {
		glog.Errorf("Failed to store key '%s', error=%v, resp=%v", tokenKey, err, resp)
		return err
	}

	glog.V(2).Infof("Stored key '%s' and '%s'", indexKey, tokenKey)
	return err
}
func (s etcdStore) generatePayloadKey(token string) string {
	return path.Join(s.payloadKeyPrefix, token)
}

func (s etcdStore) generateIndexKey(username string, token string) string {
	if s.indexKeyPrefix == "" {
		glog.Fatal("invalid setup: when indexKeyPrefix is empty, this method must not be used")
	}
	return path.Join(s.indexKeyPrefix, "users", username, token)
}
