package handlers

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	aws_v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/proxy"
	"github.com/kpn/pion/pkg/pion/proxy/utils"
)

var (
	masterAccessKey = os.Getenv("MINIO_ACCESS_KEY")
	masterSecretKey = os.Getenv("MINIO_SECRET_KEY")
)

// PopulateRequest change protocol scheme and address to the upstream, and if request to private resources, it needs
// signing
func PopulateRequest(r *http.Request) (*http.Request, error) {
	body, contentLength, err := readBody(r.Body)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = r.Body.Close()
		if err != nil {
			glog.Errorf("[proxy error]: %v", err)
		}
	}()

	url := UpstreamAddress + r.RequestURI
	glog.V(2).Infof("Target URL=%s", url)

	newReq, err := http.NewRequest(r.Method, url, body)
	if err != nil {
		return nil, err
	}
	newReq.ContentLength = contentLength
	newReq.Header = utils.CopyWhiteListHeaders(r.Header)

	if contentLength > 0 {
		newReq.Header.Set("Content-Length", strconv.FormatInt(newReq.ContentLength, 10))
	} else {
		newReq.Header.Del("Content-Length")
	}

	err = signWithBody(newReq, body)
	if err != nil {
		return nil, err
	}
	glog.V(3).Infof("Content-Length after populating: %v", newReq.ContentLength)

	return newReq, nil
}

func readBody(body io.Reader) (io.ReadSeeker, int64, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		glog.Errorf("Cannot read all request body: %v", err)
		return nil, 0, err
	}
	length := int64(len(data))
	return bytes.NewReader(data), length, nil
}

func signWithBody(req *http.Request, body io.ReadSeeker) error {
	signer := buildSigner(masterAccessKey, masterSecretKey)

	_, err := signer.Sign(req, body, "s3", proxy.DefaultRegion, time.Now())
	if err != nil {
		glog.Errorf("Failed to sign request: %v", err)
		return err
	}
	return nil
}

func buildSigner(accessKey, secretKey string) aws_v4.Signer {
	return aws_v4.Signer{
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	}
}
