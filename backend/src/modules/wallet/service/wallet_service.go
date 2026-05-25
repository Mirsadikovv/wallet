package service

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"

	crypto_utils "wallet/src/common/utils"
	wallet_model "wallet/src/modules/wallet/model"
)

type WalletService interface {
	Create(ctx context.Context, userID int64, walletType, network string) (*wallet_model.Wallet, error)
	FindByID(ctx context.Context, id int64) (*wallet_model.Wallet, error)
	FindByUserID(ctx context.Context, userID int64) ([]*wallet_model.Wallet, error)
	GetInfo(ctx context.Context, id int64) (*WalletDetailInfo, error)
	GetBalance(ctx context.Context, id int64) (string, error)
	GetTransactions(ctx context.Context, id int64, limit int) ([]*TransactionInfo, error)
	Send(ctx context.Context, id int64, recipient, amount, comment string) (*SendTransactionResult, error)
	Delete(ctx context.Context, id int64) error
}

type walletService struct {
	db            *gorm.DB
	tonService    *TONService
	encryptionKey string
}

func NewWalletService(db *gorm.DB, network string, encryptionKey string) (WalletService, error) {
	tonService, err := NewTONService(network)
	if err != nil {
		return nil, fmt.Errorf("failed to create TON service: %w", err)
	}

	return &walletService{
		db:            db,
		tonService:    tonService,
		encryptionKey: encryptionKey,
	}, nil
}

func (s *walletService) Create(ctx context.Context, userID int64, walletType, network string) (*wallet_model.Wallet, error) {
	seedWords := s.tonService.GenerateWallet()

	walletInfo, err := s.tonService.CreateWalletFromSeed(seedWords, walletType)
	if err != nil {
		return nil, fmt.Errorf("failed to create TON wallet: %w", err)
	}

	encryptedSeed, err := crypto_utils.EncryptSeed(walletInfo.SeedPhrase, s.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt seed: %w", err)
	}

	w := &wallet_model.Wallet{
		UserID:        userID,
		Address:       walletInfo.Address,
		PublicKey:     walletInfo.PublicKey,
		EncryptedSeed: encryptedSeed,
		WalletType:    walletType,
		Network:       network,
		IsActive:      true,
	}

	if err := s.db.WithContext(ctx).Create(w).Error; err != nil {
		return nil, fmt.Errorf("failed to save wallet: %w", err)
	}

	return w, nil
}

func (s *walletService) FindByID(ctx context.Context, id int64) (*wallet_model.Wallet, error) {
	var w wallet_model.Wallet

	if err := s.db.WithContext(ctx).Where("id = ?", id).First(&w).Error; err != nil {
		return nil, fmt.Errorf("wallet not found: %w", err)
	}

	return &w, nil
}

func (s *walletService) FindByUserID(ctx context.Context, userID int64) ([]*wallet_model.Wallet, error) {
	var wallets []*wallet_model.Wallet

	if err := s.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ?", userID, true).
		Order("created_at DESC").
		Find(&wallets).Error; err != nil {
		return nil, fmt.Errorf("failed to get wallets: %w", err)
	}

	return wallets, nil
}

func (s *walletService) GetInfo(ctx context.Context, id int64) (*WalletDetailInfo, error) {
	w, err := s.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	seedPhrase, err := crypto_utils.DecryptSeed(w.EncryptedSeed, s.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt seed: %w", err)
	}

	tonInfo, err := s.tonService.GetWalletInfo(ctx, strings.Split(seedPhrase, " "), w.WalletType)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet info from blockchain: %w", err)
	}

	return &WalletDetailInfo{
		ID:         w.ID,
		Address:    tonInfo.Address,
		Balance:    tonInfo.Balance,
		WalletType: tonInfo.WalletType,
		Network:    w.Network,
		Seqno:      tonInfo.Seqno,
		IsActive:   w.IsActive,
		CreatedAt:  w.CreatedAt,
	}, nil
}

func (s *walletService) GetBalance(ctx context.Context, id int64) (string, error) {
	w, err := s.FindByID(ctx, id)
	if err != nil {
		return "", err
	}

	seedPhrase, err := crypto_utils.DecryptSeed(w.EncryptedSeed, s.encryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt seed: %w", err)
	}

	balance, err := s.tonService.GetBalance(ctx, strings.Split(seedPhrase, " "), w.WalletType)
	if err != nil {
		return "", fmt.Errorf("failed to get balance: %w", err)
	}

	return balance, nil
}

func (s *walletService) GetTransactions(ctx context.Context, id int64, limit int) ([]*TransactionInfo, error) {
	w, err := s.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	seedPhrase, err := crypto_utils.DecryptSeed(w.EncryptedSeed, s.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt seed: %w", err)
	}

	transactions, err := s.tonService.GetTransactions(ctx, strings.Split(seedPhrase, " "), w.WalletType, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	return transactions, nil
}

func (s *walletService) Send(ctx context.Context, id int64, recipient, amount, comment string) (*SendTransactionResult, error) {
	w, err := s.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	seedPhrase, err := crypto_utils.DecryptSeed(w.EncryptedSeed, s.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt seed: %w", err)
	}

	result, err := s.tonService.SendTransaction(ctx, strings.Split(seedPhrase, " "), w.WalletType, recipient, amount, comment)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	return result, nil
}

func (s *walletService) Delete(ctx context.Context, id int64) error {
	if err := s.db.WithContext(ctx).
		Model(&wallet_model.Wallet{}).
		Where("id = ?", id).
		Update("is_active", false).Error; err != nil {
		return fmt.Errorf("failed to delete wallet: %w", err)
	}

	return nil
}
