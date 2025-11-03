package users

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// Mocked Repository
type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) CreateUser(u *User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *MockRepo) GetUserByEmail(email string) (*User, error) {
	args := m.Called(email)
	u := args.Get(0)
	if u == nil {
		return nil, args.Error(1)
	}
	return u.(*User), args.Error(1)
}

func (m *MockRepo) GetUserByID(id uint) (*User, error) {
	args := m.Called(id)
	u := args.Get(0)
	if u == nil {
		return nil, args.Error(1)
	}
	return u.(*User), args.Error(1)
}

func TestRegister_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	s := NewService(mockRepo)

	name := "Alex"
	email := "alex@example.com"
	password := "secret123"

	// repository expectations
	mockRepo.On("GetUserByEmail", email).Return(nil, nil)
	mockRepo.On("CreateUser", mock.AnythingOfType("*users.User")).Return(nil)

	user, err := s.Register(name, email, password)
	assert.NoError(t, err)
	assert.Equal(t, name, user.Name)
	assert.Equal(t, email, user.Email)
	assert.Empty(t, user.Password)

	mockRepo.AssertExpectations(t)
}

func TestRegister_AlreadyExists(t *testing.T) {
	mockRepo := new(MockRepo)
	s := NewService(mockRepo)

	existingUser := &User{
		ID:        1,
		Name:      "Alex",
		Email:     "alex@example.com",
		Password:  "hashed",
		CreatedAt: time.Now(),
	}

	mockRepo.On("GetUserByEmail", existingUser.Email).Return(existingUser, nil)

	user, err := s.Register(existingUser.Name, existingUser.Email, "secret123")
	assert.Nil(t, user)
	assert.Equal(t, ErrUserExists, err)
}

func TestAuthenticate_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	s := NewService(mockRepo)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	u := &User{ID: 1, Email: "alex@example.com", Password: string(hashedPassword)}

	mockRepo.On("GetUserByEmail", u.Email).Return(u, nil)

	user, err := s.Authenticate("alex@example.com", "secret123")
	assert.NoError(t, err)
	assert.Equal(t, u.Email, user.Email)
	assert.Empty(t, user.Password)
}

func TestAuthenticate_InvalidPassword(t *testing.T) {
	mockRepo := new(MockRepo)
	s := NewService(mockRepo)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	u := &User{ID: 1, Email: "alex@example.com", Password: string(hashedPassword)}

	mockRepo.On("GetUserByEmail", u.Email).Return(u, nil)

	user, err := s.Authenticate(u.Email, "wrongpassword")
	assert.Nil(t, user)
	assert.Equal(t, ErrInvalidCredentials, err)
}

func TestGetByID_UserExists(t *testing.T) {
	mockRepo := new(MockRepo)
	s := NewService(mockRepo)

	u := &User{ID: 1, Name: "Alex", Email: "alex@example.com"}
	mockRepo.On("GetUserByID", uint(1)).Return(u, nil)

	user, err := s.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, u.Email, user.Email)
}

func TestGetByID_UserNotFound(t *testing.T) {
	mockRepo := new(MockRepo)
	s := NewService(mockRepo)

	mockRepo.On("GetUserByID", uint(2)).Return(nil, nil)

	user, err := s.GetByID(2)
	assert.NoError(t, err)
	assert.Nil(t, user)
}

func TestGetByID_RepoError(t *testing.T) {
	mockRepo := new(MockRepo)
	s := NewService(mockRepo)

	mockRepo.On("GetUserByID", uint(3)).Return(nil, errors.New("DB error"))

	user, err := s.GetByID(3)
	assert.Nil(t, user)
	assert.EqualError(t, err, "DB error")
}
