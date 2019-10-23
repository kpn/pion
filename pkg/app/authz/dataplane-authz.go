package authz

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/authz"
	"github.com/kpn/pion/pkg/pion/authz/handlers"
	"github.com/kpn/pion/pkg/pion/authz/request"
	"github.com/kpn/pion/pkg/pion/shared"
	"github.com/kpn/pion/pkg/pion/shared/rbac"
	"github.com/labstack/echo"
)

// DataPlaneAuthorize implements access control mechanism for data plane requests (i.e. requests from S3 clients)
func (app *App) DataPlaneAuthorize(c echo.Context) error {
	rbacRequest, aclRequest, err := request.Build(c.Request())
	if err != nil {
		glog.Errorf("building request failed: %v", err)
		return c.NoContent(http.StatusForbidden)
	}
	glog.V(3).Infof("RBAC request: %+v", rbacRequest)
	glog.V(3).Infof("ACL request: %+v", aclRequest)

	rbachdlr, err := handlers.NewRBACHandler(shared.DefaultEtcdAddress, aclRequest.Customer)
	if err != nil {
		glog.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	rbacDecision := rbachdlr.Authorize(rbacRequest)

	if creatingBucket(rbacRequest) {
		glog.V(2).Infof("Creating bucket '%s'", aclRequest.Target)
		if rbacDecision == rbac.PermitDecision {
			return c.NoContent(http.StatusOK)
		} else {
			glog.Info("Access denied. Cannot create bucket")
			return c.String(http.StatusForbidden, "Access denied. Cannot create bucket")
		}
	}

	// Accessing existing buckets, need to verify bucket customer
	if !app.verifyBucketCustomer(*aclRequest) {
		glog.Warningf("Resource '%s' does not belong to customer '%s'", aclRequest.Target, aclRequest.Customer)
		return c.String(http.StatusForbidden, "bucket either does not exist or belong to the customer")
	}
	glog.V(2).Infof("Resource '%s' is owned by '%s'", aclRequest.Target, aclRequest.Customer)

	aclDecision := app.engine.Evaluate(*aclRequest)

	// either decision is accepted
	if (rbacDecision == rbac.PermitDecision) || (aclDecision == authz.DecisionPermit) {
		glog.V(2).Infof("Granted permission, rbac-decision='%s', acl-decision='%s'", rbacDecision, aclDecision)
		return c.NoContent(http.StatusOK)
	}

	glog.V(2).Infof("Forbidden request")
	if glog.V(3) {
		glog.Infof("rbac-request: %+v", rbacRequest)
		glog.Infof("acl-request: %+v", aclRequest)
	}

	return c.NoContent(http.StatusForbidden)
}

func creatingBucket(rbacRequest *rbac.Request) bool {
	return rbacRequest.Action == rbac.Create && rbacRequest.Target == rbac.BucketResource
}

// verifyBucketCustomer returns true if the target bucket belongs to the customer defined in the request
func (app *App) verifyBucketCustomer(r authz.Request) bool {
	bkt := app.bucketsWatcher.GetBucket(r.Target)
	if bkt == nil {
		glog.Warningf("Bucket object '%s' not found", r.Target)
		return false
	}
	return bkt.OwnedBy == r.Customer
}
