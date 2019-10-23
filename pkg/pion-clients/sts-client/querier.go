package sts_client

import (
	"encoding/json"
	"net/http"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion-clients"
	"gopkg.in/resty.v1"
)

type AccessKeyQuerier interface {
	// Query returns the single generic token from STS service
	Query(accessKey string) (*STSResponse, error)
}

type querier struct {
	urlPrefix string
}

// NewAccessKeyQuerier creates an object for querying access key
func NewAccessKeyQuerier() AccessKeyQuerier {
	return &querier{
		urlPrefix: STSAddressURL,
	}
}

// Query returns the single access key from STS service
func (c querier) Query(accessKey string) (*STSResponse, error) {
	resp, err := resty.R().Get(c.urlPrefix + "/accesskey/" + accessKey)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() == http.StatusNotFound {
		return nil, pion_clients.ErrAccessKeyNotFound
	}
	body := resp.Body()
	if resp.StatusCode() != http.StatusOK {
		glog.Errorf("Cannot query access key '%s'\nBody: %s", accessKey, string(body))
		return nil, pion_clients.ErrInternalError
	}

	var payload STSResponse
	err = json.Unmarshal(body, &payload)
	if err != nil {
		glog.Errorf("Failed to parse response from STS::GetAccessKey(): %v", err)
		return nil, pion_clients.ErrInternalError
	}
	return &payload, nil
}
