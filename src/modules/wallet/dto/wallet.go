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

type GetTransactionsRequest struct {
	Limit int `form:"limit" json:"limit" binding:"omitempty,min=1,max=100"`
}

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
}

type GetTransactionsResponse struct {
	WalletID     int64             `json:"wallet_id"`
	Address      string            `json:"address"`
	Transactions []*TransactionDTO `json:"transactions"`
	Total        int               `json:"total"`
}

type SendCoinsRequest struct {
	Recipient string `json:"recipient" binding:"required"` // Адрес получателя
	Amount    string `json:"amount" binding:"required"`    // Сумма в TON (например "1.5")
	Comment   string `json:"comment,omitempty"`            // Комментарий к транзакции
}

type SendCoinsResponse struct {
	Hash      string `json:"hash"`       // Хеш транзакции
	Lt        uint64 `json:"lt"`         // Logical time
	Address   string `json:"address"`    // Адрес отправителя
	Amount    string `json:"amount"`     // Отправленная сумма
	Fee       string `json:"fee"`        // Комиссия
	Recipient string `json:"recipient"`  // Адрес получателя
	Comment   string `json:"comment,omitempty"` // Комментарий
}
