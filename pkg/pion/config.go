package pion

import (
	"fmt"
	"time"
)

var appConfig Config

type Config struct {
	LDAP          LDAPSetting
	TokenLifetime time.Duration
}

type LDAPSetting struct {
	Host         string
	Port         int
	Type         string
	UserBase     string
	GroupBase    string
	UserClass    string
	GroupClass   string
	BindDN       string
	BindPassword string
}

func (s LDAPSetting) String() string {
	return fmt.Sprintf("Host:'%s', Port:'%d', Type:'%s', UserBase:'%s', GroupBase:'%s', "+
		"UserClass:'%s', GroupClass='%s', BindDN: '%s', BindPassword='Redacted'", s.Host, s.Port, s.Type, s.UserBase,
		s.GroupBase, s.UserClass, s.GroupClass, s.BindDN)
}

// AppConfig returns the singleton app config
func AppConfig() *Config {
	return &appConfig
}

func SetAppConfig(conf *Config) {
	appConfig = *conf
}
