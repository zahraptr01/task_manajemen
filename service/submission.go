// service/lecturer_service.go
package service

import (
	"fmt"
	"go-19/model"
	"go-19/repository"
)

type SubmissionService interface {
	GetAllSubmissions() ([]model.Submission, error)
	GradeSubmission(studentID, assignmentID int, grade float64) error
	DeleteSubmission(studentID, assignmentID int) error
}

type submissionService struct {
	Repo repository.Repository
}

func NewSubmissionService(subRepo repository.Repository) SubmissionService {
	return &submissionService{
		Repo: subRepo,
	}
}

func (submissionService *submissionService) GetAllSubmissions() ([]model.Submission, error) {
	return submissionService.Repo.SubmissionRepo.GetAllWithStudentAndAssignment()
}

func (s *submissionService) GradeSubmission(studentID, assignmentID int, grade float64) error {
	// Cek apakah submission-nya ada
	sub, err := s.Repo.SubmissionRepo.FindByStudentAndAssignment(studentID, assignmentID)
	if err != nil {
		return fmt.Errorf("submission not found: %w", err)
	}

	// Update nilai
	sub.Grade = &grade
	return s.Repo.SubmissionRepo.UpdateGrade(sub)
}

func (s *submissionService) DeleteSubmission(studentID, assignmentID int) error {
	return s.Repo.SubmissionRepo.DeleteByStudentAndAssignment(studentID, assignmentID)
}
