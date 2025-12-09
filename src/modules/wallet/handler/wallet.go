package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"wallet_test/src/modules/wallet/dto"
	"wallet_test/src/modules/wallet/service"
)

type WalletHandler struct {
	walletService *service.WalletService
}

func NewWalletHandler(walletService *service.WalletService) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
	}
}

// CreateWallet создает новый кошелек
// @Summary Создать новый кошелек
// @Description Создает новый TON кошелек для пользователя
// @Tags wallet
// @Accept json
// @Produce json
// @Param request body dto.CreateWalletRequest true "Данные для создания кошелька"
// @Success 201 {object} dto.CreateWalletResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/wallet [post]
func (h *WalletHandler) CreateWallet(c *gin.Context) {
	var req dto.CreateWalletRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	wallet, err := h.walletService.CreateWallet(c.Request.Context(), req.UserID, req.WalletType, req.Network)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "wallet_creation_failed",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusCreated, dto.CreateWalletResponse{
		ID:         wallet.ID,
		Address:    wallet.Address,
		WalletType: wallet.WalletType,
		Network:    wallet.Network,
		CreatedAt:  wallet.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// GetWalletInfo получает информацию о кошельке
// @Summary Получить информацию о кошельке
// @Description Возвращает детальную информацию о кошельке включая баланс
// @Tags wallet
// @Accept json
// @Produce json
// @Param id path int true "ID кошелька"
// @Success 200 {object} dto.GetWalletInfoResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/wallet/{id} [get]
func (h *WalletHandler) GetWalletInfo(c *gin.Context) {
	walletIDStr := c.Param("id")
	walletID, err := strconv.ParseInt(walletIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_wallet_id",
			Message: "ID кошелька должен быть числом",
			Code:    http.StatusBadRequest,
		})
		return
	}

	wallet, err := h.walletService.GetWalletByID(c.Request.Context(), walletID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "wallet_not_found",
			Message: err.Error(),
			Code:    http.StatusNotFound,
		})
		return
	}

	info, err := h.walletService.GetWalletInfo(c.Request.Context(), walletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "failed_to_get_wallet_info",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, dto.GetWalletInfoResponse{
		ID:         wallet.ID,
		Address:    info.Address,
		Balance:    info.Balance,
		WalletType: info.WalletType,
		Network:    wallet.Network,
		Seqno:      info.Seqno,
		IsActive:   wallet.IsActive,
		CreatedAt:  wallet.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// GetBalance получает баланс кошелька
// @Summary Получить баланс кошелька
// @Description Возвращает текущий баланс кошелька в TON
// @Tags wallet
// @Accept json
// @Produce json
// @Param id path int true "ID кошелька"
// @Success 200 {object} dto.GetBalanceResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/wallet/{id}/balance [get]
func (h *WalletHandler) GetBalance(c *gin.Context) {
	walletIDStr := c.Param("id")
	walletID, err := strconv.ParseInt(walletIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_wallet_id",
			Message: "ID кошелька должен быть числом",
			Code:    http.StatusBadRequest,
		})
		return
	}

	wallet, err := h.walletService.GetWalletByID(c.Request.Context(), walletID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "wallet_not_found",
			Message: err.Error(),
			Code:    http.StatusNotFound,
		})
		return
	}

	balance, err := h.walletService.GetBalance(c.Request.Context(), walletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "failed_to_get_balance",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, dto.GetBalanceResponse{
		Address: wallet.Address,
		Balance: balance,
	})
}

// ListUserWallets получает список кошельков пользователя
// @Summary Список кошельков пользователя
// @Description Возвращает все активные кошельки пользователя
// @Tags wallet
// @Accept json
// @Produce json
// @Param user_id query int true "ID пользователя"
// @Success 200 {object} dto.ListWalletsResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/wallet/list [get]
func (h *WalletHandler) ListUserWallets(c *gin.Context) {
	userIDStr := c.Query("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_user_id",
			Message: "ID пользователя должен быть числом",
			Code:    http.StatusBadRequest,
		})
		return
	}

	wallets, err := h.walletService.GetUserWallets(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "failed_to_get_wallets",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	summaries := make([]dto.WalletSummary, 0, len(wallets))
	for _, w := range wallets {
		summaries = append(summaries, dto.WalletSummary{
			ID:         w.ID,
			Address:    w.Address,
			WalletType: w.WalletType,
			Network:    w.Network,
			IsActive:   w.IsActive,
			CreatedAt:  w.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	c.JSON(http.StatusOK, dto.ListWalletsResponse{
		Wallets: summaries,
		Total:   len(summaries),
	})
}

// DeleteWallet удаляет кошелек (мягкое удаление)
// @Summary Удалить кошелек
// @Description Деактивирует кошелек (мягкое удаление)
// @Tags wallet
// @Accept json
// @Produce json
// @Param id path int true "ID кошелька"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/wallet/{id} [delete]
func (h *WalletHandler) DeleteWallet(c *gin.Context) {
	walletIDStr := c.Param("id")
	walletID, err := strconv.ParseInt(walletIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_wallet_id",
			Message: "ID кошелька должен быть числом",
			Code:    http.StatusBadRequest,
		})
		return
	}

	err = h.walletService.DeleteWallet(c.Request.Context(), walletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "failed_to_delete_wallet",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Кошелек успешно удален",
	})
}

// GetTransactions получает историю транзакций кошелька
// @Summary Получить историю транзакций
// @Description Возвращает список транзакций кошелька из блокчейна TON
// @Tags wallet
// @Accept json
// @Produce json
// @Param id path int true "ID кошелька"
// @Param limit query int false "Количество транзакций (по умолчанию 10, макс 100)" default(10)
// @Success 200 {object} dto.GetTransactionsResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/wallet/{id}/transactions [get]
func (h *WalletHandler) GetTransactions(c *gin.Context) {
	walletIDStr := c.Param("id")
	walletID, err := strconv.ParseInt(walletIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_wallet_id",
			Message: "ID кошелька должен быть числом",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var req dto.GetTransactionsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Устанавливаем значение по умолчанию
	limit := 10
	if req.Limit > 0 {
		limit = req.Limit
	}

	wallet, err := h.walletService.GetWalletByID(c.Request.Context(), walletID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "wallet_not_found",
			Message: err.Error(),
			Code:    http.StatusNotFound,
		})
		return
	}

	transactions, err := h.walletService.GetTransactions(c.Request.Context(), walletID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "failed_to_get_transactions",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Преобразуем в DTO
	txDTOs := make([]*dto.TransactionDTO, 0, len(transactions))
	for _, tx := range transactions {
		txDTOs = append(txDTOs, &dto.TransactionDTO{
			Hash:      tx.Hash,
			Lt:        tx.Lt,
			Timestamp: tx.Timestamp,
			Type:      tx.Type,
			Amount:    tx.Amount,
			Fee:       tx.Fee,
			From:      tx.From,
			To:        tx.To,
			Comment:   tx.Comment,
			Success:   tx.Success,
		})
	}

	c.JSON(http.StatusOK, dto.GetTransactionsResponse{
		WalletID:     wallet.ID,
		Address:      wallet.Address,
		Transactions: txDTOs,
		Total:        len(txDTOs),
	})
}
