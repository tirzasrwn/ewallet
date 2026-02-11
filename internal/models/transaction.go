package models

import (
	"time"

	"gorm.io/gorm"
)

type TransactionType string
type TransactionStatus string

const (
	TransactionTypeTopUp    TransactionType = "topup"
	TransactionTypeTransfer TransactionType = "transfer"

	TransactionStatusPending TransactionStatus = "pending"
	TransactionStatusSuccess TransactionStatus = "success"
	TransactionStatusFailed  TransactionStatus = "failed"
)

type Transaction struct {
	ID         uint              `gorm:"primaryKey" json:"id"`
	SenderID   *uint             `gorm:"index" json:"sender_id,omitempty"`
	ReceiverID uint              `gorm:"index;not null" json:"receiver_id"`
	Amount     float64           `gorm:"type:decimal(15,2);not null" json:"amount"`
	Type       TransactionType   `gorm:"type:varchar(20);not null" json:"type"`
	Status     TransactionStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	Sender     *User             `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
	Receiver   User              `gorm:"foreignKey:ReceiverID" json:"receiver,omitempty"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
	DeletedAt  gorm.DeletedAt    `gorm:"index" json:"-"`
}

// TransactionResponse represents the transaction data returned in API responses
type TransactionResponse struct {
	ID         uint              `json:"id"`
	SenderID   *uint             `json:"sender_id,omitempty"`
	ReceiverID uint              `json:"receiver_id"`
	Amount     float64           `json:"amount"`
	Type       TransactionType   `json:"type"`
	Status     TransactionStatus `json:"status"`
	CreatedAt  time.Time         `json:"created_at"`
}

// ToResponse converts Transaction model to TransactionResponse
func (t *Transaction) ToResponse() TransactionResponse {
	return TransactionResponse{
		ID:         t.ID,
		SenderID:   t.SenderID,
		ReceiverID: t.ReceiverID,
		Amount:     t.Amount,
		Type:       t.Type,
		Status:     t.Status,
		CreatedAt:  t.CreatedAt,
	}
}
