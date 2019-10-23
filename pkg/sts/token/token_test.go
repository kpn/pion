package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSerialize(t *testing.T) {
	p := Payload{
		Username: "canh.ngo",
		Data:     "foo",
	}

	assert.Equal(t, `{"username":"canh.ngo","data":"foo"}`, p.String())

}

func TestNewTokenPayload(t *testing.T) {
	json := `{"username":"canh.ngo","data":"foo"}`

	p, err := NewPayload(json)
	assert.NoError(t, err)
	assert.Equal(t, "canh.ngo", p.Username)
	assert.Equal(t, "foo", p.Data)
}
