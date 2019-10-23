package api

import (
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion"
	"github.com/kpn/pion/pkg/pion-clients"
	"github.com/kpn/pion/pkg/pion-clients/sts-client"
	"github.com/kpn/pion/pkg/pion/ui/multi_tenant"
	"github.com/kpn/pion/pkg/pion/ui/session"
	"github.com/kpn/pion/pkg/sts/handlers"
	"github.com/labstack/echo"
)

type AccessKey struct {
	Id        string    `json:"accessKeyId,omitempty"`
	SecretKey string    `json:"secretKeyId,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	LastUsed  time.Time `json:"lastUsed,omitempty"`
	Status    string    `json:"status,omitempty"`
}

func ListAccessKeys(c echo.Context) error {
	userInfo, err := session.GetUserInfo(c)
	if err != nil {
		glog.Errorf("Failed to get userInfo from session: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	glog.V(2).Infof("Listing access keys of user %s", userInfo.Username)

	stsClient := sts_client.New(userInfo.Customer)
	keyInfos, err := stsClient.ListAccessKeys(userInfo.Username)
	if err != nil {
		glog.Errorf("Failed to query access keys of user '%s': %v", userInfo.Username, err)
		return nil
	}

	var accessKeys []AccessKey = nil
	for _, ki := range keyInfos {
		accessKeys = append(accessKeys, AccessKey{
			Id:        ki.AccessKey,
			CreatedAt: ki.CreatedAt,
			Status:    "Active",                       // faked
			LastUsed:  time.Now().Add(-4 * time.Hour), // faked
		})
	}
	return c.JSON(http.StatusOK, accessKeys)
}

func DeleteAccessKey(c echo.Context) error {
	userInfo, err := session.GetUserInfo(c)
	if err != nil {
		glog.Error(err)
		return handlers.Response(c, http.StatusInternalServerError, err.Error())
	}

	accessKey := c.Param("id")
	if accessKey == "" {
		glog.Errorf("Empty accessKeyId")
		return handlers.Response(c, http.StatusBadRequest, "Empty accessKeyId")
	}

	stsClient := sts_client.New(userInfo.Customer)
	err = stsClient.DeleteAccessKey(accessKey)
	if err == pion_clients.ErrAccessKeyNotFound {
		return handlers.Response(c, http.StatusNotFound, "key not found")
	} else if err != nil {
		return handlers.Response(c, http.StatusInternalServerError, "Delete access key failed")
	}

	return handlers.Response(c, http.StatusOK, "Deleted key successfully")
}

func CreateAccessKey(c echo.Context) error {
	userInfo, err := session.GetUserInfo(c)
	if err != nil {
		glog.Error(err)
		return handlers.Response(c, http.StatusInternalServerError, err.Error())
	}
	glog.V(1).Infof("Creating access key for user '%s'", userInfo.Username)
	glog.V(2).Infof("UserInfo=%+v", userInfo)

	stsClient := sts_client.New(userInfo.Customer)

	attributes := make(map[string]interface{})
	attributes[multi_tenant.UserAttributeGroups] = userInfo.UserGroups
	attributes[multi_tenant.UserAttributeCustomer] = userInfo.Customer

	resp, err := stsClient.CreateAccessKey(userInfo.Username, pion.AppConfig().TokenLifetime, attributes)
	if err != nil {
		glog.Error(err)
		return handlers.Response(c, http.StatusInternalServerError, "Failed to create access key")
	}

	glog.V(2).Infof("Created access key for '%s'", userInfo.Username)

	return c.JSON(http.StatusOK, AccessKey{
		Id:        resp.AccessKey,
		SecretKey: resp.SecretKey,
		CreatedAt: resp.CreatedAt,
	})
}
