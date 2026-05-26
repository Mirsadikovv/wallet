package dto

import "time"

type CreateAccountRequest struct {
	UserID   int64  `json:"user_id" validate:"required"`
	Name     string `json:"name" validate:"required,min=1,max=255"`
	Currency string `json:"currency" validate:"required,oneof=USD EUR RUB"`
} // @Name CreateAccountRequest

type UpdateAccountRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
} // @Name UpdateAccountRequest

type AccountResponse struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
	Currency  string    `json:"currency"`
	Balance   string    `json:"balance"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
} // @Name AccountResponse

type AccountListResponse struct {
	Accounts []AccountResponse `json:"accounts"`
	Total    int               `json:"total"`
} // @Name AccountListResponse

type ErrorResponse struct {
	Code    int    `json:"code"`
	Error   string `json:"error"`
	Message string `json:"message"`
} // @Name ErrorResponse

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
} // @Name SuccessResponse
