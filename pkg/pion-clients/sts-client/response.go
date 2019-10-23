package sts_client

import (
	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/ui/multi_tenant"
)

// STSResponse is the response from STS::GetAccessKey() API, which is cloned from pkg/sts/pion-store/store.go:10
type STSResponse struct {
	SecretKey  string                 `json:"secretKey"`
	UserId     string                 `json:"userId"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

// GetUserGroups returns the groups attribute in the STS response
func (resp STSResponse) GetUserGroups() []string {
	arr, ok := resp.Attributes[multi_tenant.UserAttributeGroups].([]interface{})
	if !ok {
		glog.Errorf("Failed to get groups attribute")
		return nil
	}
	var groups []string
	for _, v := range arr {
		groups = append(groups, v.(string))
	}
	return groups
}

// GetUserCustomer returns the customer attribute in the STS response
func (resp STSResponse) GetUserCustomer() string {
	customer, ok := resp.Attributes[multi_tenant.UserAttributeCustomer].(string)
	if !ok {
		glog.Errorf("Failed to get customer attribute")
		return ""
	}
	return customer
}
