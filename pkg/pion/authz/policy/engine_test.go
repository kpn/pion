package policy_test

import (
	"testing"

	"github.com/kpn/pion/pkg/pion/authz"
	"github.com/kpn/pion/pkg/pion/authz/policy"
	"github.com/kpn/pion/pkg/pion/authz/watcher/mock"
	"github.com/kpn/pion/pkg/pion/shared/multi-tenancy/model"
	"github.com/stretchr/testify/assert"
)

func TestEngine(t *testing.T) {
	policiesJSON := `
{
  "mybucket": [
    {
      "id": "1",
      "actions": [
        "Read",
        "Write"
      ],
      "grantees": [
        {
          "type": "User",
          "value": "ngo500"
        }
      ]
    }
  ],
  "share-bucket": [
    {
      "id": "1",
      "actions": [
        "Read"
      ],
      "grantees": [
        {
          "type": "Group",
          "value": "dig_infraplatform"
        }
      ]
    }
  ]
}
`
	policies, err := authz.NewPolicies([]byte(policiesJSON))
	assert.NoError(t, err)

	mbw := mock.NewBucketsWatcher()
	mbw.On("GetBucket", "mybucket").Return(&model.Bucket{
		Name: "mybucket",
		ACLs: policies["mybucket"],
	}).Once()

	de, err := policy.NewEngine(mbw)
	assert.NoError(t, err)

	decision := de.Evaluate(authz.Request{
		Username: "ngo500",
		Groups:   []string{"infra"},
		Action:   authz.Read,
		Target:   "mybucket",
	})
	assert.Equal(t, authz.DecisionPermit, decision)
}
