# Authorization in Pion

## Resources and Actions
Pion defines following resources:
- bucket
- object
- role
- role-binding
- bucket-acl

And actions:
- List
- Get
- Create
- Delete
- Update
- Publish: action for object resource only
- Unpublish: action for object resource only

## Predefined roles
- ReadOnly: users can only list buckets and get objects
- User: users can list, add/remove buckets, list, get and add/remove objects.
- Editor: all permissions of `user` role, plus the `publish/unpublish` buckets
- Admin: all permissions of `editor`, and capable to set bucket ACLs, manage role-bindings (from user/LDAP-group to a role)

## Custom roles

Admin can also create custom roles (TBD) and bind to specific user/groups. This feature is not implemented.	

## ACLs

Beside RBACs, we provide ACLs binding to buckets, which aims for finer-grain access control. The ACLs will be like:
```json
[
	{
	  "id": "1",
	  "actions": [
		"Read"
	  ],
	  "grantees": [
		{
		  "type": "Group",
		  "value": "developers"
		}
	  ]
	},
	{
	  "id": "2",
	  "actions": [
		"Read",
		"Write"
	  ],
	  "grantees": [
		{
		  "type": "User",
		  "value": "alice"
		}
	  ]
	}
]
```

Above ALCs defines that users from `developers` can only read objects in the bucket, the user `alice` can read and 
write objects to the bucket.

### ACL Evaluations

- There are two decision values: `Permit` and `Deny`.
- By default buckets does not have ACLs. Incoming requests from S3 client (like aws-cli) are evaluated against RBAC roles.
- If any bucket has ACLs, requests are also evaluated against these ACLs. Either RBAC evaluation or ACLs evaluation return `Permit`, the final decision is `Permit`.
- There are multiple ACLs per bucket, the ACLs evaluation is permit-override. If there is no matching rule, the default decision is `Deny`.  
