package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Wallet struct {
	bun.BaseModel `bun:"table:wallets,alias:w"`

	ID            int64     `bun:"id,pk,autoincrement" json:"id"`
	UserID        int64     `bun:"user_id,notnull" json:"user_id"`
	Address       string    `bun:"address,notnull,unique" json:"address"`
	PublicKey     string    `bun:"public_key,notnull" json:"public_key"`
	EncryptedSeed string    `bun:"encrypted_seed,notnull" json:"-"`
	WalletType    string    `bun:"wallet_type,notnull" json:"wallet_type"` // V5R1Final, V4R2, etc
	Network       string    `bun:"network,notnull" json:"network"`         // mainnet, testnet
	IsActive      bool      `bun:"is_active,notnull,default:true" json:"is_active"`
	CreatedAt     time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt     time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`
	User          *User     `bun:"rel:belongs-to,join:user_id=id" json:"user,omitempty"`
}

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	Username  string    `bun:"username,unique,notnull" json:"username"`
	Email     string    `bun:"email,unique,notnull" json:"email"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`
	Wallets   []*Wallet `bun:"rel:has-many,join:id=user_id" json:"wallets,omitempty"`
}

// Transaction история транзакций
type Transaction struct {
	bun.BaseModel `bun:"table:transactions,alias:t"`

	ID          int64     `bun:"id,pk,autoincrement" json:"id"`
	WalletID    int64     `bun:"wallet_id,notnull" json:"wallet_id"`
	TxHash      string    `bun:"tx_hash,unique,notnull" json:"tx_hash"`
	FromAddress string    `bun:"from_address,notnull" json:"from_address"`
	ToAddress   string    `bun:"to_address,notnull" json:"to_address"`
	Amount      string    `bun:"amount,notnull" json:"amount"` // храним как string для точности
	Fee         string    `bun:"fee" json:"fee"`
	Status      string    `bun:"status,notnull" json:"status"` // pending, confirmed, failed
	BlockNumber int64     `bun:"block_number" json:"block_number"`
	Comment     string    `bun:"comment" json:"comment"`
	CreatedAt   time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`

	// Relations
	Wallet *Wallet `bun:"rel:belongs-to,join:wallet_id=id" json:"wallet,omitempty"`
}
