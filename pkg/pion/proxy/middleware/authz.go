package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion-clients/manager-client"
	"github.com/kpn/pion/pkg/pion/proxy/s3responses"
	"github.com/kpn/pion/pkg/pion/proxy/user-context"
	"github.com/kpn/pion/pkg/pion/shared"
	"github.com/kpn/pion/pkg/pion/shared/s3"
	"github.com/labstack/echo"
	"gopkg.in/resty.v1"
)

var (
	authzServiceURL = os.Getenv("AUTHZ_SERVICE_URL")
)

// Authorize is the middleware responsible for authorization interceptor
func Authorize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if isPublic(c) {
			return next(c)
		}

		glog.Info("Evaluating Pion authorization policies")
		r := c.Request()

		resourceId := r.URL.Path
		bucketName, keyPath, err := s3.GetResources(resourceId)
		if err != nil {
			glog.Errorf("Error parsing bucket name: %v", err)
			return c.NoContent(http.StatusForbidden)
		}

		s3Action, err := s3.ToS3Action(r.Method, keyPath)
		if err != nil {
			glog.Errorf("Authorization failed: %v", err)
			return c.NoContent(http.StatusForbidden)
		}
		glog.V(3).Infof("S3 Action=%s", s3Action)
		// FIXME Only Bucket APIs and Object APIs containing bucket name in the URI

		if (bucketName == "" || bucketName == "probe-bucket-sign") && s3Action == s3.ListBucket {
			// TODO Add authz on LIST permission
			return listCustomerBuckets(c)
		}

		userId, userGroups, customerName, err := user_context.GetData(c)
		if err != nil {
			glog.Error("failed to get authenticated user info: ", err)
			return c.NoContent(http.StatusUnauthorized)
		}

		glog.V(2).Infof("Authorizing: customer %s - user '%s' - groups '%v'", customerName, userId, userGroups)

		// TODO cache authz responses during a fixed period, e.g. 5 minutes
		resp, err := resty.R().
			SetHeader(shared.UserIdKey, userId).
			SetHeader(shared.UserGroupKey, strings.Join(userGroups, ",")).
			SetHeader(shared.CustomerKey, customerName).
			SetHeader(shared.ActionKey, string(s3Action)).
			SetHeader(shared.ResourceKey, resourceId).
			Get(authzServiceURL + "/data-authorize")
		if err != nil {
			glog.Errorf("Cannot call authorization service: %v", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		if resp.StatusCode() != http.StatusOK {
			glog.Warningf("Unauthorized decision from authz-service: %d", resp.StatusCode())
			return c.NoContent(http.StatusForbidden)
		}

		if s3Action == s3.CreateBucket || s3Action == s3.DeleteBucket {
			// Set flag in context to create/delete Pion bucket objects (via Manager APIs)
			c.Set(populateS3BucketAction, s3Action)
			c.Set(populateS3BucketTarget, bucketName)
		}
		if glog.V(4) {
			glog.Infof("Granted permission for '%+v'", r)
		} else {
			glog.Info("Granted permission")
		}
		return next(c)
	}
}

// listCustomerBuckets calls Manager service to list buckets of the customer and returns a S3 XML response
func listCustomerBuckets(context echo.Context) error {
	userId, userGroups, customerName, err := user_context.GetData(context)
	if err != nil {
		glog.Error("Get user context failed:", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "get user session failed")
	}

	mc := manager_client.New(manager_client.ManagerServiceURL, customerName, userId, userGroups)
	buckets, err := mc.ListBuckets()
	switch err {
	case manager_client.ErrUnauthorized:
		return echo.ErrUnauthorized
	case manager_client.ErrForbidden:
		return echo.ErrForbidden
	case manager_client.ErrCustomerNotfound:
		return echo.NewHTTPError(http.StatusNotFound, "customer not found")
	case nil: // ignore
	default: // not nil
		glog.Errorf("List bucket failed: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "list bucket failed")
	}

	resp, err := s3responses.NewListBucketsResult(buckets)
	if err != nil {
		glog.Errorf("Create s3 response failed: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "create s3 response failed")
	}
	return context.XML(http.StatusOK, resp)
}
