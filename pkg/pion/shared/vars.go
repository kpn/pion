package shared

import "os"

var DefaultEtcdAddress = os.Getenv("ETCD_ADDRESS")
