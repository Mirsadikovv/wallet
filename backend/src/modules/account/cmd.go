package account_cmd

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	account_handler "wallet/src/modules/account/handler"
	account_service "wallet/src/modules/account/service"
)

func Cmd(router *echo.Echo, db *gorm.DB) {
	svc := account_service.NewAccountService(db)

	apiGroup := router.Group("/api/v1")
	{
		account_handler.NewAccountHandler(apiGroup, svc)
	}
}
