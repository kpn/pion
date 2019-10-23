package pion_store

import "github.com/kpn/pion/pkg/pion/shared"

const (
	IndexKeyPrefix   = shared.DefaultKeyPrefix + "/tokens/indices" // IndexKeyPrefix is the Etcd key prefix storing index keys
	PayloadKeyPrefix = shared.DefaultKeyPrefix + "/tokens/secrets" // PayloadKeyPrefix is the Etcd key prefix storing secret keys
)
