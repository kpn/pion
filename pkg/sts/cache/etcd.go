package cache

import (
	"os"
	"time"

	"github.com/golang/glog"
	"go.etcd.io/etcd/clientv3"
)

var DefaultEtcdAddress = os.Getenv("ETCD_ADDRESS")

// NewEtcdClient simplifies creating a Etcd v3 client
func NewEtcdClient(address string) (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		Endpoints:   []string{address},
		DialTimeout: 3 * time.Second,
	})
}

// SilentClose closes the etcd-v3 client without error popup. It's used in defer calls
func SilentClose(c *clientv3.Client) {
	err := c.Close()
	if err != nil {
		glog.Error(err.Error())
	}
}
