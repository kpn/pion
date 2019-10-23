package rbac

import (
	"encoding/json"
	"time"

	"github.com/golang/glog"
)

type SubjectType string

const (
	UserType  SubjectType = "user"
	GroupType SubjectType = "group"
)

type Subject struct {
	Type  SubjectType `json:"type"`
	Value string      `json:"value"`
}

// RoleBinding defines the link from Role with specific permissions to the given subjects, which can be either users or group of users
type RoleBinding struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt,omitempty"` // default is IS8601 (RFC3339) date format
	Subjects  []Subject `json:"subjects"`
	RoleRef   string    `json:"roleRef"`
}

// NewRoleBindingFromBytes deserializes RoleBinding object from byte array
func NewRoleBindingFromBytes(data []byte) (*RoleBinding, error) {
	var obj RoleBinding
	err := json.Unmarshal(data, &obj)
	if err != nil {
		glog.Errorf("Failed to unmarshal object: '%v'", string(data))
		return nil, err
	}
	return &obj, nil
}
