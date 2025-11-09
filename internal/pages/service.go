package pages

import (
	"errors"
)

var (
	ErrPageNotFound = errors.New("page not found")
)

type Service interface {
	CreatePage(input PageInput, userID uint) (*Page, error)
	GetAllPagesByUser(userID uint) ([]Page, error)
	GetPageByID(id, userID uint) (*Page, error)
	UpdatePage(id uint, input PageInput, userID uint) (*Page, error)
	DeletePage(id, userID uint) error
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

func (s *service) GetAllPagesByUser(userID uint) ([]Page, error) {
	return s.repo.GetAllPagesByUser(userID)
}

func (s *service) GetPageByID(id, userID uint) (*Page, error) {
	page, err := s.repo.GetPageByID(id)
	if err != nil {
		return nil, err
	}
	if page == nil || page.UserID != userID {
		return nil, ErrPageNotFound
	}
	return page, nil
}

func (s *service) UpdatePage(id uint, input PageInput, userID uint) (*Page, error) {
	page, err := s.repo.GetPageByID(id)
	if err != nil {
		return nil, err
	}
	if page == nil || page.UserID != userID {
		return nil, ErrPageNotFound
	}

	page.Title = input.Title
	page.Content = input.Content

	if err := s.repo.UpdatePage(page); err != nil {
		return nil, err
	}
	return page, nil
}

func (s *service) DeletePage(id, userID uint) error {
	page, err := s.repo.GetPageByID(id)
	if err != nil {
		return err
	}
	if page == nil || page.UserID != userID {
		return ErrPageNotFound
	}
	return s.repo.DeletePage(id)
}
