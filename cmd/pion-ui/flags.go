package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

var (
	maxTokenTTL     = 6 * 30 * 24 * time.Hour // 6 months
	defaultTokenTTL = 7 * 24 * time.Hour      // 1 week
)

func parseFlags() (*pion.Config, error) {
	var (
		flags = pflag.NewFlagSet("", pflag.ExitOnError)

		ldapHost         = flags.String("ldap-host", "", "LDAP host used for user authentication")
		ldapPort         = flags.Int("ldap-port", 636, "LDAP port used for user authentication")
		ldapType         = flags.String("ldap-type", "ldaps", "LDAP type, either 'ldaps' or 'ldap'")
		ldapBindDN       = flags.String("ldap-binddn", "cn=canh,dc=example,dc=org", "Distinguished name of the binding user")
		ldapBindPassword = flags.String("ldap-bindpassword", "p4ssw0rd", "Password of the binding DN")
		ldapUserBase     = flags.String("ldap-userbase", "ou=People,dc=example,dc=org", "User base")
		ldapGroupBase    = flags.String("ldap-groupbase", "ou=Groups,dc=example,dc=org", "Group base")
		ldapUserClass    = flags.String("ldap-userclass", "inetOrgPerson", "User objectClass")
		ldapGroupClass   = flags.String("ldap-groupclass", "GroupOfNames", "Group objectClass")

		tokenLifetimeStr = flags.String("token-ttl", "", "Token lifetime, e.g. '300ms', '1.5h' or '2h45m'")
	)

	flag.Set("logtostderr", "true")
	flags.AddGoFlagSet(flag.CommandLine)
	flags.Parse(os.Args)

	if *ldapHost == "" {
		return nil, errors.New("'ldap-host' parameter is missing")
	}

	tokenLifetime, err := time.ParseDuration(*tokenLifetimeStr)
	if err != nil {
		return nil, fmt.Errorf("'token-ttl' parameter is invalid: %v", err)
	}
	if tokenLifetime > maxTokenTTL {
		return nil, fmt.Errorf("'token-ttl' is too long, must be shorter than '%v'", maxTokenTTL)
	}
	if tokenLifetime == 0 {
		glog.V(1).Info("Using default token lifetime")
		tokenLifetime = defaultTokenTTL
	}

	config := &pion.Config{
		LDAP: pion.LDAPSetting{
			Host:         *ldapHost,
			Port:         *ldapPort,
			Type:         *ldapType,
			BindDN:       *ldapBindDN,
			BindPassword: *ldapBindPassword,
			UserBase:     *ldapUserBase,
			GroupBase:    *ldapGroupBase,
			UserClass:    *ldapUserClass,
			GroupClass:   *ldapGroupClass,
		},
		TokenLifetime: tokenLifetime,
	}

	glog.V(2).Infof("App configuration: '%v'", *config)
	return config, nil
}
