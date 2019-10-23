package verifier

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	aws_v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	minio "github.com/minio/minio/cmd"
)

func buildRequest(method, url, body string, t *testing.T) (*http.Request, io.ReadSeeker) {
	reader := strings.NewReader(body)
	return buildRequestWithBodyReader(method, url, reader, t)
}

func buildRequestWithBodyReader(method, urlStr string, body io.Reader, t *testing.T) (*http.Request, io.ReadSeeker) {
	var bodyLen int
	if method == "" {
		method = "POST"
	}

	type lenner interface {
		Len() int
	}
	if lr, ok := body.(lenner); ok {
		bodyLen = lr.Len()
	}

	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		t.Fatalf("Cannot create request: %v", err)
	}
	// req.URL.Opaque = "//example.org/bucket/key-._~,!@#$%^&*()"
	req.Header.Set("X-Amz-Target", "prefix.Operation")
	req.Header.Set("Content-Type", "application/x-amz-json-1.0")

	if bodyLen > 0 {
		req.Header.Set("Content-Length", strconv.Itoa(bodyLen))
	}

	req.Header.Set("X-Amz-Meta-Other-Header", "some-value=!@#$%^&* (+)")
	req.Header.Add("X-Amz-Meta-Other-Header_With_Underscore", "some-value=!@#$%^&* (+)")

	var seeker io.ReadSeeker
	if sr, ok := body.(io.ReadSeeker); ok {
		seeker = sr
	} else {
		t.Fatalf("Cannot typecast")
	}

	return req, seeker

}

func buildSigner(accessKey, secretKey string) aws_v4.Signer {
	return aws_v4.Signer{
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, "SESSION"),
	}
}

func buildAndSignRequest(t *testing.T, accessKey string, secretKey string, region string, method string, url string, bodyStr string) *http.Request {
	req, body := buildRequest(method, url, bodyStr, t)
	signer := buildSigner(accessKey, secretKey)

	_, err := signer.Sign(req, body, "s3", region, time.Now())
	if err != nil {
		t.Fatalf("failed to sign: %v", err)
	}
	return req
}

func TestAWSSignatureV4(t *testing.T) {
	region := "us-east-1"
	verifier, err := minio.NewAWSV4Verifier("myuser", "mypassword", region)
	if err != nil {
		t.Fatalf("Cannot create AWS signature-v4 verifier: %v", err)
	}

	testCases := []struct {
		Request     *http.Request
		HashContent string
		ErrCode     minio.APIErrorCode
	}{
		{
			buildAndSignRequest(t, "myuser", "mypassword", region,
				"GET", "http://localhost:9000/mybucket", ""),
			"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			minio.ErrNone,
		},
		{
			buildAndSignRequest(t, "myuser", "invalid-password", region,
				"GET", "http://localhost:9000/mybucket", "hello"),
			"2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
			minio.ErrSignatureDoesNotMatch,
		},
		{
			buildAndSignRequest(t, "invalid-user", "password", region,
				"POST", "http://localhost:9000/mybucket", "hello"),
			"2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
			minio.ErrInvalidAccessKeyID,
		},
	}
	ctx := context.Background()
	for i, testCase := range testCases {
		if e, a := testCase.HashContent, testCase.Request.Header.Get("X-Amz-Content-Sha256"); e != a {
			t.Errorf("Test %d: Invalid hash body, expected %v, got %v", i, e, a)
		}
		if s3Error := verifier.CheckAdminRequestAuthType(ctx, testCase.Request, region); s3Error != testCase.ErrCode {
			t.Errorf("Test %d: Unexpected s3error returned wanted %d, got %d", i, testCase.ErrCode, s3Error)
		}
	}

}
