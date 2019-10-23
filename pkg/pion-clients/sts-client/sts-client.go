package sts_client

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion-clients"
	"gopkg.in/resty.v1"
)

// CreateAccessKeyResponse returns the generated access/secret keys to client
type CreateAccessKeyResponse struct {
	AccessKey string    `json:"accessKey"`
	SecretKey string    `json:"secretKey"`
	CreatedAt time.Time `json:"createdAt"`
}

// KeyInfo struct only contains insensitive information, which can be displayed at UI
type KeyInfo struct {
	AccessKey string    `json:"accessKey"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}

// STSClient invokes STS service APIs
type STSClient interface {
	// List lists all existing access key of the given user
	ListAccessKeys(username string) ([]KeyInfo, error)

	// DeleteAccessKey removes the access key from STS service
	DeleteAccessKey(accessKey string) error

	// Create invokes STS service to generate a new access key with given lifetime and user's attributes
	CreateAccessKey(username string, lifetime time.Duration, attributes map[string]interface{}) (*CreateAccessKeyResponse, error)
}

// STSAddressURL variable refers to the remote STS service URL via the env-var STS_SERVICE_URL
var STSAddressURL = os.Getenv("STS_SERVICE_URL")

type defaultSTSClient struct {
	stsURLPrefix string
}

func New(customerName string) STSClient {
	if customerName == "" {
		glog.Fatal("Invalid initialization, customer name must not be empty")
	}
	return &defaultSTSClient{
		stsURLPrefix: STSAddressURL + "/customers/" + customerName,
	}
}

// List lists all existing access key of the given user
func (c defaultSTSClient) ListAccessKeys(username string) ([]KeyInfo, error) {
	if username == "" {
		return nil, errors.New("empty username")
	}
	resp, err := resty.R().
		Get(c.stsURLPrefix + "/users/" + username)
	if err != nil {
		return nil, err
	}

	body := resp.Body()
	if resp.StatusCode() != http.StatusOK {
		glog.Errorf("Cannot query username '%s'\nBody: %s", username, string(body))
		return nil, pion_clients.ErrInternalError
	}

	var keyInfos []KeyInfo
	err = json.Unmarshal(body, &keyInfos)
	if err != nil {
		glog.Errorf("Failed to parse response from STS::List(): %v", err)
		return nil, err
	}
	return keyInfos, nil
}

// DeleteAccessKey removes the access key from STS service
func (c defaultSTSClient) DeleteAccessKey(accessKey string) error {
	payload := map[string]string{
		"accessKey": accessKey,
	}
	resp, err := resty.R().SetBody(payload).Delete(c.stsURLPrefix + "/accesskey")
	if err != nil {
		return err
	}
	code := resp.StatusCode()
	switch code {
	case http.StatusNotFound:
		return pion_clients.ErrAccessKeyNotFound
	case http.StatusOK:
		return nil
	default:
		glog.Errorf("Failed to delete accesskey='%s': status='%d', body='%v'", accessKey, code, string(resp.Body()))
		return pion_clients.ErrInternalError
	}
}

// Create invokes STS service to generate a new access key with given lifetime and user's attributes
func (c defaultSTSClient) CreateAccessKey(username string, lifetime time.Duration, attributes map[string]interface{}) (*CreateAccessKeyResponse, error) {
	type Payload struct {
		UserId     string                 `json:"userId,omitempty"`
		Lifetime   string                 `json:"lifetime,omitempty"`
		Attributes map[string]interface{} `json:"attributes,omitempty"`
	}

	p := Payload{
		UserId:     username,
		Attributes: attributes,
		Lifetime:   lifetime.String(),
	}
	resp, err := resty.R().SetBody(p).Post(c.stsURLPrefix + "/accesskey")
	if err != nil {
		return nil, err
	}
	body := resp.Body()
	if resp.StatusCode() != http.StatusOK {
		glog.Errorf("Failed to create accesskey for user '%s': status='%d', body='%v'", username, resp.StatusCode(), string(body))
		return nil, pion_clients.ErrInternalError
	}
	var respObj CreateAccessKeyResponse
	err = json.Unmarshal(body, &respObj)
	if err != nil {
		glog.Errorf("Failed to parse response from STS::Create(): %v", err)
		return nil, err
	}
	return &respObj, nil
}
