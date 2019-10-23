package policy

import (
	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion/authz"
	"github.com/kpn/pion/pkg/pion/authz/watcher"
)

// Engine defines API for the ACL evaluation engine
type Engine interface {
	Evaluate(r authz.Request) authz.DecisionType
}

type defaultEngine struct {
	bucketsWatcher watcher.BucketsWatcher
}

// NewEngine creates a policy engine backed by latest ACLs in buckets at the given path
func NewEngine(bw watcher.BucketsWatcher) (Engine, error) {
	return &defaultEngine{
		bucketsWatcher: bw,
	}, nil
}

// Evaluate goes through ACLs of the bucket and check if request is permit. Evaluation is PermitOverride, Deny by default
func (e defaultEngine) Evaluate(r authz.Request) authz.DecisionType {
	acls := e.getBucketACLs(r.Target)
	if acls == nil {
		// if no ACLs is attached to the bucket, the bucket is deny
		glog.V(2).Infof("ACLs for bucket '%s' not found", r.Target)
		return authz.DecisionDeny
	}
	glog.V(2).Infof("Found ACLs of the bucket '%s': %+v", r.Target, acls)
	// evaluate list of ACLs, Permit overridden
	for _, acl := range acls {
		decision := acl.Evaluate(r)
		if decision == authz.DecisionPermit {
			return authz.DecisionPermit
		}
	}
	return authz.DecisionDeny
}

func (e defaultEngine) getBucketACLs(bucketName string) authz.ACLList {
	bkt := e.bucketsWatcher.GetBucket(bucketName)
	if bkt == nil {
		glog.Warningf("Bucket '%s' not found", bucketName)
		return nil
	}
	return bkt.ACLs
}
