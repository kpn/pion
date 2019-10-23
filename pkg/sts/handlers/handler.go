package handlers

import (
	"net/http"
	"time"

	"github.com/golang/glog"
	store "github.com/kpn/pion/pkg/sts/pion-store"
	"github.com/kpn/pion/pkg/sts/secure_rand"
	"github.com/labstack/echo"
)

const (
	accessKeyLength = 16
	secretKeyLength = 32
)

var (
	maxTokenTTL     = 6 * 30 * 24 * time.Hour // 6 months
	defaultTokenTTL = 7 * 24 * time.Hour      // 1 week
)

// AccessKeyHandler contains methods to manage access keys
type AccessKeyHandler interface {
	// Create handles API creating a new access key
	Create(c echo.Context) error

	// List handles API listing all access keys of a user
	List(c echo.Context) error

	// Remove handles API removing a given access key
	Revoke(c echo.Context) error
}

type accessKeyHandler struct {
	keystore store.KeyStore
}

// NewAccessKeyHandler returns an object implemented AccessKeyHandler
func NewAccessKeyHandler(ks store.KeyStore) AccessKeyHandler {
	return &accessKeyHandler{
		keystore: ks,
	}
}

func (h accessKeyHandler) Create(c echo.Context) error {
	type Payload struct {
		UserId     string                 `json:"userId,omitempty"`
		Lifetime   string                 `json:"lifetime,omitempty"`
		Attributes map[string]interface{} `json:"attributes,omitempty"`
	}

	var p Payload
	err := c.Bind(&p)
	if err != nil {
		glog.Warningf("Payload is invalid or not found: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Payload is invalid or not found")
	}
	if p.UserId == "" {
		glog.V(1).Info("empty userId attribute in payload")
		return echo.NewHTTPError(http.StatusBadRequest, "empty userId attribute in payload")
	}

	lifetime, err := time.ParseDuration(p.Lifetime)
	if err != nil {
		glog.V(1).Infof("Cannot generate access key: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Cannot generate access key, invalid duration")
	}
	if lifetime == 0 {
		glog.V(1).Info("Token lifetime is 0, use default value")
		lifetime = defaultTokenTTL
	}
	if lifetime > maxTokenTTL {
		glog.V(1).Infof("Request token lifetime is too long: %v", lifetime)
		return echo.NewHTTPError(http.StatusBadRequest, "Cannot generate access key, duration is too long")
	}

	accessKey, err := secure_rand.SecureRandomString(accessKeyLength)
	if err != nil {
		glog.Errorf("Cannot generate access key: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Cannot generate access key")
	}
	accessKey = "pion-" + accessKey

	secretKey, err := secure_rand.SecureRandomString(secretKeyLength)
	if err != nil {
		glog.Errorf("Cannot generate secret key: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Cannot generate secret key")
	}

	createdAt := time.Now().UTC()
	// write to Etcd
	glog.V(2).Infof("Saved access key='%v' of user '%s'", accessKey, p.UserId)

	h.keystore.SaveKey(p.UserId, accessKey, lifetime, store.SecretData{
		UserId:     p.UserId,
		Attributes: p.Attributes,
		SecretKey:  secretKey,
		CreatedAt:  createdAt,
	})
	if err != nil {
		glog.Errorf("Cannot save Pion access key: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Cannot save Pion access key")
	}

	type Response struct {
		AccessKey string    `json:"accessKey"`
		SecretKey string    `json:"secretKey"`
		CreatedAt time.Time `json:"createdAt"`
	}
	return c.JSON(http.StatusOK, &Response{
		AccessKey: accessKey,
		SecretKey: secretKey,
		CreatedAt: createdAt,
	})
}

func (h accessKeyHandler) List(c echo.Context) error {
	type KeyInfo struct {
		AccessKey string    `json:"accessKey"`
		CreatedAt time.Time `json:"createdAt,omitempty"`
	}

	username := c.Param("username")
	if username == "" {
		return c.JSON(http.StatusBadRequest, "Missing username")
	}
	glog.V(2).Infof("Listing access keys of user '%s'", username)

	keys, err := h.keystore.ListAccessKeys(username)
	if err != nil {
		glog.Errorf("Cannot list access keys of user '%s': %v", username, err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Cannot list access keys")
	}

	var accessKeyInfos []KeyInfo
	for _, key := range keys {
		data, err := h.keystore.Query(key)
		if err != nil {
			glog.Errorf("Failed to query data of key '%s': %v", key, err)
			continue
		}
		accessKeyInfos = append(accessKeyInfos, KeyInfo{
			AccessKey: key,
			CreatedAt: data.CreatedAt,
		})
	}

	return c.JSON(http.StatusOK, accessKeyInfos)
}

func (h accessKeyHandler) Revoke(c echo.Context) error {
	type Payload struct {
		AccessKey string `json:"accessKey"`
	}

	var p Payload
	err := c.Bind(&p)
	if err != nil {
		glog.Errorf("Failed to parse payload: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse payload")
	}

	if p.AccessKey == "" {
		glog.Error("Empty access key")
		return echo.NewHTTPError(http.StatusBadRequest, "Empty access key")
	}

	err = h.keystore.DeleteAccessKey(p.AccessKey)
	if err != nil {
		glog.Errorf("Failed to delete access key: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete access key")
	}
	return c.NoContent(http.StatusOK)
}
