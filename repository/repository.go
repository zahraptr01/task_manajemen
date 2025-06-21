package repository

import "database/sql"

type Repository struct {
	AssignmentRepo AssignmentRepository
	SubmissionRepo SubmissionRepo
	UserRepo       UserRepository
}

func NewRepository(db *sql.DB) Repository {
	return Repository{
		AssignmentRepo: NewAssignmentRepository(db),
		SubmissionRepo: NewSubmissionRepo(db),
		UserRepo:       NewUserRepository(db),
	}
}
