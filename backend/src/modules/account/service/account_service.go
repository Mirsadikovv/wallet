package service

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	account_model "wallet/src/modules/account/model"
)

type AccountService interface {
	Create(ctx context.Context, userID int64, name, currency string) (*account_model.Account, error)
	FindByID(ctx context.Context, id int64) (*account_model.Account, error)
	FindByUserID(ctx context.Context, userID int64) ([]*account_model.Account, error)
	Update(ctx context.Context, id int64, name string) (*account_model.Account, error)
	DeleteOrRestore(ctx context.Context, id int64) error
}

type accountService struct {
	db *gorm.DB
}

func NewAccountService(db *gorm.DB) AccountService {
	return &accountService{db: db}
}

func (s *accountService) Create(ctx context.Context, userID int64, name, currency string) (*account_model.Account, error) {
	a := &account_model.Account{
		UserID:   userID,
		Name:     name,
		Currency: currency,
		Balance:  "0",
		IsActive: true,
	}

	if err := s.db.WithContext(ctx).Create(a).Error; err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return a, nil
}

func (s *accountService) FindByID(ctx context.Context, id int64) (*account_model.Account, error) {
	var a account_model.Account

	if err := s.db.WithContext(ctx).
		Select("id", "user_id", "name", "currency", "balance", "is_active", "created_at", "updated_at").
		Where("id = ?", id).
		First(&a).Error; err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}

	return &a, nil
}

func (s *accountService) FindByUserID(ctx context.Context, userID int64) ([]*account_model.Account, error) {
	var accounts []*account_model.Account

	if err := s.db.WithContext(ctx).
		Select("id", "user_id", "name", "currency", "balance", "is_active", "created_at").
		Where("user_id = ? AND is_active = ?", userID, true).
		Order("created_at DESC").
		Find(&accounts).Error; err != nil {
		return nil, fmt.Errorf("failed to get accounts: %w", err)
	}

	return accounts, nil
}

func (s *accountService) Update(ctx context.Context, id int64, name string) (*account_model.Account, error) {
	if err := s.db.WithContext(ctx).
		Model(&account_model.Account{}).
		Where("id = ?", id).
		Update("name", name).Error; err != nil {
		return nil, fmt.Errorf("failed to update account: %w", err)
	}

	return s.FindByID(ctx, id)
}

func (s *accountService) DeleteOrRestore(ctx context.Context, id int64) error {
	var a account_model.Account

	if err := s.db.WithContext(ctx).Unscoped().Where("id = ?", id).First(&a).Error; err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	if a.DeletedAt.Valid {
		return s.db.WithContext(ctx).
			Model(&account_model.Account{}).
			Where("id = ?", id).
			Update("deleted_at", nil).Error
	}

	return s.db.WithContext(ctx).Where("id = ?", id).Delete(&account_model.Account{}).Error
}
