package s3responses

import (
	"time"

	"github.com/kpn/pion/pkg/pion/shared/multi-tenancy/model"
)

type ListAllMyBucketsResult struct {
	Owner   Owner   `xml:"Owner"`
	Buckets Buckets `xml:"Buckets"`
}

type Owner struct {
	ID          string `xml:"ID"`
	DisplayName string `xml:"DisplayName"`
}

type Buckets struct {
	Bucket []Bucket `xml:"Bucket"`
}

type Bucket struct {
	Location     string `xml:"Location"`
	Name         string `xml:"Name"`
	CreationDate string `xml:"CreationDate"`
}

func NewListBucketsResult(buckets []model.Bucket) (*ListAllMyBucketsResult, error) {
	var result ListAllMyBucketsResult
	for _, bkt := range buckets {
		result.Buckets.Bucket = append(result.Buckets.Bucket, Bucket{
			Name:         bkt.Name,
			CreationDate: bkt.CreatedAt.Format(time.RFC3339),
		})
	}
	return &result, nil
}
