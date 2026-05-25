package service

import "time"

type WalletInfo struct {
	Address    string
	PublicKey  string
	SeedPhrase string
	WalletType string
}

type TONWalletInfo struct {
	Address    string
	Balance    string
	WalletType string
	Seqno      int64
}

type WalletDetailInfo struct {
	ID         int64
	Address    string
	Balance    string
	WalletType string
	Network    string
	Seqno      int64
	IsActive   bool
	CreatedAt  time.Time
}

type TransactionInfo struct {
	Hash      string
	Lt        uint64
	Timestamp int64
	Type      string
	Amount    string
	Fee       string
	From      string
	To        string
	Comment   string
	Success   bool
}

type SendTransactionResult struct {
	Hash      string
	Lt        uint64
	Address   string
	Amount    string
	Fee       string
	Recipient string
	Comment   string
}
