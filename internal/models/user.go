package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Email     string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	Wallet    Wallet         `gorm:"foreignKey:UserID" json:"wallet,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate hook to generate UUID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// HashPassword hashes the user password
func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword checks if the provided password matches the user's password
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// UserResponse represents the user data returned in API responses
type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

// ToResponse converts User model to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}
