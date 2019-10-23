package middleware

import (
	"context"
	"net/http"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion-clients"
	"github.com/kpn/pion/pkg/pion-clients/sts-client"
	"github.com/kpn/pion/pkg/pion/proxy"
	"github.com/kpn/pion/pkg/pion/proxy/debug"
	"github.com/kpn/pion/pkg/pion/shared"
	"github.com/labstack/echo"
	minio "github.com/minio/minio/cmd"
)

// Authenticate middleware does the request authentication and signature verification. Authenticated request is then
// inserted some authentication headers, e.g. X-User-Id, X-User-Groups
func Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if isPublic(c) {
			glog.V(2).Infof("Request to a public URL: '%s'", c.Request().URL.Path)
			return next(c)
		}

		errCode := verifySignatureAndSetAuthenticatedUser(c)

		if errCode == minio.ErrNone {
			return next(c)
		}

		glog.Warningf("Authentication failed. ErrCode=%v", errCode)
		// TODO return proper AWS signature v4 response
		switch errCode {
		case minio.ErrNoSuchKey:
			return c.NoContent(http.StatusBadRequest)
		case minio.ErrAccessDenied:
			return c.NoContent(http.StatusForbidden)
		case minio.ErrInternalError:
			return c.NoContent(http.StatusInternalServerError)
		default:
			return c.NoContent(http.StatusUnauthorized)
		}
	}
}

func verifySignatureAndSetAuthenticatedUser(c echo.Context) minio.APIErrorCode {
	r := c.Request()
	glog.V(3).Infof("Request before signing:\n%s", debug.FormatRequest(r, false))

	accessKey, s3Err := minio.AccessKeyFromRequest(r, proxy.DefaultRegion)
	if s3Err != minio.ErrNone {
		glog.Warningf("Invalid request, no access key found: ErrCode=%v", s3Err)
		return s3Err
	}
	// TODO Filter invalid accessKey values
	glog.V(2).Infof("Verifying request with accessKey='%v'", accessKey)

	keyQuerier := sts_client.NewAccessKeyQuerier()
	stsResp, err := keyQuerier.Query(accessKey)
	if err == pion_clients.ErrAccessKeyNotFound {
		glog.Warningf("Secret key of access key '%s' is empty", accessKey)
		return minio.ErrNoSuchKey
	}
	if err != nil {
		glog.Errorf("Cannot read secret key of the given access key '%s': %v", accessKey, err)
		return minio.ErrInternalError
	}

	verifier, err := minio.NewAWSV4Verifier(accessKey, stsResp.SecretKey, proxy.DefaultRegion)
	if err != nil {
		glog.Errorf("Failed to create AWS signature-v4 verifier: %v", err)
		return minio.ErrInternalError
	}

	ctx := context.Background()
	s3Err = verifier.IsReqAuthenticated(ctx, r, verifier.GetRegion())
	if s3Err != minio.ErrNone {
		glog.Warningf("Request error with s3Err=%d", s3Err)
		return s3Err
	}
	glog.V(2).Info("valid signature")

	groups := stsResp.GetUserGroups()
	customer := stsResp.GetUserCustomer()
	// Insert authentication info to the context
	c.Set(shared.UserIdKey, stsResp.UserId)
	c.Set(shared.UserGroupKey, groups)
	c.Set(shared.CustomerKey, customer)

	glog.V(2).Infof("Set to request context: X-User-Id='%s', X-User-Groups='%s', X-Customer='%s'", stsResp.UserId, groups, customer)
	return minio.ErrNone
}
