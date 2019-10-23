package framework

import (
	"context"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/golang/glog"
)

func (f *Framework) UploadObject(accessKey, secretKey, srcFilePath, targetBucketName, targetKeyPath string) (location string, err error) {
	sess := f.createS3Session(accessKey, secretKey)
	uploader := s3manager.NewUploader(sess)

	glog.Info("Uploading file to S3")

	ctx, cancelFn := context.WithTimeout(context.Background(), timeout)
	defer cancelFn()

	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		glog.Errorf("Failed to open file '%s': %v", srcFilePath, err)
		return "", err
	}
	defer srcFile.Close()

	result, err := uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(targetBucketName),
		Key:    aws.String(path.Join(targetKeyPath, path.Base(srcFilePath))),
		Body:   srcFile,
	})
	if err != nil {
		return "", err
	}
	return result.Location, nil
}

func (f *Framework) DeleteObject(accessKey, secretKey, bucketName, keyPath string) error {
	svc := s3.New(f.createS3Session(accessKey, secretKey))
	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyPath),
	})

	return err
}
