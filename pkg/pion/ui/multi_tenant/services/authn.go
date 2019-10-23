package services

import (
	"github.com/kpn/pion/pkg/pion/ui/multi_tenant"
)

type AuthnService interface {
	Login(customerName string, username, password string) (*multi_tenant.UserInfo, error)
}
