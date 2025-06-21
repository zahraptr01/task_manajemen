package service

import (
	"go-19/model"
	"go-19/repository"
)

type UserService interface {
	GetUserByID(id int) (model.User, error)
	CreateUser(user *model.User) error
	FindByEmail(email string) (*model.User, error)
}

type userService struct {
	Repo repository.Repository
}

func NewUserService(repo repository.Repository) UserService {
	return &userService{Repo: repo}
}

func (s *userService) GetUserByID(id int) (model.User, error) {
	return s.Repo.UserRepo.GetUserByID(id)
}

func (s *userService) CreateUser(user *model.User) error {
	return s.Repo.UserRepo.Create(user)
}

func (s *userService) FindByEmail(email string) (*model.User, error) {
	return s.Repo.UserRepo.FindByEmail(email)
}
