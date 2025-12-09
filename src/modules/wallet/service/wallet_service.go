package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/uptrace/bun"
	"wallet_test/src/modules/wallet/model"
)

type WalletService struct {
	db            *bun.DB
	tonService    *TONService
	encryptionKey string
}

func NewWalletService(db *bun.DB, network string, encryptionKey string) (*WalletService, error) {
	tonService, err := NewTONService(network)
	if err != nil {
		return nil, fmt.Errorf("failed to create TON service: %w", err)
	}

	return &WalletService{
		db:            db,
		tonService:    tonService,
		encryptionKey: encryptionKey,
	}, nil
}

func (s *WalletService) CreateWallet(ctx context.Context, userID int64, walletType, network string) (*model.Wallet, error) {
	seedWords := s.tonService.GenerateWallet()

	walletInfo, err := s.tonService.CreateWalletFromSeed(seedWords, walletType)
	if err != nil {
		return nil, fmt.Errorf("failed to create TON wallet: %w", err)
	}

	encryptedSeed, err := EncryptSeed(walletInfo.SeedPhrase, s.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt seed: %w", err)
	}

	wallet := &model.Wallet{
		UserID:        userID,
		Address:       walletInfo.Address,
		EncryptedSeed: encryptedSeed,
		WalletType:    walletType,
		Network:       network,
		IsActive:      true,
	}

	_, err = s.db.NewInsert().Model(wallet).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to save wallet: %w", err)
	}

	return wallet, nil
}

func (s *WalletService) GetWalletByID(ctx context.Context, walletID int64) (*model.Wallet, error) {
	wallet := &model.Wallet{}
	err := s.db.NewSelect().
		Model(wallet).
		Where("id = ?", walletID).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	return wallet, nil
}

func (s *WalletService) GetWalletByAddress(ctx context.Context, address string) (*model.Wallet, error) {
	wallet := &model.Wallet{}
	err := s.db.NewSelect().
		Model(wallet).
		Where("address = ?", address).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	return wallet, nil
}

func (s *WalletService) GetUserWallets(ctx context.Context, userID int64) ([]*model.Wallet, error) {
	var wallets []*model.Wallet
	err := s.db.NewSelect().
		Model(&wallets).
		Where("user_id = ?", userID).
		Where("is_active = ?", true).
		Order("created_at DESC").
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get wallets: %w", err)
	}

	return wallets, nil
}

func (s *WalletService) GetWalletInfo(ctx context.Context, walletID int64) (*WalletDetailInfo, error) {
	wallet, err := s.GetWalletByID(ctx, walletID)
	if err != nil {
		return nil, err
	}

	seedPhrase, err := DecryptSeed(wallet.EncryptedSeed, s.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt seed: %w", err)
	}

	seedWords := strings.Split(seedPhrase, " ")

	info, err := s.tonService.GetWalletInfo(ctx, seedWords, wallet.WalletType)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet info from blockchain: %w", err)
	}

	return info, nil
}

func (s *WalletService) GetBalance(ctx context.Context, walletID int64) (string, error) {
	wallet, err := s.GetWalletByID(ctx, walletID)
	if err != nil {
		return "", err
	}

	seedPhrase, err := DecryptSeed(wallet.EncryptedSeed, s.encryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt seed: %w", err)
	}

	seedWords := strings.Split(seedPhrase, " ")

	balance, err := s.tonService.GetBalance(ctx, seedWords, wallet.WalletType)
	if err != nil {
		return "", fmt.Errorf("failed to get balance: %w", err)
	}

	return balance, nil
}

func (s *WalletService) DeleteWallet(ctx context.Context, walletID int64) error {
	_, err := s.db.NewUpdate().
		Model((*model.Wallet)(nil)).
		Set("is_active = ?", false).
		Where("id = ?", walletID).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to delete wallet: %w", err)
	}

	return nil
}

func (s *WalletService) GetTransactions(ctx context.Context, walletID int64, limit int) ([]*TransactionInfo, error) {
	wallet, err := s.GetWalletByID(ctx, walletID)
	if err != nil {
		return nil, err
	}

	seedPhrase, err := DecryptSeed(wallet.EncryptedSeed, s.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt seed: %w", err)
	}

	seedWords := strings.Split(seedPhrase, " ")

	transactions, err := s.tonService.GetTransactions(ctx, seedWords, wallet.WalletType, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	return transactions, nil
}

func (s *WalletService) SendCoins(ctx context.Context, walletID int64, recipient, amount, comment string) (*SendTransactionResult, error) {
	wallet, err := s.GetWalletByID(ctx, walletID)
	if err != nil {
		return nil, err
	}

	seedPhrase, err := DecryptSeed(wallet.EncryptedSeed, s.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt seed: %w", err)
	}

	seedWords := strings.Split(seedPhrase, " ")

	result, err := s.tonService.SendTransaction(ctx, seedWords, wallet.WalletType, recipient, amount, comment)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	return result, nil
}
