package model

import (
	"encoding/json"
	"time"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/authz"
)

// Bucket struct defines a Bucket object information
type Bucket struct {
	Name       string        `json:"name"`
	OwnedBy    string        `json:"ownedBy"`              // name of customer owning the bucket
	CreatedAt  time.Time     `json:"createdAt,omitempty"`  // default is IS8601 (RFC3339) date format
	ModifiedAt time.Time     `json:"modifiedAt,omitempty"` // default is IS8601 (RFC3339) date format
	Creator    string        `json:"creator"`
	ACLs       authz.ACLList `json:"acls"`
}

// NewBucketFromBytes creates Bucket object from serialized byte array.
func NewBucketFromBytes(data []byte) (*Bucket, error) {
	var obj Bucket
	err := json.Unmarshal(data, &obj)
	if err != nil {
		glog.Errorf("Failed to unmarshal object: '%v'", string(data))
		return nil, err
	}
	return &obj, nil
}
