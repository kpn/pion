package s3_test

import (
	"testing"

	"github.com/kpn/pion/pkg/pion/shared/s3"
	"github.com/stretchr/testify/assert"
)

func TestGetBucketName(t *testing.T) {
	testCases := []struct {
		URI        string
		BucketName string
		KeyPath    string
		Error      bool
	}{
		{
			"/foo",
			"foo",
			"",
			false,
		},
		{
			"/foo/",
			"foo",
			"",
			false,
		},
		{
			"/",
			"",
			"",
			false,
		},
		{
			"",
			"",
			"",
			true,
		},
		{
			"/foo/bah",
			"foo",
			"bah",
			false,
		},
	}
	for i, c := range testCases {
		actualBucketName, actualKeyPath, err := s3.GetResources(c.URI)
		if c.Error {
			assert.Error(t, err, "Test-case %d, expect error but not", i)
		} else {
			assert.NoError(t, err, "Test-case %d, expect no error but occurred", i)
			assert.Equal(t, c.BucketName, actualBucketName, "Test-case %d, expect bucket name '%s' but got '%s'", i, c.BucketName, actualBucketName)
			assert.Equal(t, c.KeyPath, actualKeyPath, "Test-case %d, expect key-path '%s' but got '%s'", i, c.KeyPath, actualKeyPath)
		}
	}

}
