package services

import (
	"errors"
	"time"

	"github.com/AllanKariuki/fiber-go-rest-api/config"
	"github.com/AllanKariuki/fiber-go-rest-api/models"
	"github.com/AllanKariuki/fiber-go-rest-api/repositories"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	Register(dto *models.RegisterDTO) (*models.User, string, error)
	Login(dto *models.LoginDTO) (*models.User, string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(dto *models.RegisterDTO) (*models.User, string, error) {
	// Check if user exists
	existingUser, _ := s.userRepo.FindByEmail(dto.Email)
	if existingUser != nil {
		return nil, "", errors.New("user already exists")
	}

	user := &models.User{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password,
		Role:     "user",
	}

	if err := user.Hashpassword(); err != nil {
		return nil, "", err
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, "", err
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *authService) Login(dto *models.LoginDTO) (*models.User, string, error) {
	user, err := s.userRepo.FindByEmail(dto.Email)
	if err != nil {
		return nil, "", errors.New("Invalid credentials")
	}

	if !user.CheckPassword(dto.Password) {
		return nil, "", errors.New("Invalid credentials")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *authService) generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims {
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // 7 days
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}

func (s *authService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(config.AppConfig.JWTSecret), nil
	})
}