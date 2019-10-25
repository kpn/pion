!/usr/bin/env bash

ldapmodify -a -x -D "cn=admin,dc=example,dc=org" -w ${LDAP_ADMIN_PASSWORD} -H ldap://localhost -f /example.ldif

echo "Verify LDAP bootstrap"
ldapsearch -x -H ldap://localhost -b dc=example,dc=org -D "cn=admin,dc=example,dc=org" -w ${LDAP_ADMIN_PASSWORD}