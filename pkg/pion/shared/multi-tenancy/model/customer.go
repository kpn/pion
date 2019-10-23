package model

import (
	"time"
)

// Customer struct represents a customer, that can manage multiple groups of users
type Customer struct {
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"createdAt,omitempty"`  // default is IS8601 (RFC3339) date format
	ModifiedAt time.Time `json:"modifiedAt,omitempty"` // default is IS8601 (RFC3339) date format
	Groups     []string  `json:"groups"`               // array of user groups binding to the customer account
	UserIDs    []string  `json:"userIDs"`              // array of individual userIDs binding to the customer account

	// TODO Minio server endpoint for each customer, which supports multi minio servers at back-end
}

// ContainsUser checks if the given userID is in the customer's individual UserIDs list
func (c Customer) ContainsUser(userID string) bool {
	for _, u := range c.UserIDs {
		if u == userID {
			return true
		}
	}
	return false
}

// ContainsGroup checks if the given group-name is in the customer
func (c Customer) ContainsGroup(groupName string) bool {
	for _, g := range c.Groups {
		if g == groupName {
			return true
		}
	}
	return false
}

// HasAnyGroup checks if any of given groups belongs to the customer
func (c Customer) HasAnyGroup(checkingGroups []string) bool {
	var mapGroups = make(map[string]bool)
	for _, g := range c.Groups {
		mapGroups[g] = true
	}

	for _, cg := range checkingGroups {
		if mapGroups[cg] {
			return true
		}
	}
	return false
}
