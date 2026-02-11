package service

import (
	"errors"
	"ewallet/internal/models"
	"ewallet/internal/repository"

	"bytes"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransactionService interface {
	Transfer(senderID, receiverID uuid.UUID, amount float64) (*models.Transaction, error)
	GetHistory(userID uuid.UUID, limit int) ([]models.Transaction, error)
}

type transactionService struct {
	walletRepo      repository.WalletRepository
	transactionRepo repository.TransactionRepository
	userRepo        repository.UserRepository
	db              *gorm.DB
}

func NewTransactionService(
	walletRepo repository.WalletRepository,
	transactionRepo repository.TransactionRepository,
	userRepo repository.UserRepository,
	db *gorm.DB,
) TransactionService {
	return &transactionService{
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
		db:              db,
	}
}

func (s *transactionService) Transfer(senderID uuid.UUID, receiverID uuid.UUID, amount float64) (*models.Transaction, error) {
	// Validate amount
	if amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	// Validate not transferring to self
	if senderID == receiverID {
		return nil, errors.New("cannot transfer to yourself")
	}

	// Check if receiver exists
	receiver, err := s.userRepo.FindByID(receiverID)
	if err != nil {
		return nil, errors.New("receiver not found")
	}
	if receiver == nil {
		return nil, errors.New("receiver not found")
	}

	var transaction *models.Transaction

	// Use database transaction to ensure atomicity and handle race conditions
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Get sender wallet with row lock (FOR UPDATE)
		var senderWallet models.Wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", senderID).
			First(&senderWallet).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("sender wallet not found")
			}
			return err
		}

		// Check sufficient balance
		if senderWallet.Balance < amount {
			return errors.New("insufficient balance")
		}

		// Get receiver wallet with row lock (FOR UPDATE)
		var receiverWallet models.Wallet
		// Lock wallets in consistent order to prevent deadlock (lower ID first)
		if bytes.Compare(senderID[:], receiverID[:]) < 0 {
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("user_id = ?", receiverID).
				First(&receiverWallet).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New("receiver wallet not found")
				}
				return err
			}
		} else {
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("user_id = ?", receiverID).
				First(&receiverWallet).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New("receiver wallet not found")
				}
				return err
			}
		}

		// Update sender balance
		newSenderBalance := senderWallet.Balance - amount
		if err := s.walletRepo.UpdateBalanceWithLock(tx, senderWallet.ID, newSenderBalance); err != nil {
			return err
		}

		// Update receiver balance
		newReceiverBalance := receiverWallet.Balance + amount
		if err := s.walletRepo.UpdateBalanceWithLock(tx, receiverWallet.ID, newReceiverBalance); err != nil {
			return err
		}

		// Create transaction record
		transaction = &models.Transaction{
			SenderID:   &senderID,
			ReceiverID: receiverID,
			Amount:     amount,
			Type:       models.TransactionTypeTransfer,
			Status:     models.TransactionStatusSuccess,
		}

		if err := s.transactionRepo.Create(tx, transaction); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		// Create failed transaction record
		failedTransaction := &models.Transaction{
			SenderID:   &senderID,
			ReceiverID: receiverID,
			Amount:     amount,
			Type:       models.TransactionTypeTransfer,
			Status:     models.TransactionStatusFailed,
		}
		s.transactionRepo.Create(nil, failedTransaction)

		return nil, err
	}

	return transaction, nil
}

func (s *transactionService) GetHistory(userID uuid.UUID, limit int) ([]models.Transaction, error) {
	transactions, err := s.transactionRepo.FindByUserID(userID, limit)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
