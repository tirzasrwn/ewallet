package repository

import (
	"errors"
	"ewallet/internal/models"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type WalletRepository interface {
	Create(wallet *models.Wallet) error
	FindByUserID(userID uuid.UUID) (*models.Wallet, error)
	UpdateBalance(walletID uuid.UUID, amount float64) error
	UpdateBalanceWithLock(tx *gorm.DB, walletID uuid.UUID, amount float64) error
}

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) Create(wallet *models.Wallet) error {
	return r.db.Create(wallet).Error
}

func (r *walletRepository) FindByUserID(userID uuid.UUID) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("wallet not found")
		}
		return nil, err
	}
	return &wallet, nil
}

func (r *walletRepository) UpdateBalance(walletID uuid.UUID, amount float64) error {
	return r.db.Model(&models.Wallet{}).Where("id = ?", walletID).Update("balance", amount).Error
}

func (r *walletRepository) UpdateBalanceWithLock(tx *gorm.DB, walletID uuid.UUID, amount float64) error {
	return tx.Model(&models.Wallet{}).
		Where("id = ?", walletID).
		Update("balance", amount).Error
}
