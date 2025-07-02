package service

import (
	"errors"

	"github.com/1sh-repalto/url-monitoring-api/internal/model"
	"github.com/1sh-repalto/url-monitoring-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{userRepo}
}

func (s *AuthService) Login(email, password string) (*model.User, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *AuthService) Signup(email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := &model.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	return s.userRepo.CreateUser(user)
}
