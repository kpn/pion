package secure_rand

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomString(t *testing.T) {
	for i := 1; i <= 10; i++ {
		key, err := SecureRandomString(32)
		assert.NoError(t, err)
		assert.Equal(t, 32*8, len(key)*4) // as each hexa char is 4 bits
		fmt.Printf("Round %d:%s\n", i, key)
	}
}
