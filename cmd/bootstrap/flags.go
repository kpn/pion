package main

import (
	"flag"
	"os"

	"github.com/spf13/pflag"
)

var (
	adminUser string
)

func parseFlags() {
	var (
		flags = pflag.NewFlagSet("", pflag.ExitOnError)
		user  = flags.String("admin-user", "admin", "Admin user of all customers")
	)

	flag.Set("logtostderr", "true")
	flags.AddGoFlagSet(flag.CommandLine)
	flags.Parse(os.Args)

	adminUser = *user
}
