package api

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion-clients/manager-client"
	"github.com/kpn/pion/pkg/pion/authz"
	"github.com/kpn/pion/pkg/pion/ui/util"
	"github.com/labstack/echo"
)

// ListBuckets handles requests listing all buckets in the UI
func ListBuckets(c echo.Context) error {
	mc, err := createManagerClient(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	buckets, err := mc.ListBuckets()
	switch err {
	case manager_client.ErrForbidden:
		return echo.ErrForbidden
	case manager_client.ErrUnauthorized:
		return echo.ErrUnauthorized
	case manager_client.ErrCustomerNotfound:
		return echo.NewHTTPError(http.StatusBadRequest, "customer not found")
	}

	if err != nil {
		glog.Errorf("Failed to list buckets: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list buckets")
	}
	glog.V(2).Infof("Listing buckets: %v", buckets)
	return c.JSON(http.StatusOK, buckets)
}

// UpdateBucketACLs handles requests updating bucket's ACLs from UI
func UpdateBucketACLs(c echo.Context) error {
	bucketName := c.Param("id")
	if bucketName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "empty bucket name parameter")
	}

	var payload authz.ACLList
	err := c.Bind(&payload)
	if err != nil {
		glog.Warningf("Updating bucet ACLs: parsing body failed: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid payload")
	}

	mc, err := createManagerClient(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	bucket, err := mc.UpdateACLs(bucketName, payload)
	if err != nil {
		glog.Errorf("Updating bucket ACLs failed: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "updating ACLs failed")
	}
	return c.JSON(http.StatusOK, bucket)
}

// AddBucket handles requests creating new bucket from UI. If error, it returns the error response with either error codes
// BucketAlreadyExists or InternalError
func AddBucket(c echo.Context) error {
	const (
		ErrCodeInternalError = "InternalError"
	)

	type Payload struct {
		Name string `json:"name"`
	}
	type ErrorResponse struct {
		ErrorCode string `json:"error_code"`
	}
	var p Payload
	err := c.Bind(&p)
	if err != nil {
		glog.Warningf("Add bucket - parsing body failed: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid payload")
	}
	if err = util.ValidateBucketName(p.Name); err != nil {
		glog.V(2).Infof("Bucket name '%s' is invalid: %v", p.Name, err)
		return c.JSON(http.StatusBadRequest, "Invalid name")
	}

	// create bucket in the upstream server
	err = util.CreateS3Bucket(p.Name)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			awsErrcode := aerr.Code()
			switch awsErrcode {
			case s3.ErrCodeBucketAlreadyExists:
				fallthrough
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				glog.V(1).Infof("Failed to create bucket in the upstream: %v", awsErrcode)
				// error code 'BucketAlreadyOwnedByYou' is incorrect as all tenants are sharing the same Minio server
				return c.JSON(http.StatusConflict, ErrorResponse{s3.ErrCodeBucketAlreadyExists})
			default:
				glog.Errorf("Failed to create bucket: %v", err)
				return c.JSON(http.StatusInternalServerError, ErrorResponse{awsErrcode})
			}
		} else {
			glog.Errorf("Failed to create bucket in the upstream: %v", err)
			return c.JSON(http.StatusInternalServerError, ErrorResponse{ErrCodeInternalError})
		}
	}

	// create bucket meta-data in Manager
	mc, err := createManagerClient(c)
	if err != nil {
		glog.Errorf("Failed to create manager client%v", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{ErrCodeInternalError})
	}

	bkt, err := mc.CreateBucket(p.Name)
	if err != nil {
		var statusCode int
		var errorCode = err.Error()
		if errorCode == s3.ErrCodeBucketAlreadyExists {
			statusCode = http.StatusConflict
		} else {
			statusCode = http.StatusInternalServerError
		}
		glog.V(1).Infof("Failed to create bucket in the Manager: %v", err)
		return c.JSON(statusCode, ErrorResponse{errorCode})
	}

	glog.V(2).Infof("Created bucket '%s': %v", bkt.Name, bkt)
	return c.JSON(http.StatusOK, bkt)
}
