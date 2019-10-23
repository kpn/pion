package util

import (
	"errors"
	"os"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var UpstreamAddress = os.Getenv("MINIO_SERVICE_URL")

// CreateS3Bucket makes a new bucket in the Upstream S3 server defined at $MINIO_SERVICE_URL
func CreateS3Bucket(bucketName string) error {
	if UpstreamAddress == "" {
		return errors.New("missing env-var 'MINIO_SERVICE_URL' setting")
	}

	accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
	secretKeyID := os.Getenv("MINIO_SECRET_KEY")
	if accessKeyID == "" || secretKeyID == "" {
		return errors.New("missing access keys configurations 'MINIO_ACCESS_KEY' and 'MINIO_SECRET_KEY'")
	}

	sess := session.Must(session.NewSession())
	svc := s3.New(sess, &aws.Config{
		Region:           aws.String(endpoints.UsEast1RegionID),
		Endpoint:         aws.String(UpstreamAddress),
		Credentials:      credentials.NewStaticCredentials(accessKeyID, secretKeyID, ""),
		S3ForcePathStyle: aws.Bool(true),
	})
	_, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	return err
}

// ValidateBucketName checks if the bucketName follows restrictions of S3 bucket name in
// https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-s3-bucket-naming-requirements.html
func ValidateBucketName(bucketName string) error {
	const ipAddressRegexValue = `^(\d+\.)+\d+$`
	const nameRegexValue = `^(([a-z0-9]|[a-z0-9][a-z0-9\-]*[a-z0-9])\.)*([a-z0-9]|[a-z0-9][a-z0-9\-]*[a-z0-9])$`

	l := len(bucketName)
	if l < 3 || l > 63 {
		return errors.New("name must be between 3-63 characters")
	}
	ipRegex := regexp.MustCompile(ipAddressRegexValue)
	if ipRegex.MatchString(bucketName) {
		return errors.New("name must not be IP address style")
	}
	nameRegex := regexp.MustCompile(nameRegexValue)
	if !nameRegex.MatchString(bucketName) {
		return errors.New("invalid bucket naming")
	}
	return nil
}
