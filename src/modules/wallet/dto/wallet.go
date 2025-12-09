package dto

type CreateWalletRequest struct {
	UserID     int64  `json:"user_id" binding:"required"`
	WalletType string `json:"wallet_type" binding:"required,oneof=V5R1Final V4R2"`
	Network    string `json:"network" binding:"required,oneof=mainnet testnet"`
}

type CreateWalletResponse struct {
	ID         int64  `json:"id"`
	Address    string `json:"address"`
	WalletType string `json:"wallet_type"`
	Network    string `json:"network"`
	CreatedAt  string `json:"created_at"`
}

type GetWalletInfoResponse struct {
	ID         int64  `json:"id"`
	Address    string `json:"address"`
	Balance    string `json:"balance"`
	WalletType string `json:"wallet_type"`
	Network    string `json:"network"`
	Seqno      int64  `json:"seqno"`
	IsActive   bool   `json:"is_active"`
	CreatedAt  string `json:"created_at"`
}

type GetBalanceResponse struct {
	Address string `json:"address"`
	Balance string `json:"balance"`
}

type ListWalletsResponse struct {
	Wallets []WalletSummary `json:"wallets"`
	Total   int             `json:"total"`
}

type WalletSummary struct {
	ID         int64  `json:"id"`
	Address    string `json:"address"`
	WalletType string `json:"wallet_type"`
	Network    string `json:"network"`
	IsActive   bool   `json:"is_active"`
	CreatedAt  string `json:"created_at"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}
