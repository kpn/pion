package token

import (
	"errors"
	"time"
)

const KeyLength = 32

var NotFoundError = errors.New("token not found")

// Store contains methods to manage tokens
type Store interface {
	// Generate creates a random token for user with payload and saves to store
	Generate(username string, data string, lifetime time.Duration) (string, error)

	// StoreToken saves a token and data to store with limited lifetime
	StoreToken(username, token, data string, lifetime time.Duration) error

	// DeleteToken removes a sigle token from store
	DeleteToken(t string) error

	// GetData returns the stored data of the token
	GetData(token string) (data string, err error)

	// GetUserTokens returns list of tokens of the given user
	GetUserTokens(username string) ([]string, error)

	// DeleteUserTokens removes all tokens of the user
	DeleteUserTokens(username string) error
}
