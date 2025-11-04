package pages

import (
	"errors"
	"time"
)

var ErrPageNotFound = errors.New("page not found")

type Service interface {
	CreatePage(input PageInput, userID uint) (*Page, error)
	GetAllPages(userID uint) ([]Page, error)

	GetPageByID(id uint, userID uint) (*Page, error)
	UpdatePage(id uint, input PageInput, userID uint) (*Page, error)
	DeletePage(id uint, userID uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreatePage(input PageInput, userID uint) (*Page, error) {
	page := &Page{
		Title:   input.Title,
		Content: input.Content,
		UserID:  userID,
	}
	return s.repo.CreatePage(page)
}

func (s *service) GetAllPages(userID uint) ([]Page, error) {
	return s.repo.GetAllPagesByUser(userID)
}

func (s *service) GetPageByID(id uint, userID uint) (*Page, error) {
	p, err := s.repo.GetPageByID(id)
	if err != nil {
		return nil, err
	}
	if p == nil || p.UserID != userID {
		return nil, ErrPageNotFound
	}
	return p, nil
}

func (s *service) UpdatePage(id uint, input PageInput, userID uint) (*Page, error) {
	p, err := s.repo.GetPageByID(id)
	if err != nil {
		return nil, err
	}
	if p == nil || p.UserID != userID {
		return nil, ErrPageNotFound
	}

	p.Title = input.Title
	p.Content = input.Content
	p.UpdatedAt = time.Now()

	if err := s.repo.UpdatePage(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *service) DeletePage(id uint, userID uint) error {
	p, err := s.repo.GetPageByID(id)
	if err != nil {
		return err
	}
	if p == nil || p.UserID != userID {
		return ErrPageNotFound
	}
	return s.repo.DeletePage(id)
}
