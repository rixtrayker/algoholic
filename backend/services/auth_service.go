package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/yourusername/algoholic/config"
	"github.com/yourusername/algoholic/models"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidToken       = errors.New("invalid token")
)

// Claims represents JWT claims
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// AuthService handles authentication operations
type AuthService struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewAuthService creates a new auth service
func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{
		db:  db,
		cfg: cfg,
	}
}

// Register creates a new user account
func (s *AuthService) Register(username, email, password string) (*models.User, error) {
	// Check if user already exists
	var existingUser models.User
	if err := s.db.Where("username = ? OR email = ?", username, email).First(&existingUser).Error; err == nil {
		return nil, ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), s.cfg.Auth.BCryptCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(username, password string) (string, *models.User, error) {
	var user models.User
	if err := s.db.Where("username = ? OR email = ?", username, username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, ErrInvalidCredentials
		}
		return "", nil, err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := s.GenerateToken(&user)
	if err != nil {
		return "", nil, err
	}

	// Update last active
	s.db.Model(&user).Update("last_active_at", time.Now())

	return token, &user, nil
}

// GenerateToken creates a JWT token for a user
func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	expirationTime := time.Now().Add(time.Duration(s.cfg.Auth.JWTExpiry) * time.Hour)

	claims := &Claims{
		UserID:   user.UserID,
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    s.cfg.App.Name,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.Auth.JWTSecret))
}

// ValidateToken validates a JWT token and returns the claims
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.Auth.JWTSecret), nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(userID int) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// UpdatePassword changes a user's password
func (s *AuthService) UpdatePassword(userID int, oldPassword, newPassword string) error {
	user, err := s.GetUserByID(userID)
	if err != nil {
		return err
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return ErrInvalidCredentials
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), s.cfg.Auth.BCryptCost)
	if err != nil {
		return err
	}

	// Update password
	return s.db.Model(&user).Update("password_hash", string(hashedPassword)).Error
}
