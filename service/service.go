package service

import "go-19/repository"

type Service struct {
	AssignmentService AssignmentService
	SubmissionService SubmissionService
	UserService       UserService
}

func NewService(repo repository.Repository) Service {
	return Service{
		AssignmentService: NewAssignmentService(repo),
		SubmissionService: NewSubmissionService(repo),
		UserService:       NewUserService(repo),
	}
}
