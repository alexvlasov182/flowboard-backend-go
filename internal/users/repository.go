package users

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	// Define repository methods here, e.g.:
	CreateUser(u *User) error
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id uint) (*User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// CreateUser implements Repository.
func (r *repository) CreateUser(u *User) error {
	return r.db.Create(u).Error
}

// GetUserByEmail implements Repository.
func (r *repository) GetUserByEmail(email string) (*User, error) {
	var u User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *repository) GetUserByID(id uint) (*User, error) {
	var u User
	if err := r.db.First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}
