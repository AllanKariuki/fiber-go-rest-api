package services

import (
	"errors"
	"github.com/AllanKariuki/fiber-go-rest-api/models"
	"github.com/AllanKariuki/fiber-go-rest-api/repositories"
)

type UserService interface {
	GetAllUsers() ([]models.User, error)
	GetUserByID(id uint) (*models.User, error)
	CreateUser(dto *models.RegisterDTO) (*models.User, error)
	UpdateUser(id uint, dto *models.UpdateUserDTO) (*models.User, error)
	DeleteUser(id uint) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.repo.FindAll()
}

func (s *userService) GetUserByID(id uint) (*models.User, error) {
	return s.repo.FindByID(id)
}

func (s *userService) CreateUser(dto *models.RegisterDTO) (*models.User, error) {
	// Check if user already exists
	existingUser, _ := s.repo.FindByEmail(dto.Email)
	if existingUser != nil {
		return nil, errors.New("User with this email already exists")
	}

	user := &models.User{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password,
		Role:    "user",
	}

	// Hash password
	if err := user.Hashpassword(); err != nil {
		return nil, err
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateUser(id uint, dto *models.UpdateUserDTO) (*models.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("User not found")
	}

	if dto.Name != "" {
		user.Name = dto.Name
	}

	if dto.Email != "" {
		user.Email = dto.Email
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) DeleteUser(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("User not found")
	}

	return s.repo.Delete(id)
}

