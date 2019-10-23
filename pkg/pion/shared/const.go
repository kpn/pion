package shared

const (
	UserIdKey        = "X-User-Id"     // UserIdKey is the http header having UserID
	UserGroupKey     = "X-User-Groups" // UserGroupKey is the http header having UserGroups
	CustomerKey      = "X-Customer"    // CustomerKey is the http header having customerID
	ActionKey        = "X-Action"      // ActionKey is the http header having actionID
	ResourceKey      = "X-Resource"    // ResourceKey is the http header having resourceID
	DefaultKeyPrefix = "/pion/oss"     // DefaultKeyPrefix set the common Etcd key prefix for all Pion configs
	SystemKey        = "_system"
)
