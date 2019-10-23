package user_context

import (
	"errors"

	"github.com/kpn/pion/pkg/pion/shared"
	"github.com/labstack/echo"
)

func GetData(context echo.Context) (userId string, userGroups []string, customerName string, err error) {
	userId, ok := context.Get(shared.UserIdKey).(string)
	if !ok {
		return userId, userGroups, customerName, errors.New("failed to get userId from context")
	}
	userGroups, ok = context.Get(shared.UserGroupKey).([]string)
	if !ok {
		return userId, userGroups, customerName, errors.New("failed to get userGroups from context")
	}
	customerName, ok = context.Get(shared.CustomerKey).(string)
	if !ok {
		return userId, userGroups, customerName, errors.New("failed to get customer from context")
	}
	return userId, userGroups, customerName, nil
}
