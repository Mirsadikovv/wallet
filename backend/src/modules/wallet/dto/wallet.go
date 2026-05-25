package dto

import "time"

type CreateWalletRequest struct {
	UserID     int64  `json:"user_id" validate:"required"`
	WalletType string `json:"wallet_type" validate:"required,oneof=V5R1Final V4R2"`
	Network    string `json:"network" validate:"required,oneof=mainnet testnet"`
} // @Name CreateWalletRequest

type CreateWalletResponse struct {
	ID         int64     `json:"id"`
	Address    string    `json:"address"`
	WalletType string    `json:"wallet_type"`
	Network    string    `json:"network"`
	CreatedAt  time.Time `json:"created_at"`
} // @Name CreateWalletResponse

type GetWalletInfoResponse struct {
	ID         int64     `json:"id"`
	Address    string    `json:"address"`
	Balance    string    `json:"balance"`
	WalletType string    `json:"wallet_type"`
	Network    string    `json:"network"`
	Seqno      int64     `json:"seqno"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
} // @Name GetWalletInfoResponse

type GetBalanceResponse struct {
	Address string `json:"address"`
	Balance string `json:"balance"`
} // @Name GetBalanceResponse

type WalletSummary struct {
	ID         int64     `json:"id"`
	Address    string    `json:"address"`
	WalletType string    `json:"wallet_type"`
	Network    string    `json:"network"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
} // @Name WalletSummary

type ListWalletsResponse struct {
	Wallets []WalletSummary `json:"wallets"`
	Total   int             `json:"total"`
} // @Name ListWalletsResponse

type ErrorResponse struct {
	Code    int    `json:"code"`
	Error   string `json:"error"`
	Message string `json:"message"`
} // @Name ErrorResponse

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
} // @Name SuccessResponse

type GetTransactionsRequest struct {
	Limit int `query:"limit" validate:"omitempty,min=1,max=100"`
} // @Name GetTransactionsRequest

type TransactionDTO struct {
	Hash      string `json:"hash"`
	Lt        uint64 `json:"lt"`
	Timestamp int64  `json:"timestamp"`
	Type      string `json:"type"`
	Amount    string `json:"amount"`
	Fee       string `json:"fee"`
	From      string `json:"from"`
	To        string `json:"to"`
	Comment   string `json:"comment,omitempty"`
	Success   bool   `json:"success"`
} // @Name TransactionDTO

type GetTransactionsResponse struct {
	WalletID     int64             `json:"wallet_id"`
	Address      string            `json:"address"`
	Transactions []*TransactionDTO `json:"transactions"`
	Total        int               `json:"total"`
} // @Name GetTransactionsResponse

type SendCoinsRequest struct {
	Recipient string `json:"recipient" validate:"required"`
	Amount    string `json:"amount" validate:"required"`
	Comment   string `json:"comment,omitempty"`
} // @Name SendCoinsRequest

type SendCoinsResponse struct {
	Hash      string `json:"hash"`
	Lt        uint64 `json:"lt"`
	Address   string `json:"address"`
	Amount    string `json:"amount"`
	Fee       string `json:"fee"`
	Recipient string `json:"recipient"`
	Comment   string `json:"comment,omitempty"`
} // @Name SendCoinsResponse
