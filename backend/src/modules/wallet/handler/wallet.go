package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	wallet_dto "wallet/src/modules/wallet/dto"
	wallet_service "wallet/src/modules/wallet/service"
)

type walletHandler struct {
	walletService wallet_service.WalletService
}

func NewWalletHandler(router *echo.Group, svc wallet_service.WalletService) {
	h := &walletHandler{
		walletService: svc,
	}

	group := router.Group("/wallet")
	{
		group.POST("", h.Create)
		group.GET("/list", h.List)
		group.GET("/:id", h.GetInfo)
		group.GET("/:id/balance", h.GetBalance)
		group.GET("/:id/transactions", h.GetTransactions)
		group.POST("/:id/send", h.Send)
		group.DELETE("/:id", h.Delete)
	}
}

// Create godoc
// @Summary      Create wallet
// @Description  Create new TON wallet for user
// @Tags         wallet
// @Accept       json
// @Produce      json
// @Param        input body wallet_dto.CreateWalletRequest true "Wallet data"
// @Success      201 {object} wallet_dto.CreateWalletResponse
// @Failure      400 {object} wallet_dto.ErrorResponse
// @Failure      500 {object} wallet_dto.ErrorResponse
// @Router       /api/v1/wallet [POST]
func (h *walletHandler) Create(c echo.Context) error {
	var req wallet_dto.CreateWalletRequest
	{
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, wallet_dto.ErrorResponse{
				Code:    http.StatusBadRequest,
				Error:   "invalid_request",
				Message: err.Error(),
			})
		}
	}

	{
		if err := c.Validate(&req); err != nil {
			return c.JSON(http.StatusBadRequest, wallet_dto.ErrorResponse{
				Code:    http.StatusBadRequest,
				Error:   "validation_failed",
				Message: err.Error(),
			})
		}
	}

	w, err := h.walletService.Create(c.Request().Context(), req.UserID, req.WalletType, req.Network)
	{
		if err != nil {
			return c.JSON(http.StatusInternalServerError, wallet_dto.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Error:   "wallet_creation_failed",
				Message: err.Error(),
			})
		}
	}

	return c.JSON(http.StatusCreated, wallet_dto.CreateWalletResponse{
		ID:         w.ID,
		Address:    w.Address,
		WalletType: w.WalletType,
		Network:    w.Network,
		CreatedAt:  w.CreatedAt,
	})
}

// GetInfo godoc
// @Summary      Get wallet info
// @Description  Returns wallet details including blockchain balance
// @Tags         wallet
// @Produce      json
// @Param        id path int true "Wallet ID"
// @Success      200 {object} wallet_dto.GetWalletInfoResponse
// @Failure      400 {object} wallet_dto.ErrorResponse
// @Failure      404 {object} wallet_dto.ErrorResponse
// @Failure      500 {object} wallet_dto.ErrorResponse
// @Router       /api/v1/wallet/{id} [GET]
func (h *walletHandler) GetInfo(c echo.Context) error {
	walletID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	{
		if err != nil {
			return c.JSON(http.StatusBadRequest, wallet_dto.ErrorResponse{
				Code:    http.StatusBadRequest,
				Error:   "invalid_wallet_id",
				Message: "wallet ID must be a number",
			})
		}
	}

	info, err := h.walletService.GetInfo(c.Request().Context(), walletID)
	{
		if err != nil {
			return c.JSON(http.StatusNotFound, wallet_dto.ErrorResponse{
				Code:    http.StatusNotFound,
				Error:   "wallet_not_found",
				Message: err.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, wallet_dto.GetWalletInfoResponse{
		ID:         info.ID,
		Address:    info.Address,
		Balance:    info.Balance,
		WalletType: info.WalletType,
		Network:    info.Network,
		Seqno:      info.Seqno,
		IsActive:   info.IsActive,
		CreatedAt:  info.CreatedAt,
	})
}

// GetBalance godoc
// @Summary      Get wallet balance
// @Description  Returns current TON balance of the wallet
// @Tags         wallet
// @Produce      json
// @Param        id path int true "Wallet ID"
// @Success      200 {object} wallet_dto.GetBalanceResponse
// @Failure      400 {object} wallet_dto.ErrorResponse
// @Failure      404 {object} wallet_dto.ErrorResponse
// @Failure      500 {object} wallet_dto.ErrorResponse
// @Router       /api/v1/wallet/{id}/balance [GET]
func (h *walletHandler) GetBalance(c echo.Context) error {
	walletID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	{
		if err != nil {
			return c.JSON(http.StatusBadRequest, wallet_dto.ErrorResponse{
				Code:    http.StatusBadRequest,
				Error:   "invalid_wallet_id",
				Message: "wallet ID must be a number",
			})
		}
	}

	w, err := h.walletService.FindByID(c.Request().Context(), walletID)
	{
		if err != nil {
			return c.JSON(http.StatusNotFound, wallet_dto.ErrorResponse{
				Code:    http.StatusNotFound,
				Error:   "wallet_not_found",
				Message: err.Error(),
			})
		}
	}

	balance, err := h.walletService.GetBalance(c.Request().Context(), walletID)
	{
		if err != nil {
			return c.JSON(http.StatusInternalServerError, wallet_dto.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Error:   "failed_to_get_balance",
				Message: err.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, wallet_dto.GetBalanceResponse{
		Address: w.Address,
		Balance: balance,
	})
}

// List godoc
// @Summary      List user wallets
// @Description  Returns all active wallets for a user
// @Tags         wallet
// @Produce      json
// @Param        user_id query int true "User ID"
// @Success      200 {object} wallet_dto.ListWalletsResponse
// @Failure      400 {object} wallet_dto.ErrorResponse
// @Failure      500 {object} wallet_dto.ErrorResponse
// @Router       /api/v1/wallet/list [GET]
func (h *walletHandler) List(c echo.Context) error {
	userID, err := strconv.ParseInt(c.QueryParam("user_id"), 10, 64)
	{
		if err != nil {
			return c.JSON(http.StatusBadRequest, wallet_dto.ErrorResponse{
				Code:    http.StatusBadRequest,
				Error:   "invalid_user_id",
				Message: "user ID must be a number",
			})
		}
	}

	wallets, err := h.walletService.FindByUserID(c.Request().Context(), userID)
	{
		if err != nil {
			return c.JSON(http.StatusInternalServerError, wallet_dto.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Error:   "failed_to_get_wallets",
				Message: err.Error(),
			})
		}
	}

	summaries := make([]wallet_dto.WalletSummary, 0, len(wallets))
	for _, w := range wallets {
		summaries = append(summaries, wallet_dto.WalletSummary{
			ID:         w.ID,
			Address:    w.Address,
			WalletType: w.WalletType,
			Network:    w.Network,
			IsActive:   w.IsActive,
			CreatedAt:  w.CreatedAt,
		})
	}

	return c.JSON(http.StatusOK, wallet_dto.ListWalletsResponse{
		Wallets: summaries,
		Total:   len(summaries),
	})
}

// Delete godoc
// @Summary      Delete wallet
// @Description  Deactivates wallet (soft delete)
// @Tags         wallet
// @Produce      json
// @Param        id path int true "Wallet ID"
// @Success      200 {object} wallet_dto.SuccessResponse
// @Failure      400 {object} wallet_dto.ErrorResponse
// @Failure      500 {object} wallet_dto.ErrorResponse
// @Router       /api/v1/wallet/{id} [DELETE]
func (h *walletHandler) Delete(c echo.Context) error {
	walletID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	{
		if err != nil {
			return c.JSON(http.StatusBadRequest, wallet_dto.ErrorResponse{
				Code:    http.StatusBadRequest,
				Error:   "invalid_wallet_id",
				Message: "wallet ID must be a number",
			})
		}
	}

	{
		if err := h.walletService.Delete(c.Request().Context(), walletID); err != nil {
			return c.JSON(http.StatusInternalServerError, wallet_dto.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Error:   "failed_to_delete_wallet",
				Message: err.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, wallet_dto.SuccessResponse{
		Success: true,
		Message: "wallet deleted successfully",
	})
}

// GetTransactions godoc
// @Summary      Get transaction history
// @Description  Returns list of wallet transactions from TON blockchain
// @Tags         wallet
// @Produce      json
// @Param        id    path  int true  "Wallet ID"
// @Param        limit query int false "Number of transactions (default 10, max 100)"
// @Success      200 {object} wallet_dto.GetTransactionsResponse
// @Failure      400 {object} wallet_dto.ErrorResponse
// @Failure      404 {object} wallet_dto.ErrorResponse
// @Failure      500 {object} wallet_dto.ErrorResponse
// @Router       /api/v1/wallet/{id}/transactions [GET]
func (h *walletHandler) GetTransactions(c echo.Context) error {
	walletID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	{
		if err != nil {
			return c.JSON(http.StatusBadRequest, wallet_dto.ErrorResponse{
				Code:    http.StatusBadRequest,
				Error:   "invalid_wallet_id",
				Message: "wallet ID must be a number",
			})
		}
	}

	var req wallet_dto.GetTransactionsRequest
	{
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, wallet_dto.ErrorResponse{
				Code:    http.StatusBadRequest,
				Error:   "invalid_request",
				Message: err.Error(),
			})
		}
	}

	limit := 10
	if req.Limit > 0 {
		limit = req.Limit
	}

	w, err := h.walletService.FindByID(c.Request().Context(), walletID)
	{
		if err != nil {
			return c.JSON(http.StatusNotFound, wallet_dto.ErrorResponse{
				Code:    http.StatusNotFound,
				Error:   "wallet_not_found",
				Message: err.Error(),
			})
		}
	}

	transactions, err := h.walletService.GetTransactions(c.Request().Context(), walletID, limit)
	{
		if err != nil {
			return c.JSON(http.StatusInternalServerError, wallet_dto.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Error:   "failed_to_get_transactions",
				Message: err.Error(),
			})
		}
	}

	txDTOs := make([]*wallet_dto.TransactionDTO, 0, len(transactions))
	for _, tx := range transactions {
		txDTOs = append(txDTOs, &wallet_dto.TransactionDTO{
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

	return c.JSON(http.StatusOK, wallet_dto.GetTransactionsResponse{
		WalletID:     w.ID,
		Address:      w.Address,
		Transactions: txDTOs,
		Total:        len(txDTOs),
	})
}

// Send godoc
// @Summary      Send TON coins
// @Description  Sends TON from wallet to recipient address
// @Tags         wallet
// @Accept       json
// @Produce      json
// @Param        id    path int                          true "Wallet ID"
// @Param        input body wallet_dto.SendCoinsRequest true "Transaction data"
// @Success      200 {object} wallet_dto.SendCoinsResponse
// @Failure      400 {object} wallet_dto.ErrorResponse
// @Failure      404 {object} wallet_dto.ErrorResponse
// @Failure      500 {object} wallet_dto.ErrorResponse
// @Router       /api/v1/wallet/{id}/send [POST]
func (h *walletHandler) Send(c echo.Context) error {
	walletID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	{
		if err != nil {
			return c.JSON(http.StatusBadRequest, wallet_dto.ErrorResponse{
				Code:    http.StatusBadRequest,
				Error:   "invalid_wallet_id",
				Message: "wallet ID must be a number",
			})
		}
	}

	var req wallet_dto.SendCoinsRequest
	{
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, wallet_dto.ErrorResponse{
				Code:    http.StatusBadRequest,
				Error:   "invalid_request",
				Message: err.Error(),
			})
		}
	}

	{
		if err := c.Validate(&req); err != nil {
			return c.JSON(http.StatusBadRequest, wallet_dto.ErrorResponse{
				Code:    http.StatusBadRequest,
				Error:   "validation_failed",
				Message: err.Error(),
			})
		}
	}

	result, err := h.walletService.Send(c.Request().Context(), walletID, req.Recipient, req.Amount, req.Comment)
	{
		if err != nil {
			return c.JSON(http.StatusInternalServerError, wallet_dto.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Error:   "failed_to_send_coins",
				Message: err.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, wallet_dto.SendCoinsResponse{
		Hash:      result.Hash,
		Lt:        result.Lt,
		Address:   result.Address,
		Amount:    result.Amount,
		Fee:       result.Fee,
		Recipient: result.Recipient,
		Comment:   result.Comment,
	})
}
