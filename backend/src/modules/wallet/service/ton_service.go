package service

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

type TONService struct {
	client          *liteclient.ConnectionPool
	api             ton.APIClientWrapped
	config          *liteclient.GlobalConfig
	networkGlobalID int32
}

func NewTONService(network string) (*TONService, error) {
	client := liteclient.NewConnectionPool()

	configURL := "https://ton.org/global-config.json"
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

	networkGlobalID := int32(wallet.MainnetGlobalID)
	if network == "testnet" {
		networkGlobalID = int32(wallet.TestnetGlobalID)
	}

	return &TONService{
		client:          client,
		api:             api,
		config:          cfg,
		networkGlobalID: networkGlobalID,
	}, nil
}

func (s *TONService) GenerateWallet() []string {
	return wallet.NewSeed()
}

func (s *TONService) walletFromSeed(seedWords []string, walletType string) (*wallet.Wallet, error) {
	switch walletType {
	case "V4R2":
		return wallet.FromSeed(s.api, seedWords, wallet.V4R2)
	default:
		return wallet.FromSeed(s.api, seedWords, wallet.ConfigV5R1Final{
			NetworkGlobalID: s.networkGlobalID,
		})
	}
}

func (s *TONService) CreateWalletFromSeed(seedWords []string, walletType string) (*WalletInfo, error) {
	w, err := s.walletFromSeed(seedWords, walletType)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet from seed: %w", err)
	}

	pubKey := ed25519.PublicKey(w.PrivateKey().Public().(ed25519.PublicKey))

	return &WalletInfo{
		Address:    w.WalletAddress().String(),
		PublicKey:  hex.EncodeToString(pubKey),
		SeedPhrase: strings.Join(seedWords, " "),
		WalletType: walletType,
	}, nil
}

func (s *TONService) GetBalance(ctx context.Context, seedWords []string, walletType string) (string, error) {
	w, err := s.walletFromSeed(seedWords, walletType)
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

func (s *TONService) GetWalletInfo(ctx context.Context, seedWords []string, walletType string) (*TONWalletInfo, error) {
	w, err := s.walletFromSeed(seedWords, walletType)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	block, err := s.api.CurrentMasterchainInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get masterchain info: %w", err)
	}

	balance, err := w.GetBalance(ctx, block)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	return &TONWalletInfo{
		Address:    w.WalletAddress().String(),
		Balance:    balance.String(),
		WalletType: walletType,
		Seqno:      0,
	}, nil
}

func (s *TONService) GetTransactions(ctx context.Context, seedWords []string, walletType string, limit int) ([]*TransactionInfo, error) {
	w, err := s.walletFromSeed(seedWords, walletType)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	addr := w.WalletAddress()

	txList, err := s.api.ListTransactions(ctx, addr, uint32(limit), 0, nil)
	if err != nil {
		if strings.Contains(err.Error(), "no transactions were found") {
			return []*TransactionInfo{}, nil
		}
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	transactions := make([]*TransactionInfo, 0, len(txList))
	for _, tx := range txList {
		success := true
		if desc, ok := tx.Description.(tlb.TransactionDescriptionOrdinary); ok {
			success = !desc.Aborted
		}

		txInfo := &TransactionInfo{
			Hash:      base64.StdEncoding.EncodeToString(tx.Hash),
			Lt:        tx.LT,
			Timestamp: int64(tx.Now),
			Success:   success,
		}

		if tx.TotalFees.Coins.Nano() != nil {
			txInfo.Fee = tx.TotalFees.Coins.TON()
		}

		if tx.IO.In != nil && tx.IO.In.MsgType == tlb.MsgTypeInternal {
			intMsg := tx.IO.In.AsInternal()
			if intMsg != nil {
				txInfo.Type = "in"
				txInfo.Amount = intMsg.Amount.TON()
				if intMsg.SrcAddr != nil {
					txInfo.From = intMsg.SrcAddr.String()
				}
				txInfo.To = addr.String()

				if intMsg.Body != nil {
					if payload, err := intMsg.Body.BeginParse(); err == nil {
						if op, err := payload.LoadUInt(32); err == nil && op == 0 {
							if comment, err := payload.LoadStringSnake(); err == nil {
								txInfo.Comment = comment
							}
						}
					}
				}
			}
		}

		if tx.IO.Out != nil {
			list, err := tx.IO.Out.ToSlice()
			if err == nil {
				totalNano := new(big.Int)
				for _, msg := range list {
					if msg.MsgType == tlb.MsgTypeInternal {
						intMsg := msg.AsInternal()
						if intMsg != nil {
							if txInfo.Type == "" {
								txInfo.Type = "out"
								txInfo.From = addr.String()
								if intMsg.DstAddr != nil {
									txInfo.To = intMsg.DstAddr.String()
								}
								if intMsg.Body != nil {
									if payload, err := intMsg.Body.BeginParse(); err == nil {
										if op, err := payload.LoadUInt(32); err == nil && op == 0 {
											if comment, err := payload.LoadStringSnake(); err == nil {
												txInfo.Comment = comment
											}
										}
									}
								}
							}
							if intMsg.Amount.Nano() != nil {
								totalNano.Add(totalNano, intMsg.Amount.Nano())
							}
						}
					}
				}
				if txInfo.Type == "out" {
					txInfo.Amount = nanoToTON(totalNano)
				}
			}
		}

		transactions = append(transactions, txInfo)
	}

	return transactions, nil
}

func nanoToTON(nano *big.Int) string {
	if nano == nil || nano.Sign() == 0 {
		return "0"
	}
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil)
	quotient := new(big.Int)
	remainder := new(big.Int)
	quotient.DivMod(nano, divisor, remainder)
	if remainder.Sign() == 0 {
		return quotient.String()
	}
	fracStr := strings.TrimRight(fmt.Sprintf("%09d", remainder.Int64()), "0")
	return quotient.String() + "." + fracStr
}

func (s *TONService) SendTransaction(ctx context.Context, seedWords []string, walletType, recipient, amount, comment string) (*SendTransactionResult, error) {
	w, err := s.walletFromSeed(seedWords, walletType)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	addr, err := address.ParseAddr(recipient)
	if err != nil {
		return nil, fmt.Errorf("invalid recipient address: %w", err)
	}

	coins, err := tlb.FromTON(amount)
	if err != nil {
		return nil, fmt.Errorf("invalid amount: %w", err)
	}

	var body *cell.Cell
	if comment != "" {
		body, err = wallet.CreateCommentCell(comment)
		if err != nil {
			return nil, fmt.Errorf("failed to create comment: %w", err)
		}
	}

	tx, _, err := w.SendWaitTransaction(ctx, &wallet.Message{
		Mode: 3,
		InternalMessage: &tlb.InternalMessage{
			IHRDisabled: true,
			Bounce:      addr.IsBounceable(),
			DstAddr:     addr,
			Amount:      coins,
			Body:        body,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	fee := "0"
	if tx.TotalFees.Coins.Nano() != nil {
		fee = tx.TotalFees.Coins.TON()
	}

	return &SendTransactionResult{
		Hash:      base64.StdEncoding.EncodeToString(tx.Hash),
		Lt:        tx.LT,
		Address:   w.WalletAddress().String(),
		Amount:    amount,
		Fee:       fee,
		Recipient: recipient,
		Comment:   comment,
	}, nil
}
