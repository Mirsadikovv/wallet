package service

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

type TONService struct {
	client *liteclient.ConnectionPool
	api    ton.APIClientWrapped
	config *liteclient.GlobalConfig
}

func NewTONService(network string) (*TONService, error) {
	client := liteclient.NewConnectionPool()

	configURL := "https://ton.org/global-config.json" // mainnet
	if network == "testnet" {
		configURL = "https://ton.org/testnet-global.config.json"
	}

	cfg, err := liteclient.GetConfigFromUrl(context.Background(), configURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	err = client.AddConnectionsFromConfig(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	api := ton.NewAPIClient(client, ton.ProofCheckPolicyFast).WithRetry()
	api.SetTrustedBlockFromConfig(cfg)

	return &TONService{
		client: client,
		api:    api,
		config: cfg,
	}, nil
}

func (s *TONService) GenerateWallet() []string {
	seed := wallet.NewSeed()
	return seed
}

func (s *TONService) CreateWalletFromSeed(seedWords []string, walletType string) (*WalletInfo, error) {
	config := wallet.ConfigV5R1Final{
		NetworkGlobalID: wallet.MainnetGlobalID,
	}

	w, err := wallet.FromSeed(s.api, seedWords, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet from seed: %w", err)
	}

	address := w.WalletAddress()

	return &WalletInfo{
		Address:    address.String(),
		SeedPhrase: strings.Join(seedWords, " "),
		WalletType: walletType,
	}, nil
}

func (s *TONService) GetBalance(ctx context.Context, seedWords []string, walletType string) (string, error) {
	config := wallet.ConfigV5R1Final{
		NetworkGlobalID: wallet.MainnetGlobalID,
	}

	w, err := wallet.FromSeed(s.api, seedWords, config)
	if err != nil {
		return "", fmt.Errorf("failed to create wallet: %w", err)
	}

	block, err := s.api.CurrentMasterchainInfo(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get masterchain info: %w", err)
	}

	balance, err := w.GetBalance(ctx, block)
	if err != nil {
		return "", fmt.Errorf("failed to get balance: %w", err)
	}

	return balance.String(), nil
}

func (s *TONService) GetWalletInfo(ctx context.Context, seedWords []string, walletType string) (*WalletDetailInfo, error) {
	config := wallet.ConfigV5R1Final{
		NetworkGlobalID: wallet.MainnetGlobalID,
	}

	w, err := wallet.FromSeed(s.api, seedWords, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	address := w.WalletAddress()

	block, err := s.api.CurrentMasterchainInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get masterchain info: %w", err)
	}

	balance, err := w.GetBalance(ctx, block)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	var seqno uint32 = 0

	return &WalletDetailInfo{
		Address:    address.String(),
		Balance:    balance.String(),
		WalletType: walletType,
		Seqno:      int64(seqno),
	}, nil
}

func EncryptSeed(seedPhrase, encryptionKey string) (string, error) {
	block, err := aes.NewCipher([]byte(encryptionKey))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(seedPhrase), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func DecryptSeed(encryptedSeed, encryptionKey string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedSeed)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(encryptionKey))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

type WalletInfo struct {
	Address    string `json:"address"`
	SeedPhrase string `json:"seed_phrase"`
	WalletType string `json:"wallet_type"`
}

type WalletDetailInfo struct {
	Address    string `json:"address"`
	Balance    string `json:"balance"`
	WalletType string `json:"wallet_type"`
	Seqno      int64  `json:"seqno"`
}

func TONAmount(amount string) (tlb.Coins, error) {
	return tlb.FromTON(amount)
}
