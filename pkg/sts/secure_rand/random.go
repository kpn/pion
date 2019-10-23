package secure_rand

import (
	"crypto/rand"
	"fmt"
)

func randomKey(keyLength int) ([]byte, error) {
	key := make([]byte, keyLength)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// SecureRandomString generates a random key with given keyLength in bytes
func SecureRandomString(keyLength int) (string, error) {
	key, err := randomKey(keyLength)
	if err != nil {
		return "", err
	}
	id := fmt.Sprintf("%0x", key[0:keyLength])
	return id, nil
}
