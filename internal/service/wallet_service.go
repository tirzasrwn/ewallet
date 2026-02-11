package service

import (
	"errors"
	"ewallet/internal/models"
	"ewallet/internal/repository"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletService interface {
	GetBalance(userID uint) (*models.Wallet, error)
	TopUp(userID uint, amount float64) (*models.Wallet, error)
}

type walletService struct {
	walletRepo      repository.WalletRepository
	transactionRepo repository.TransactionRepository
	db              *gorm.DB
}

func NewWalletService(
	walletRepo repository.WalletRepository,
	transactionRepo repository.TransactionRepository,
	db *gorm.DB,
) WalletService {
	return &walletService{
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
		db:              db,
	}
}

func (s *walletService) GetBalance(userID uint) (*models.Wallet, error) {
	wallet, err := s.walletRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (s *walletService) TopUp(userID uint, amount float64) (*models.Wallet, error) {
	// Validate amount
	if amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	var wallet *models.Wallet

	// Use transaction to ensure atomicity
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Get wallet with lock
		var w models.Wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", userID).
			First(&w).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("wallet not found")
			}
			return err
		}

		// Update balance
		newBalance := w.Balance + amount
		if err := s.walletRepo.UpdateBalanceWithLock(tx, w.ID, newBalance); err != nil {
			return err
		}

		// Create transaction record
		transaction := &models.Transaction{
			ReceiverID: userID,
			Amount:     amount,
			Type:       models.TransactionTypeTopUp,
			Status:     models.TransactionStatusSuccess,
		}

		if err := s.transactionRepo.Create(tx, transaction); err != nil {
			return err
		}

		// Get updated wallet
		w.Balance = newBalance
		wallet = &w
		return nil
	})

	if err != nil {
		return nil, err
	}

	return wallet, nil
}
