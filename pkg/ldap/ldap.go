package ldap

import (
	"crypto/tls"
	"errors"
	"fmt"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/pion"
	"gopkg.in/ldap.v2"
)

type LDAPClient struct {
	setting pion.LDAPSetting
	conn    *ldap.Conn
}

// NewLDAPClient creates a new LDAP client connecting to KPN LDAP server. Caller must call client.Close() after using.
func NewLDAPClient(setting pion.LDAPSetting) (*LDAPClient, error) {
	client := &LDAPClient{
		setting: setting,
	}

	if err := client.connect(); err != nil {
		return nil, err
	}
	return client, nil
}

// NewAndBindLDAPClient creates a new LDAP client and binds with the configured system credential. Caller must call
// client.Close() after using.
func NewAndBindLDAPClient(setting pion.LDAPSetting) (*LDAPClient, error) {
	c, err := NewLDAPClient(setting)
	if err != nil {
		return nil, err
	}
	err = c.conn.Bind(setting.BindDN, setting.BindPassword)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c LDAPClient) Close() {
	c.conn.Close()
}

func (c *LDAPClient) connect() (err error) {
	c.conn, err = c.dial()
	return err
}

func (c LDAPClient) dial() (*ldap.Conn, error) {
	connString := fmt.Sprintf("%s:%d", c.setting.Host, c.setting.Port)
	glog.V(2).Infof("LDAP dialing to '%s'", connString)

	if c.setting.Type == "ldaps" {
		return ldap.DialTLS("tcp", connString, &tls.Config{
			// TODO install KPN private CA certificate to remove this option
			InsecureSkipVerify: true,
		})
	}
	conn, err := ldap.Dial("tcp", connString)
	if err != nil {
		return nil, err
	}
	if conn == nil {
		return nil, errors.New("cannot dial LDAP, nil connection")
	}
	return conn, nil
}

// do initializes the connection to LDAP server and invokes the provided function with the connection.
func (c LDAPClient) do(f func(c *ldap.Conn) error) error {
	if err := c.conn.Bind(c.setting.BindDN, c.setting.BindPassword); err != nil {
		if c.setting.BindDN == "" && c.setting.BindPassword == "" {
			return fmt.Errorf("ldap: initial anonymous bind failed: %v", err)
		}
		return fmt.Errorf("ldap: initial bind for user %q failed: %v", c.setting.BindDN, err)
	}
	return f(c.conn)
}

// Login authenticates against the LDAP server with the given username and password
func (c LDAPClient) Login(username, password string) error {
	if password == "" {
		return errors.New("prevented unauthenticated binding to LDAP server")
	}
	err := c.do(func(conn *ldap.Conn) error {
		userDN := fmt.Sprintf("cn=%s,%s", username, c.setting.UserBase)

		// authenticate as distinguished name
		if err := conn.Bind(userDN, password); err != nil {
			if ldapErr, ok := err.(*ldap.Error); ok {
				switch ldapErr.ResultCode {
				case ldap.LDAPResultInvalidCredentials:
					glog.Errorf("ldap: invalid password for user %q", userDN)
					return fmt.Errorf("ldap: invalid password")
				case ldap.LDAPResultConstraintViolation:
					glog.Errorf("ldap: constraint violation for user %q: %s", userDN, ldapErr.Error())
					return fmt.Errorf("ldap: constraint violation")
				}
			}
			return fmt.Errorf("ldap: failed to bind as dn %q: %v", userDN, err)
		}
		glog.V(1).Infof("User %s authenticated successfully!", username)
		return nil
	})
	return err
}

func (c LDAPClient) QueryLDAP(query string, attributes []string, base string) (*ldap.SearchResult, error) {
	var (
		result *ldap.SearchResult
	)

	err := c.do(func(conn *ldap.Conn) error {
		var err error
		req := ldap.NewSearchRequest(
			base,
			ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0,
			false,
			query,
			attributes,
			nil,
		)

		result, err = conn.Search(req)
		return err
	})
	if err != nil {
		glog.Errorf("ldap: query failed: %v", err)
		return nil, err
	}
	return result, nil
}
