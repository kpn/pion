package cmd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/minio/minio/pkg/auth"
)

// Decorator to expose AWS Signature V4 verification functions
type AWSV4Verifier struct {
	creds  auth.Credentials
	region string
}

func NewAWSV4Verifier(accessKey string, secretKey string, region string) (*AWSV4Verifier, error) {
	creds, err := auth.CreateCredentials(accessKey, secretKey)
	if err != nil {
		return nil, fmt.Errorf("unable create credential, %s", err)
	}

	return &AWSV4Verifier{
		creds:  creds,
		region: region,
	}, nil
}

func (verifier AWSV4Verifier) GetCredential() auth.Credentials {
	return verifier.creds
}

func (verifier AWSV4Verifier) IsReqAuthenticated(ctx context.Context, r *http.Request, region string) (s3Error APIErrorCode) {
	return verifier.isReqAuthenticated(ctx, r, region)
}

func (verifier AWSV4Verifier) CheckAdminRequestAuthType(ctx context.Context, r *http.Request, region string) APIErrorCode {
	return verifier.checkAdminRequestAuthType(ctx, r, region)
}

func (verifier AWSV4Verifier) GetRegion() string {
	return verifier.region
}

// AccessKeyFromRequest returns the accessKey used in the request
func AccessKeyFromRequest(req *http.Request, region string) (string, APIErrorCode) {
	v4Auth := req.Header.Get("Authorization")

	// Parse signature version '4' header.
	signV4Values, errCode := parseSignV4(v4Auth, region)
	if errCode != ErrNone {
		return "", errCode
	}

	return signV4Values.Credential.accessKey, ErrNone

}
