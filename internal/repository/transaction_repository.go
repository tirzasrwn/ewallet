package repository

import (
	"ewallet/internal/models"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(tx *gorm.DB, transaction *models.Transaction) error
	FindByUserID(userID uint, limit int) ([]models.Transaction, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(tx *gorm.DB, transaction *models.Transaction) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Create(transaction).Error
}

func (r *transactionRepository) FindByUserID(userID uint, limit int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	query := r.db.Where("sender_id = ? OR receiver_id = ?", userID, userID).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&transactions).Error
	return transactions, err
}
