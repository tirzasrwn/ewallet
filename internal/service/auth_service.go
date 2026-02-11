package service

import (
	"errors"
	"ewallet/internal/models"
	"ewallet/internal/repository"
	"ewallet/pkg/utils"

	"gorm.io/gorm"
)

type AuthService interface {
	Register(name, email, password string) (*models.User, error)
	Login(email, password string) (string, *models.User, error)
}

type authService struct {
	userRepo   repository.UserRepository
	walletRepo repository.WalletRepository
	jwtUtil    *utils.JWTUtil
	db         *gorm.DB
}

func NewAuthService(
	userRepo repository.UserRepository,
	walletRepo repository.WalletRepository,
	jwtUtil *utils.JWTUtil,
	db *gorm.DB,
) AuthService {
	return &authService{
		userRepo:   userRepo,
		walletRepo: walletRepo,
		jwtUtil:    jwtUtil,
		db:         db,
	}
}

func (s *authService) Register(name, email, password string) (*models.User, error) {
	// Validate input
	if name == "" || email == "" || password == "" {
		return nil, errors.New("all fields are required")
	}

	if len(password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}

	// Check if user already exists
	existingUser, _ := s.userRepo.FindByEmail(email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Create user with transaction
	var user models.User
	err := s.db.Transaction(func(tx *gorm.DB) error {
		user = models.User{
			Name:  name,
			Email: email,
		}

		if err := user.HashPassword(password); err != nil {
			return err
		}

		if err := s.userRepo.Create(&user); err != nil {
			return err
		}

		// Create wallet for user
		wallet := models.Wallet{
			UserID:  user.ID,
			Balance: 0,
		}

		if err := s.walletRepo.Create(&wallet); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *authService) Login(email, password string) (string, *models.User, error) {
	// Validate input
	if email == "" || password == "" {
		return "", nil, errors.New("email and password are required")
	}

	// Find user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	// Check password
	if err := user.CheckPassword(password); err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := s.jwtUtil.GenerateToken(user.ID, user.Email)
	if err != nil {
		return "", nil, errors.New("failed to generate token")
	}

	return token, user, nil
}
