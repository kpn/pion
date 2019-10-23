package manager_client

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion-clients"
	"github.com/kpn/pion/pkg/pion/authz"
	"github.com/kpn/pion/pkg/pion/shared"
	"github.com/kpn/pion/pkg/pion/shared/multi-tenancy/model"
	"github.com/kpn/pion/pkg/pion/shared/rbac"
	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
)

var (
	// TODO Put this url config via app args
	// ManagerServiceURL refers to the URL of the Manager service
	ManagerServiceURL = os.Getenv("MANAGER_SERVICE_URL")

	ErrCustomerNotfound    = errors.New("customer not found")
	ErrForbidden           = errors.New("Forbidden")
	ErrUnauthorized        = errors.New("Unauthorized")
	ErrRoleBindingNotFound = errors.New("role-binding not found")
)

// ManagerClient interface defines Manager APIs
type ManagerClient interface {
	// Query Manager service to get all buckets of the tenant
	ListBuckets() ([]model.Bucket, error)

	// Create a bucket object at Manager service
	CreateBucket(bucketName string) (*model.Bucket, error)

	DeleteBucket(bucketName string) error

	UpdateACLs(bucketName string, acls authz.ACLList) (*model.Bucket, error)

	ListRoleBindings() ([]rbac.RoleBinding, error)

	CreateRoleBinding(roleBinding rbac.RoleBinding) (*rbac.RoleBinding, error)

	DeleteRoleBinding(name string) error

	ListRoles() ([]rbac.Role, error)
}
type managerClient struct {
	serverURL    string
	userId       string
	userGroups   []string
	customerName string
}

func New(serverURL string, customerName, userId string, userGroups []string) ManagerClient {
	return managerClient{
		serverURL:    serverURL,
		userId:       userId,
		userGroups:   userGroups,
		customerName: customerName,
	}
}

func (client managerClient) ListBuckets() ([]model.Bucket, error) {
	var payload []model.Bucket
	err := client.getManagerAPI("/buckets", &payload)
	return payload, err
}

func (client managerClient) CreateBucket(bucketName string) (*model.Bucket, error) {
	type Payload struct {
		Name string `json:"name"`
	}

	resp, err := resty.R().
		SetHeader(shared.UserIdKey, client.userId).
		SetHeader(shared.UserGroupKey, strings.Join(client.userGroups, ",")).
		SetBody(Payload{
			Name: bucketName,
		}).
		Post(client.serverURL + "/customers/" + client.customerName + "/buckets")
	if err != nil {
		glog.Errorf("Called create bucket API failed: %v", err)
		return nil, err
	}

	body := resp.Body()
	code := resp.StatusCode()
	switch code {
	case http.StatusOK:
		return model.NewBucketFromBytes(body)
	case http.StatusConflict:
		errorCode := string(body)
		glog.Errorf("Cannot create bucket for customer '%s': ErrorCode=%s", client.customerName, errorCode)
		return nil, errors.New(errorCode)
		// Add more specific error cases here
	}
	glog.Errorf("Cannot create bucket for customer '%s'", client.customerName)
	if glog.V(2) {
		glog.Infof("Dumped body: %v", string(body))
	}
	return nil, errors.New("Cannot create buckets for customer")
}

func (client managerClient) DeleteBucket(bucketName string) error {
	type Payload struct {
		Name string `json:"name"`
	}

	_, err := resty.R().
		SetHeader(shared.UserIdKey, client.userId).
		SetHeader(shared.UserGroupKey, strings.Join(client.userGroups, ",")).
		SetBody(Payload{
			Name: bucketName,
		}).Delete(client.serverURL + "/customers/" + client.customerName + "/buckets")
	if err != nil {
		glog.Errorf("Called delete bucket API failed: %v", err)
		return err
	}
	glog.V(2).Infof("Deleted bucket object '%s' of customer '%s'", bucketName, client.customerName)
	return nil
}

func (client managerClient) UpdateACLs(bucketName string, acls authz.ACLList) (*model.Bucket, error) {
	resp, err := resty.R().
		SetHeader(shared.UserIdKey, client.userId).
		SetHeader(shared.UserGroupKey, strings.Join(client.userGroups, ",")).
		SetBody(acls).
		Put(client.serverURL + "/customers/" + client.customerName + "/buckets/" + bucketName + "/acl")
	if err != nil {
		glog.Errorf("Called delete bucket API failed: %v", err)
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		glog.Errorf("Cannot update ACLs for bucket='%s',customer=%s", bucketName, client.customerName)
		if glog.V(2) {
			body := resp.Body()
			glog.Infof("Dumped body: %v", string(body))
		}
		return nil, errors.New("Cannot update ACLs")
	}

	// fetch new bucket object

	var bucket model.Bucket
	err = client.getManagerAPI("/buckets/"+bucketName, &bucket)
	if err != nil {
		return nil, err
	}
	return &bucket, nil
}

func (client managerClient) ListRoleBindings() ([]rbac.RoleBinding, error) {
	var payload []rbac.RoleBinding
	err := client.getManagerAPI("/role-bindings", &payload)
	return payload, err
}

func (client managerClient) ListRoles() ([]rbac.Role, error) {
	var payload []rbac.Role
	err := client.getManagerAPI("/roles", &payload)
	return payload, err
}

func (client managerClient) CreateRoleBinding(roleBinding rbac.RoleBinding) (*rbac.RoleBinding, error) {
	resp, err := resty.R().
		SetHeader(shared.UserIdKey, client.userId).
		SetHeader(shared.UserGroupKey, strings.Join(client.userGroups, ",")).
		SetBody(roleBinding).
		Post(client.serverURL + "/customers/" + client.customerName + "/role-bindings")
	if err != nil {
		glog.Errorf("Called create role-binding API failed: %v", err)
		return nil, err
	}

	body := resp.Body()
	if resp.StatusCode() != http.StatusOK {
		glog.Errorf("Cannot create role-binding for customer '%s'", client.customerName)
		if glog.V(2) {
			glog.Infof("Dumped body: %v", string(body))
		}
		return nil, errors.New("Cannot create role-binding for customer")
	}
	return rbac.NewRoleBindingFromBytes(body)
}

func (client managerClient) DeleteRoleBinding(name string) error {
	type Payload struct {
		Name string `json:"name"`
	}

	resp, err := resty.R().
		SetHeader(shared.UserIdKey, client.userId).
		SetHeader(shared.UserGroupKey, strings.Join(client.userGroups, ",")).
		SetBody(Payload{
			Name: name,
		}).Delete(client.serverURL + "/customers/" + client.customerName + "/role-bindings")
	if err != nil {
		glog.Errorf("Called delete role-binding API failed: %v", err)
		return err
	}
	code := resp.StatusCode()
	switch code {
	case http.StatusNotFound:
		return ErrRoleBindingNotFound
	case http.StatusOK:
		glog.V(2).Infof("Deleted role-binding object '%s' of customer '%s'", name, client.customerName)
		return nil
	}
	glog.Errorf("Failed to delete role-binding '%s': status='%d', body='%v'", name, code, string(resp.Body()))
	return pion_clients.ErrInternalError
}

func (client managerClient) getManagerAPI(apiSuffix string, payload interface{}) (err error) {
	resp, err := resty.R().
		SetHeader(shared.UserIdKey, client.userId).
		SetHeader(shared.UserGroupKey, strings.Join(client.userGroups, ",")).
		Get(client.serverURL + "/customers/" + client.customerName + apiSuffix)
	if err != nil {
		return err
	}

	if resp.StatusCode() == http.StatusNotFound {
		return ErrCustomerNotfound
	}

	body := resp.Body()
	switch resp.StatusCode() {
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusForbidden:
		return ErrForbidden
	case http.StatusOK:
		err = json.Unmarshal(body, &payload)
		if err != nil {
			glog.Errorf("Failed to parse response from Manager API '%s': %v", apiSuffix, err)
			return err
		}
		return nil
	}

	glog.Errorf("Cannot invoke Manager API '%s' with customer '%s'", apiSuffix, client.customerName)
	if glog.V(2) {
		glog.Infof("Dumped body: %v", string(body))
	}
	return errors.New("Cannot invoke Manager API")
}
