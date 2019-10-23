package token

import (
	"encoding/json"

	"github.com/golang/glog"
	"github.com/pkg/errors"
)

type Payload struct {
	Username string `json:"username"`
	Data     string `json:"data"`
}

func NewPayload(jsonString string) (p Payload, err error) {
	if jsonString == "" {
		return p, errors.New("Empty json value")
	}
	err = json.Unmarshal([]byte(jsonString), &p)
	if err != nil {
		return p, err
	}
	return p, nil
}

// serialize to string
func (t Payload) String() string {
	jsonString, err := json.Marshal(t)
	if err != nil {
		glog.Errorf("Failed to serialize JSON: %v", err)
		return ""
	}
	return string(jsonString[:])
}
