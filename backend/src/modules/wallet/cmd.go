package wallet_cmd

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	wallet_handler "wallet/src/modules/wallet/handler"
	wallet_service "wallet/src/modules/wallet/service"
)

func Cmd(router *echo.Echo, db *gorm.DB, network, encryptionKey string) error {
	walletSvc, err := wallet_service.NewWalletService(db, network, encryptionKey)
	if err != nil {
		return fmt.Errorf("failed to initialize wallet service: %w", err)
	}

	apiGroup := router.Group("/api/v1")
	{
		wallet_handler.NewWalletHandler(apiGroup, walletSvc)
	}

	return nil
}
