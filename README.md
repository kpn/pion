# Pion - Object Storage Service Gateway

## Features
- Support S3 HTTP protocol.
- Multi-tenancy: Pion supports multi-tenancy. Buckets in customer accounts are isolated. When logging in, users must provide 
the customer account they want to login to.
- Unique bucket names among customers. It means if a customer already had a bucket name `myawesomebucket`, in the other 
customer account, you cannot create a bucket with that name anymore.
- Default access/secret token lifetime is 90 days.

### User management and Access Control
- Customer account: is a group of users in one or multiple LDAP groups. Buckets among customers are isolated.
- RBAC management: a user can be assigned to differet predefined roles in the account. For more detail, please see [here](docs/authorization.md).
 	
### Object Storage Service
The solution consists of following components:
 - Security Token Service: this service allows to create and verify tokens binding to authenticated users.
 - UI: the dashboard to manage user access keys and authorization policies (TBD). Users can login to the dashboard
 by their credentials.
 - Proxy: The proxy runs in front of the Minio cluster to authenticate (via STS) and authorize (via Authz service) 
 incoming requests from clients (Minio client or AWS-CLI S3). Validated requests are forwarded to the upstream Minio cluster.
 - Authorization service: this service manages authorization policies for buckets. It has an authorization API endpoint 
 serving request from the Proxy
 - Manager service: to manage public buckets, which can be accessed directly via URLs.
 
## Deployment
Instructions for deploying Pion can be found [here](docs/deploy.md).

You can also find example deployment at [k8s folder](k8s/dev) 

## Build
Requirements for building
- Go (built with 1.12.4)
- dep (v0.5.4) for dependency management.
- UI: npm (v6.12.0), angular-cli (v7.0.3), Node (v11.14.0)

For detail, please find [here](docs/build.md).