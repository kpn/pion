# Manager

Admins can manage customers via customers APIs in the Manager component.

## Customer APIs
### List customers

```bash
curl -X GET http://localhost:8080/_internal/customers
```

### Add new customer

```bash
curl -X POST \
  http://localhost:8080/_internal/customers \
  -H 'Content-Type: application/json' \
  -d '{
	"name": "customer1",
	"groups": ["ldap_group1", "ldap_group2"],
	"userIDs": ["user_id1", "user_id2]
}'
```

### Delete customer

```bash
curl -X DELETE \
  http://localhost:8080/_internal/customers \
  -H 'Content-Type: application/json' \
  -d '{
	"name": "customer1"
}'
```

### Update customer

```bash
curl -X PUT \
  http://localhost:8080/_internal/customers/customer1 \
  -H 'Content-Type: application/json' \
  -d '{
	"groups": ["new_ldap_group1", "new_ldap_group2"],
	"userIDs": ["new_user1", "new_user2"]
}'
```

## Customer APIs