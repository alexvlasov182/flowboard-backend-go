package users

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists         = errors.New("user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Service interface {
	Register(name, email, password string) (*User, error)
	Authenticate(email, password string) (*User, error)
	GetByID(id uint) (*User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// Register implements Service.
func (s *service) Register(name, email, password string) (*User, error) {
	// Check if user already exists
	ex, err := s.repo.GetUserByEmail(email)

	if err != nil {
		return nil, err
	}
	if ex != nil {
		return nil, ErrUserExists
	}

	// hash
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &User{
		Name:      name,
		Email:     email,
		Password:  string(hashed),
		CreatedAt: time.Now(),
	}

	err = s.repo.CreateUser(u)
	if err != nil {
		return nil, err
	}

	u.Password = ""

	return u, nil
}

// Authenticate implements Service.
func (s *service) Authenticate(email string, password string) (*User, error) {
	u, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	u.Password = ""

	return u, nil
}

// GetByID implements Service.
func (s *service) GetByID(id uint) (*User, error) {
	u, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, nil
	}
	u.Password = ""
	return u, nil
}
