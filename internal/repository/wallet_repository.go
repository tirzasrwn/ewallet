package repository

import (
	"errors"
	"ewallet/internal/models"

	"gorm.io/gorm"
)

type WalletRepository interface {
	Create(wallet *models.Wallet) error
	FindByUserID(userID uint) (*models.Wallet, error)
	UpdateBalance(walletID uint, amount float64) error
	UpdateBalanceWithLock(tx *gorm.DB, walletID uint, amount float64) error
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

func (r *walletRepository) FindByUserID(userID uint) (*models.Wallet, error) {
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

func (r *walletRepository) UpdateBalance(walletID uint, amount float64) error {
	return r.db.Model(&models.Wallet{}).Where("id = ?", walletID).Update("balance", amount).Error
}

func (r *walletRepository) UpdateBalanceWithLock(tx *gorm.DB, walletID uint, amount float64) error {
	return tx.Model(&models.Wallet{}).
		Where("id = ?", walletID).
		Update("balance", amount).Error
}
