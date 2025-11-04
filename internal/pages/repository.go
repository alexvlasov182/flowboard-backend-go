package pages

import (
	"gorm.io/gorm"
)

type Repository interface {
	CreatePage(page *Page) (*Page, error)
	GetAllPages() ([]Page, error)
	GetAllPagesByUser(userID uint) ([]Page, error)
	GetPageByID(id uint) (*Page, error)
	UpdatePage(page *Page) error
	DeletePage(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreatePage(page *Page) (*Page, error) {
	if err := r.db.Create(page).Error; err != nil {
		return nil, err
	}
	return page, nil
}

func (r *repository) GetAllPages() ([]Page, error) {
	var pages []Page
	if err := r.db.Find(&pages).Error; err != nil {
		return nil, err
	}
	return pages, nil
}

func (r *repository) GetPageByID(id uint) (*Page, error) {
	var page Page
	if err := r.db.First(&page, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &page, nil
}

func (r *repository) GetAllPagesByUser(userID uint) ([]Page, error) {
	var pages []Page
	if err := r.db.Where("user_id = ?", userID).Find(&pages).Error; err != nil {
		return nil, err
	}
	return pages, nil
}

func (r *repository) UpdatePage(page *Page) error {
	return r.db.Save(page).Error
}

func (r *repository) DeletePage(id uint) error {
	return r.db.Delete(&Page{}, id).Error
}
