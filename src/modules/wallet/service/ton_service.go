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

type TransactionInfo struct {
	Hash      string `json:"hash"`
	Lt        uint64 `json:"lt"`
	Timestamp int64  `json:"timestamp"`
	Type      string `json:"type"`      // "in" или "out"
	Amount    string `json:"amount"`    // в TON
	Fee       string `json:"fee"`       // в TON
	From      string `json:"from"`      // адрес отправителя
	To        string `json:"to"`        // адрес получателя
	Comment   string `json:"comment"`   // комментарий к транзакции
	Success   bool   `json:"success"`   // успешна ли транзакция
}

func (s *TONService) GetTransactions(ctx context.Context, seedWords []string, walletType string, limit int) ([]*TransactionInfo, error) {
	config := wallet.ConfigV5R1Final{
		NetworkGlobalID: wallet.MainnetGlobalID,
	}

	w, err := wallet.FromSeed(s.api, seedWords, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	address := w.WalletAddress()

	// Получаем список транзакций (lt=0, hash=nil означает "с последней транзакции")
	txList, err := s.api.ListTransactions(ctx, address, uint32(limit), 0, nil)
	if err != nil {
		// Если транзакций нет, возвращаем пустой массив вместо ошибки
		if strings.Contains(err.Error(), "no transactions were found") {
			return []*TransactionInfo{}, nil
		}
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	var transactions []*TransactionInfo
	for _, tx := range txList {
		txInfo := &TransactionInfo{
			Hash:      base64.StdEncoding.EncodeToString(tx.Hash),
			Lt:        tx.LT,
			Timestamp: int64(tx.Now),
			Success:   true, // По умолчанию считаем успешной
		}

		// Обрабатываем комиссию
		if tx.TotalFees.Coins.Nano() != nil {
			txInfo.Fee = tx.TotalFees.Coins.TON()
		}

		// Обрабатываем входящие сообщения
		if tx.IO.In != nil && tx.IO.In.MsgType == tlb.MsgTypeInternal {
			intMsg := tx.IO.In.AsInternal()
			if intMsg != nil {
				txInfo.Type = "in"
				txInfo.Amount = intMsg.Amount.TON()
				txInfo.From = intMsg.SrcAddr.String()
				txInfo.To = address.String()

				// Пытаемся извлечь комментарий
				if intMsg.Body != nil {
					payload := intMsg.Body.BeginParse()
					if op, err := payload.LoadUInt(32); err == nil && op == 0 {
						if comment, err := payload.LoadStringSnake(); err == nil {
							txInfo.Comment = comment
						}
					}
				}
			}
		}

		// Обрабатываем исходящие сообщения
		if tx.IO.Out != nil {
			list, err := tx.IO.Out.ToSlice()
			if err == nil {
				for _, msg := range list {
					if msg.MsgType == tlb.MsgTypeInternal {
						intMsg := msg.AsInternal()
						if intMsg != nil {
							txInfo.Type = "out"
							txInfo.Amount = intMsg.Amount.TON()
							txInfo.From = address.String()
							txInfo.To = intMsg.DstAddr.String()

							// Пытаемся извлечь комментарий
							if intMsg.Body != nil {
								payload := intMsg.Body.BeginParse()
								if op, err := payload.LoadUInt(32); err == nil && op == 0 {
									if comment, err := payload.LoadStringSnake(); err == nil {
										txInfo.Comment = comment
									}
								}
							}
							break // Берем только первое исходящее сообщение
						}
					}
				}
			}
		}

		transactions = append(transactions, txInfo)
	}

	return transactions, nil
}

func TONAmount(amount string) (tlb.Coins, error) {
	return tlb.FromTON(amount)
}
