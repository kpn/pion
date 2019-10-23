package manager_client

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/shared/multi-tenancy/model"
	"gopkg.in/resty.v1"
)

// CustomerClient interface defines Customer APIs
type CustomerClient interface {
	Get(name string) (*model.Customer, error)
}

type customerClient struct {
	serverURL string
}

// NewCustomerClient makes a client calling Customer APIs
func NewCustomerClient(serverURL string) CustomerClient {
	return &customerClient{
		serverURL: serverURL,
	}
}

func (client customerClient) Get(name string) (*model.Customer, error) {
	if name == "" {
		return nil, errors.New("empty name")
	}
	resp, err := resty.R().
		Get(client.serverURL + "/_internal/customers/" + name)
	if err != nil {
		return nil, err
	}

	body := resp.Body()
	switch resp.StatusCode() {
	case http.StatusNotFound:
		return nil, ErrCustomerNotfound
	case http.StatusOK:
		var customer model.Customer
		err = json.Unmarshal(body, &customer)
		if err != nil {
			glog.Errorf("Failed to parse response from Manager::CustomerBuckets(): %v", err)
			return nil, err
		}
		return &customer, nil
	}

	glog.Errorf("Cannot get customer '%s, status='%d", name, resp.StatusCode())
	if glog.V(2) {
		glog.Infof("Dumped body: %v", string(body))
	}
	return nil, errors.New("Cannot get customer")
}
