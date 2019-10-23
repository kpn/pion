package middleware

import (
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion-clients/manager-client"
	"github.com/kpn/pion/pkg/pion/proxy/s3responses"
	"github.com/kpn/pion/pkg/pion/proxy/user-context"
	"github.com/kpn/pion/pkg/pion/shared/s3"
	"github.com/labstack/echo"
)

const (
	populateS3BucketAction = "populate_s3_bucket_action"
	populateS3BucketTarget = "populate_s3_bucket_target"
)

// PopulateBucket calls Manager APIs to create/delete bucket when the request context having flag. This middleware runs
// BEFORE the API handler.
func PopulateBucket(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		s3Action, ok := c.Get(populateS3BucketAction).(s3.ActionType)
		if !ok || s3Action == "" {
			return next(c)
		}
		glog.V(2).Info("Bucket post processing, action=", s3Action)
		bucketName, ok := c.Get(populateS3BucketTarget).(string)
		if !ok {
			glog.Error("Failed to get bucket name from context")
			return c.XML(http.StatusInternalServerError, s3responses.InternalErrorResponse)
		}

		glog.V(2).Infof("Bucket post processing: action=%s, bucket=%s", s3Action, bucketName)
		userId, userGroups, customerName, err := user_context.GetData(c)
		if err != nil {
			glog.Errorf("Get user context failed: %v", err)
			return c.XML(http.StatusInternalServerError, s3responses.InternalErrorResponse)
		}
		mc := manager_client.New(manager_client.ManagerServiceURL, customerName, userId, userGroups)

		switch s3Action {
		case s3.CreateBucket:
			bkt, err := mc.CreateBucket(bucketName)
			httpCode, apiErrResponse := handleS3APIError(err, bucketName)
			if httpCode != http.StatusOK {
				logWithLevel(httpCode, "creating bucket failed: %v", err)
				return c.XML(httpCode, apiErrResponse)
			}
			glog.V(2).Infof("Created bucket '%s': %v", bkt.Name, bkt)
		case s3.DeleteBucket:
			err := mc.DeleteBucket(bucketName)
			httpCode, apiErrResponse := handleS3APIError(err, bucketName)
			if httpCode != http.StatusOK {
				logWithLevel(httpCode, "delete bucket failed: %v", err)
				return c.XML(httpCode, apiErrResponse)
			}
			glog.V(2).Infof("Deleted bucket '%s'", bucketName)
		default:
			glog.Errorf("Populating s3 bucket error. Unknown action '%s", s3Action)
			return c.XML(http.StatusInternalServerError, s3responses.InternalErrorResponse)
		}
		return next(c)
	}
}

// logWithLevel logs message with level based on http-code: 5xx - error, 4xx warning, otherwise info
func logWithLevel(httpCode int, format string, args ...interface{}) {
	logMessage := fmt.Sprintf(format, args...)
	if httpCode >= 500 {
		glog.ErrorDepth(1, logMessage)
	} else if httpCode > 400 {
		glog.WarningDepth(1, logMessage)
	} else {
		glog.InfoDepth(1, logMessage)
	}
}

func handleS3APIError(err error, bucketName string) (int, *s3responses.Error) {
	if err == nil {
		return http.StatusOK, nil
	}

	errCode := err.Error()
	httpStatusCode := s3responses.MapHttpCodes[errCode]
	if httpStatusCode == 0 {
		httpStatusCode = http.StatusInternalServerError
	}
	return httpStatusCode, s3responses.NewError(errCode, bucketName)
}
