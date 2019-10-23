package framework

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/golang/glog"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const DefaultRegion = "us-east-1"

const (
	timeout = 1 * time.Minute
)

type Framework struct {
	s3Endpoint string
}

func NewFramework() *Framework {
	s3Endpoint := os.Getenv("S3_ENDPOINT")
	if s3Endpoint == "" {
		glog.Fatal("Missing s3 endpoint config")
	}

	f := &Framework{
		s3Endpoint: s3Endpoint,
	}

	BeforeEach(f.BeforeEach)
	return f
}

func (f *Framework) BeforeEach() {
	Expect(flag.Set("stderrthreshold", "INFO")).To(BeNil())
}

func DescribeOSS(text string, body func()) bool {
	return Describe("[pion] "+text, body)
}

func (f *Framework) createS3Session(accessKey string, secretKey string) *session.Session {
	return session.Must(session.NewSession(&aws.Config{
		Region:           aws.String(DefaultRegion),
		Endpoint:         aws.String(f.s3Endpoint),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, "TOKEN"),
	}))
}

func printAWSErr(err error) error {
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				glog.Error(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
	}

	return err
}
