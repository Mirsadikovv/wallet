package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	account_dto "wallet/src/modules/account/dto"
	account_service "wallet/src/modules/account/service"
)

type accountHandler struct {
	accountService account_service.AccountService
}

func NewAccountHandler(router *echo.Group, svc account_service.AccountService) {
	h := &accountHandler{
		accountService: svc,
	}

	group := router.Group("/account")
	{
		group.POST("", h.Create)
		group.GET("/list", h.List)
		group.GET("/:id", h.GetByID)
		group.PUT("/:id", h.Update)
		group.DELETE("/:id", h.DeleteOrRestore)
	}
}

// Create godoc
// @Summary      Create account
// @Description  Create new account for user
// @Tags         account
// @Accept       json
// @Produce      json
// @Param        input body account_dto.CreateAccountRequest true "Account data"
// @Success      201 {object} account_dto.AccountResponse
// @Failure      400 {object} account_dto.ErrorResponse
// @Failure      500 {object} account_dto.ErrorResponse
// @Router       /api/v1/account [POST]
func (h *accountHandler) Create(c echo.Context) error {
	var req account_dto.CreateAccountRequest
	{
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, account_dto.ErrorResponse{
				Code:    http.StatusBadRequest,
				Error:   "invalid_request",
				Message: err.Error(),
			})
		}
	}

	{
		if err := c.Validate(&req); err != nil {
			return c.JSON(http.StatusBadRequest, account_dto.ErrorResponse{
				Code:    http.StatusBadRequest,
				Error:   "validation_failed",
				Message: err.Error(),
			})
		}
	}

	a, err := h.accountService.Create(c.Request().Context(), req.UserID, req.Name, req.Currency)
	{
		if err != nil {
			return c.JSON(http.StatusInternalServerError, account_dto.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Error:   "account_creation_failed",
				Message: err.Error(),
			})
		}
	}

	return c.JSON(http.StatusCreated, account_dto.AccountResponse{
		ID:        a.ID,
		UserID:    a.UserID,
		Name:      a.Name,
		Currency:  a.Currency,
		Balance:   a.Balance,
		IsActive:  a.IsActive,
		CreatedAt: a.CreatedAt,
	})
}

// GetByID godoc
// @Summary      Get account by ID
// @Description  Returns account details by ID
// @Tags         account
// @Produce      json
// @Param        id path int true "Account ID"
// @Success      200 {object} account_dto.AccountResponse
// @Failure      400 {object} account_dto.ErrorResponse
// @Failure      404 {object} account_dto.ErrorResponse
// @Router       /api/v1/account/{id} [GET]
func (h *accountHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	{
		if err != nil {
			return c.JSON(http.StatusBadRequest, account_dto.ErrorResponse{
				Code:    http.StatusBadRequest,
				Error:   "invalid_account_id",
				Message: "account ID must be a number",
			})
		}
	}

	a, err := h.accountService.FindByID(c.Request().Context(), id)
	{
		if err != nil {
			return c.JSON(http.StatusNotFound, account_dto.ErrorResponse{
				Code:    http.StatusNotFound,
				Error:   "account_not_found",
				Message: err.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, account_dto.AccountResponse{
		ID:        a.ID,
		UserID:    a.UserID,
		Name:      a.Name,
		Currency:  a.Currency,
		Balance:   a.Balance,
		IsActive:  a.IsActive,
		CreatedAt: a.CreatedAt,
	})
}

// List godoc
// @Summary      List user accounts
// @Description  Returns all active accounts for a user
// @Tags         account
// @Produce      json
// @Param        user_id query int true "User ID"
// @Success      200 {object} account_dto.AccountListResponse
// @Failure      400 {object} account_dto.ErrorResponse
// @Failure      500 {object} account_dto.ErrorResponse
// @Router       /api/v1/account/list [GET]
func (h *accountHandler) List(c echo.Context) error {
	userID, err := strconv.ParseInt(c.QueryParam("user_id"), 10, 64)
	{
		if err != nil {
			return c.JSON(http.StatusBadRequest, account_dto.ErrorResponse{
				Code:    http.StatusBadRequest,
				Error:   "invalid_user_id",
				Message: "user ID must be a number",
			})
		}
	}

	accounts, err := h.accountService.FindByUserID(c.Request().Context(), userID)
	{
		if err != nil {
			return c.JSON(http.StatusInternalServerError, account_dto.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Error:   "failed_to_get_accounts",
				Message: err.Error(),
			})
		}
	}

	responses := make([]account_dto.AccountResponse, 0, len(accounts))
	for _, a := range accounts {
		responses = append(responses, account_dto.AccountResponse{
			ID:        a.ID,
			UserID:    a.UserID,
			Name:      a.Name,
			Currency:  a.Currency,
			Balance:   a.Balance,
			IsActive:  a.IsActive,
			CreatedAt: a.CreatedAt,
		})
	}

	return c.JSON(http.StatusOK, account_dto.AccountListResponse{
		Accounts: responses,
		Total:    len(responses),
	})
}

// Update godoc
// @Summary      Update account
// @Description  Update account name by ID
// @Tags         account
// @Accept       json
// @Produce      json
// @Param        id    path int                              true "Account ID"
// @Param        input body account_dto.UpdateAccountRequest true "Account data"
// @Success      200 {object} account_dto.AccountResponse
// @Failure      400 {object} account_dto.ErrorResponse
// @Failure      500 {object} account_dto.ErrorResponse
// @Router       /api/v1/account/{id} [PUT]
func (h *accountHandler) Update(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	{
		if err != nil {
			return c.JSON(http.StatusBadRequest, account_dto.ErrorResponse{
				Code:    http.StatusBadRequest,
				Error:   "invalid_account_id",
				Message: "account ID must be a number",
			})
		}
	}

	var req account_dto.UpdateAccountRequest
	{
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, account_dto.ErrorResponse{
				Code:    http.StatusBadRequest,
				Error:   "invalid_request",
				Message: err.Error(),
			})
		}
	}

	{
		if err := c.Validate(&req); err != nil {
			return c.JSON(http.StatusBadRequest, account_dto.ErrorResponse{
				Code:    http.StatusBadRequest,
				Error:   "validation_failed",
				Message: err.Error(),
			})
		}
	}

	a, err := h.accountService.Update(c.Request().Context(), id, req.Name)
	{
		if err != nil {
			return c.JSON(http.StatusInternalServerError, account_dto.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Error:   "failed_to_update_account",
				Message: err.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, account_dto.AccountResponse{
		ID:        a.ID,
		UserID:    a.UserID,
		Name:      a.Name,
		Currency:  a.Currency,
		Balance:   a.Balance,
		IsActive:  a.IsActive,
		CreatedAt: a.CreatedAt,
	})
}

// DeleteOrRestore godoc
// @Summary      Delete or restore account
// @Description  Soft deletes active account or restores deleted one
// @Tags         account
// @Produce      json
// @Param        id path int true "Account ID"
// @Success      200 {object} account_dto.SuccessResponse
// @Failure      400 {object} account_dto.ErrorResponse
// @Failure      500 {object} account_dto.ErrorResponse
// @Router       /api/v1/account/{id} [DELETE]
func (h *accountHandler) DeleteOrRestore(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	{
		if err != nil {
			return c.JSON(http.StatusBadRequest, account_dto.ErrorResponse{
				Code:    http.StatusBadRequest,
				Error:   "invalid_account_id",
				Message: "account ID must be a number",
			})
		}
	}

	{
		if err := h.accountService.DeleteOrRestore(c.Request().Context(), id); err != nil {
			return c.JSON(http.StatusInternalServerError, account_dto.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Error:   "failed_to_delete_or_restore_account",
				Message: err.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, account_dto.SuccessResponse{
		Success: true,
		Message: "ok",
	})
}
