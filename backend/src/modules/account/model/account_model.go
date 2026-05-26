package model

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	ID        int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int64          `json:"user_id" gorm:"not null;index"`
	Name      string         `json:"name" gorm:"type:varchar(255);not null"`
	Currency  string         `json:"currency" gorm:"type:varchar(10);not null;default:'USD'"`
	Balance   string         `json:"balance" gorm:"type:varchar(50);not null;default:'0'"`
	IsActive  bool           `json:"is_active" gorm:"default:true;not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
} // @name Account
