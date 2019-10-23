package framework

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func (f *Framework) ListBuckets(accessKey, secretKey string) (buckets []string, err error) {
	svc := s3.New(f.createS3Session(accessKey, secretKey))

	ctx, cancelFn := context.WithTimeout(context.Background(), timeout)
	defer cancelFn()

	result, err := svc.ListBucketsWithContext(ctx, &s3.ListBucketsInput{})
	printAWSErr(err)
	if err != nil {
		return nil, err
	}

	buckets = []string{}
	for _, bkt := range result.Buckets {
		buckets = append(buckets, *bkt.Name)
	}
	return buckets, nil
}

func (f *Framework) CreateBucket(accessKey, secretKey, bucketName string) (err error) {
	svc := s3.New(f.createS3Session(accessKey, secretKey))

	ctx, cancelFn := context.WithTimeout(context.Background(), timeout)
	defer cancelFn()

	_, err = svc.CreateBucketWithContext(ctx, &s3.CreateBucketInput{Bucket: aws.String(bucketName)})
	return printAWSErr(err)
}

func (f *Framework) DeleteBucket(accessKey, secretKey, bucketName string) (err error) {
	svc := s3.New(f.createS3Session(accessKey, secretKey))

	ctx, cancelFn := context.WithTimeout(context.Background(), timeout)
	defer cancelFn()

	_, err = svc.DeleteBucketWithContext(ctx, &s3.DeleteBucketInput{Bucket: aws.String(bucketName)})
	return printAWSErr(err)
}

func (f *Framework) DeleteAllBucketItems(accessKey, secretKey, bucketName string) (err error) {
	// TODO this function does not work well with Pion
	svc := s3.New(f.createS3Session(accessKey, secretKey))
	iter := s3manager.NewDeleteListIterator(svc, &s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
	})
	if err := s3manager.NewBatchDeleteWithClient(svc).Delete(aws.BackgroundContext(), iter); err != nil {
		return fmt.Errorf("unable to delete objects from bucket %q: %v", bucketName, err)
	}
	return nil
}
