package authz

import (
	"encoding/json"

	"github.com/golang/glog"
)

type PoliciesType map[string]ACLList

func NewPolicies(data []byte) (policies PoliciesType, err error) {
	err = json.Unmarshal(data, &policies)
	if err != nil {
		glog.Errorf("Failed to reload policies: %v", err)
		return nil, err
	}
	return policies, nil
}

func (policies PoliciesType) String() (string, error) {
	output, err := json.Marshal(policies)
	if err != nil {
		return "", err
	}
	return string(output), nil
}
