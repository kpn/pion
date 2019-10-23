package main

import (
	"encoding/gob"

	"github.com/golang/glog"
	app "github.com/kpn/pion/pkg/app/ui"
	"github.com/kpn/pion/pkg/pion"
	"github.com/kpn/pion/pkg/pion/ui/multi_tenant"
)

func init() {
	gob.Register(multi_tenant.UserInfo{})
}

func main() {
	conf, err := parseFlags()
	if err != nil {
		glog.Fatal(err.Error())
	}
	pion.SetAppConfig(conf)

	app.NewApp().Start()
}
