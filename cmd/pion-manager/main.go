package main

import (
	"flag"
	"os"

	mgr "github.com/kpn/pion/pkg/app/manager"
	"github.com/spf13/pflag"
)

var DefaultEtcdAddress = os.Getenv("ETCD_ADDRESS")

func init() {
	flags := pflag.NewFlagSet("", pflag.ExitOnError)
	flag.Set("logtostderr", "true")
	flags.AddGoFlagSet(flag.CommandLine)
	flags.Parse(os.Args)
}

func main() {

	mgr.NewApp().Start(DefaultEtcdAddress)
}
