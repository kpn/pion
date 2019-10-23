import {InMemoryDbService} from 'angular-in-memory-web-api';
import {AccessKey} from '../model/accesskey';
import {Bucket} from '../model/bucket';
import {RoleBinding} from "../model/role-binding";
import {Role} from "../model/role";

export class InMemoryBackendService implements InMemoryDbService {
  createDb() {
    const ACCESS_KEYS: AccessKey[] = [
      {
        createdAt: new Date('2018-09-18T17:44:38+02:00'),
        accessKeyId: 'toss-0c493013c5af61f218b087aa2c46d21f',
        lastUsed: 'N/A',
        status: 'Active',
      },
      {
        createdAt: new Date('2018-09-18T17:44:38+02:00'),
        accessKeyId: 'toss-0c493013c5af61f218b087aa2c46d21f',
        lastUsed: 'N/A',
        status: 'Active',
      },
      {
        createdAt: new Date('2018-09-18T17:44:38+02:00'),
        accessKeyId: 'toss-0c493013c5af61f218b087aa2c46d21f',
        lastUsed: 'N/A',
        status: 'Active',
      },
    ];

    const BUCKETS: Bucket[] = [
      {
        name: 'shared',
        ownedBy: 'foo',
        createdAt: '2018-09-18T17:44:38+02:00',
        creator: 'ngo500',
        modifiedAt: '2018-09-20T17:44:38+02:00',
        acls: ''
      },
      {
        name: 'backups',
        ownedBy: 'foo',
        createdAt: '2018-09-17T17:44:38+02:00',
        creator: 'admin',
        modifiedAt: '2018-09-21T17:44:38+02:00',
        acls: ''
      }
    ];

    const ROLE_BINDINGS: RoleBinding[] = [
      {
        name: "editor-rb",
        createdAt: new Date("2018-11-06T15:22:18.544446432Z"),
        subjects: [
          {
            "type": "group",
            "value": "dig_foo"
          }
        ],
        roleRef: "editor"
      },
      {
        name: "dev-users",
        createdAt: new Date("2018-11-06T16:34:48.979812508Z"),
        "subjects": [
          {
            "type": "group",
            "value": "dig_foo"
          },
          {
            "type": "user",
            "value": "john"
          },
          {
            "type": "group",
            "value": "dig_bah"
          },
          {
            "type": "group",
            "value": "dig_data"
          }
        ],
        roleRef: "user"
      }
    ];

    const roleDate = new Date("2018-11-06T12:00:00Z")
    const ROLES: Role[] = [
      {
        name: "readonly",
        displayName: "ReadOnly",
        createdAt: roleDate,
        rules: [
          {
            resources: ["oss:bucket", "oss:object"],
            actions: ["list", "get"]
          }
        ]
      },
      {
        name: "editor",
        displayName: "Editor",
        createdAt: roleDate,
        rules: [
          {
            resources: ["oss:bucket"],
            actions: ["list", "get", "create", "delete", "update", "publish", "unpublish"]
          },
          {
            resources: [
              "oss:object"
            ],
            actions: ["list", "get", "create", "delete", "update"]
          }
        ]
      },
    ];

    return {
      access_keys: ACCESS_KEYS,
      buckets: BUCKETS,
      role_bindings: ROLE_BINDINGS,
      roles: ROLES,
    };
  }

}
