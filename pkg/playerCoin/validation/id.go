package validation

import (
	"github.com/labstack/echo/v4"
	_adminException "github.com/ritchie-gr8/fantasy-shop-api/pkg/admin/exception"
	_playerException "github.com/ritchie-gr8/fantasy-shop-api/pkg/player/exception"
)

// userType must be "adminID" or "playerID"
func GetUserID(pctx echo.Context, userType string) (string, error) {
	userID, ok := pctx.Get(userType).(string)
	if !ok || userID == "" {
		if userType == "adminID" {
			return "", &_adminException.AdminNotFound{AdminID: "unknow"}
		} else {
			return "", &_playerException.PlayerNotFound{PlayerID: "unknow"}
		}
	}

	return userID, nil
}
