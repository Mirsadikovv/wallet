package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID               int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	TelegramID       *int64         `json:"telegram_id" gorm:"uniqueIndex"`
	TelegramUsername string         `json:"telegram_username" gorm:"type:varchar(255)"`
	Username         string         `json:"username" gorm:"type:varchar(255);uniqueIndex;not null"`
	Email            string         `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	CreatedAt        time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	Wallets          []*Wallet      `json:"wallets,omitempty" gorm:"foreignKey:UserID"`
} // @name User

type Wallet struct {
	ID            int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID        int64          `json:"user_id" gorm:"not null;index"`
	Address       string         `json:"address" gorm:"type:varchar(255);uniqueIndex;not null"`
	PublicKey     string         `json:"public_key" gorm:"type:varchar(255);not null"`
	EncryptedSeed string         `json:"-" gorm:"type:text;not null"`
	WalletType    string         `json:"wallet_type" gorm:"type:varchar(50);not null"`
	Network       string         `json:"network" gorm:"type:varchar(20);not null"`
	IsActive      bool           `json:"is_active" gorm:"default:true;not null"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	User          *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
} // @name Wallet

type Transaction struct {
	ID          int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	WalletID    int64          `json:"wallet_id" gorm:"not null;index"`
	TxHash      string         `json:"tx_hash" gorm:"type:varchar(255);uniqueIndex;not null"`
	FromAddress string         `json:"from_address" gorm:"type:varchar(255);not null"`
	ToAddress   string         `json:"to_address" gorm:"type:varchar(255);not null"`
	Amount      string         `json:"amount" gorm:"type:varchar(50);not null"`
	Fee         string         `json:"fee" gorm:"type:varchar(50)"`
	Status      string         `json:"status" gorm:"type:varchar(20);not null"`
	BlockNumber int64          `json:"block_number"`
	Comment     string         `json:"comment" gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	Wallet      *Wallet        `json:"wallet,omitempty" gorm:"foreignKey:WalletID"`
} // @name Transaction
