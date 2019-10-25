# Deployment instruction

## Deploy LDAP

Create LDAP service in the namespace:
```bash
kubectl apply -f dev-openldap-manifest.yaml
```

Add the LDAP schema and sample LDAP entities:
```bash
OPENLDAP_POD=$(kubectl get pod -l app=openldap -o jsonpath={.items[*].metadata.name})
kubectl cp example.ldif ${OPENLDAP_POD}:/
kubectl cp ldap-bootstrap.sh ${OPENLDAP_POD}:/
kubectl exec ${OPENLDAP_POD} bash /ldap-bootstrap.sh
```

Verify if users and groups are added:
```bash
ldapsearch -x -H ldap://localhost -b dc=example,dc=org -D "cn=admin,dc=example,dc=org" -w ${LDAP_ADMIN_PASSWORD}
```

## Deploy Minio cluster

```bash
kubectl apply -f dev-minio-manifest.yaml
```

## Deploy Etcd cluster

```bash
kubectl apply -f etcd-persistence-cluster.yaml
```

## Deploy Pion services

- In the `dev-pion-proxy` ingress, change the host to proper value. The default is `pion-gw.example.org`
- In the `dev-pion-ui` ingress, change the host to proper value. The default is `pion.example.org`

Deploy services with following command:

```bash
kubectl apply -f dev-pion-manifest.yaml
```

### Create example customers

Forward port of the manager service:

```bash
$ k port-forward dev-pion-manager-588d9c784c-w8qsp 8080
```

Create customer-1:
```bash
curl -X POST \
  http://localhost:8080/_internal/customers \
  -H 'Content-Type: application/json' \
  -d '{
	"name": "customer1",
	"groups": ["group_customer1"],
	"userIDs": []
}'

curl -X POST \
  http://localhost:8080/_internal/customers \
  -H 'Content-Type: application/json' \
  -d '{
	"name": "customer2",
	"groups": ["group_customer2"],
	"userIDs": []
}'

```

List current customers:
```bash
$ curl http://localhost:8080/_internal/customers

[
  {
    "name": "customer1",
    "createdAt": "2019-10-07T14:42:42.64511046Z",
    "modifiedAt": "2019-10-07T14:42:42.64511046Z",
    "groups": [
      "group_customer1"
    ],
    "userIDs": []
  },
  {
    "name": "customer2",
    "createdAt": "2019-10-07T14:43:02.861132768Z",
    "modifiedAt": "2019-10-07T14:43:02.861132768Z",
    "groups": [
      "group_customer2"
    ],
    "userIDs": []
  },
  {
    "name": "pion",
    "createdAt": "2019-10-04T11:01:46.885781632Z",
    "modifiedAt": "2019-10-04T11:01:46.885781632Z",
    "groups": [
      "pion_users",
      "pion_admins"
    ],
    "userIDs": null
  }
]
```

### Bootstrap the default administrator

This command set `billy` as the admin of all customers
```bash
$ kubectl exec $(kubectl get pod -l component=manager -o jsonpath={.items[*].metadata.name}) -- /opt/bootstrap --admin-user billy
I1009 14:58:26.621235      10 bootstrap.go:22] Adding 'billy' as the Admin of the system
I1009 14:58:26.630980      10 bootstrap.go:30] Adding system level role-binding
I1009 14:58:26.645739      10 bootstrap.go:34] Role-binding 'default-admin-rb' existed, deleting
I1009 14:58:26.703291      10 bootstrap.go:45] Created role-binding: {default-admin-rb 0001-01-01 00:00:00 +0000 UTC [{user billy}] admin}
```

## Testing
1. Login to customer `customer1` with user/password `billy/billy`.
2. Create a role-binding for user alice as editor.
3. Login to customer `customer1` with user `alice/alice`
4. Create a new bucket: `customer1-testbucket`
5. Generate access/secret keys for billy.
6. Remember to use correct region and signature version:
    ```bash
    $ cat ~/.aws/config
    [default]
    region = us-east-1
    signature_version = s3v4
    ```

7. Use AWS CLI to access the created bucket: 
    ```bash
    export AWS_ACCESS_KEY_ID=pion-0f286e8a39e083fa840919bcd8545329
    export AWS_SECRET_ACCESS_KEY=276957a66127ea8f171a45445995af3ebd7f0ef4c5d5be416d0fc400c95b0e5a
    
    aws s3 --endpoint=https://pion-gw.example.org --no-verify-ssl ls 
    ```
   
8. Login to `customer2` with the super-admid `billy/billy`.
9. Create role-binding for group `group_customer2` as `User` role
10. Logout and login again to `customer2` with user `dave/dave`.
11. Generate access key for `dave`
12. Use AWS CLI to access the buckets of the `customer2`:
    ```bash
        export AWS_ACCESS_KEY_ID=pion-efef06346b8150c5e1d19d3f2281e1c5
        export AWS_SECRET_ACCESS_KEY=663da98784afdd0f71eba4474e36dc358c13f3adfe06c6dbcfb6135c72002dfc
        
        aws s3 --endpoint=https://pion-gw.example.org --no-verify-ssl mb s3://customer2-bucket
    aws s3 --endpoint=https://pion-gw.example.org --no-verify-ssl ls
    ```
13. Verify that buckets created by `customer1` does not show to users of `customer2`
   