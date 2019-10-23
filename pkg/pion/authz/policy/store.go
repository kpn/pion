package policy

import "github.com/kpn/pion/pkg/pion/authz"

// Store defines APIs for the policy store
type Store interface {
	// Get returns list of ACLs of the bucket, nil if no policy is found
	Get(bucketName string) authz.ACLList
	// Set replaces old ACLs of the bucket (if existed) with the new list
	Set(bucketName string, acls authz.ACLList) error
	// Add appends new ACLs to the existing ACLs of the bucket
	Add(bucketName string, acls authz.ACLList) error
	// Delete removes ACLs of the specific bucket
	Delete(bucketName string)

	// SerializePolicies returns the serialized string of policies
	SerializePolicies() (string, error)
}

type defaultStore struct {
	policies authz.PoliciesType
}

// TODO remove when no usage
func NewStore(policies authz.PoliciesType) Store {
	return &defaultStore{
		policies: policies,
	}
}

// Get returns list of ACLs of the bucket, nil if no policy is found
func (s defaultStore) Get(bucketName string) authz.ACLList {
	return s.policies[bucketName]
}

// Set replaces old ACLs of the bucket (if existed) with the new list
func (s defaultStore) Set(bucketName string, acls authz.ACLList) error {
	s.policies[bucketName] = acls
	return nil
}

// Add appends new ACLs to the existing ACLs of the bucket
func (s defaultStore) Add(bucketName string, acls authz.ACLList) error {
	oldAcls := s.policies[bucketName]
	if oldAcls == nil {
		oldAcls = []authz.ACL{}
	}
	s.policies[bucketName] = append(oldAcls, acls...)
	return nil
}

// Delete removes ACLs of the specific bucket
func (s defaultStore) Delete(bucketName string) {
	delete(s.policies, bucketName)
}

// SerializePolicies returns the serialized string of policies
func (s defaultStore) SerializePolicies() (string, error) {
	return s.policies.String()
}
