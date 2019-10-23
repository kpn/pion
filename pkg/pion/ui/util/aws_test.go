package util_test

import (
	"testing"

	"github.com/kpn/pion/pkg/pion/ui/util"
)

func TestBucketNames(t *testing.T) {
	testCases := []struct {
		Name  string
		Valid bool
	}{
		{
			Name:  "bucketname",
			Valid: true,
		},
		{
			Name:  "nam",
			Valid: true,
		},
		{
			Name:  "12369873-4577-40df-ae2a-3e4323d5f478",
			Valid: true,
		},
		{
			Name:  "hyphen-.dottest",
			Valid: false,
		},
		{
			Name:  "dot.-hyphentest",
			Valid: false,
		},
		{
			Name:  "dot..dottest",
			Valid: false,
		},
		{
			Name:  "192.168.10.5",
			Valid: false,
		},
	}

	for i, c := range testCases {
		err := util.ValidateBucketName(c.Name)
		if c.Valid && err != nil {
			t.Errorf("test-case %d failed: bucketName='%s', expected valid, actual err='%v'", i, c.Name, err)
		}
		if !c.Valid && err == nil {
			if err == nil {
				t.Errorf("test-case %d failed: bucketName='%s', expected invalid, actual no error", i, c.Name)
			}
		}
	}
}
