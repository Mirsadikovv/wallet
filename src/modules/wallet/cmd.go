package wallet_cmd

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"

	"wallet_test/src/modules/wallet/handler"
	"wallet_test/src/modules/wallet/service"
)

// @title Wallet Service API
// @version 1.0
// @description This is a Wallet Service API for managing TON blockchain wallets
// @BasePath  /api/v1
// @Schemes http https
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func Cmd(router *gin.Engine, db *bun.DB, network string, encryptionKey string) {
	walletService, err := service.NewWalletService(db, network, encryptionKey)
	if err != nil {
		log.Fatalf("Failed to create wallet service: %v", err)
	}

	walletHandler := handler.NewWalletHandler(walletService)

	walletGroup := router.Group("/api/v1/wallet")
	{
		// Создать кошелек
		walletGroup.POST("", walletHandler.CreateWallet)

		// Получить информацию о кошельке
		walletGroup.GET("/:id", walletHandler.GetWalletInfo)

		// Получить баланс кошелька
		walletGroup.GET("/:id/balance", walletHandler.GetBalance)

		// Получить историю транзакций кошелька
		walletGroup.GET("/:id/transactions", walletHandler.GetTransactions)

		// Отправить TON монеты
		walletGroup.POST("/:id/send", walletHandler.SendCoins)

		// Список кошельков пользователя
		walletGroup.GET("/list", walletHandler.ListUserWallets)

		// Удалить кошелек
		walletGroup.DELETE("/:id", walletHandler.DeleteWallet)
	}
}
