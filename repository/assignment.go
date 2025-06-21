package repository

import (
	"database/sql"
	"errors"
	"go-19/model"
)

type AssignmentRepository interface {
	Create(assignment *model.Assignment) error
	FindByID(id int) (*model.Assignment, error)
	FindAll() ([]model.Assignment, error)
	Update(assignment *model.Assignment) error
	Delete(id int) error
}

type assignmentRepository struct {
	db *sql.DB
}

func NewAssignmentRepository(db *sql.DB) AssignmentRepository {
	return &assignmentRepository{db: db}
}

func (r *assignmentRepository) Create(a *model.Assignment) error {
	query := `
		INSERT INTO assignments (course_id, lecturer_id, title, description, deadline, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id
	`
	return r.db.QueryRow(query, a.CourseID, a.LecturerID, a.Title, a.Description, a.Deadline).Scan(&a.ID)
}

func (r *assignmentRepository) FindByID(id int) (*model.Assignment, error) {
	query := `
		SELECT id, course_id, lecturer_id, title, description, deadline, created_at, updated_at, deleted_at
		FROM assignments
		WHERE id = $1 AND deleted_at IS NULL
	`
	var a model.Assignment
	err := r.db.QueryRow(query, id).Scan(
		&a.ID, &a.CourseID, &a.LecturerID,
		&a.Title, &a.Description, &a.Deadline,
		&a.CreatedAt, &a.UpdatedAt, &a.DeletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &a, err
}

func (r *assignmentRepository) FindAll() ([]model.Assignment, error) {
	query := `
		SELECT id, course_id, lecturer_id, title, description, deadline, created_at, updated_at, deleted_at
		FROM assignments
		WHERE deleted_at IS NULL
		ORDER BY deadline ASC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assignments []model.Assignment
	for rows.Next() {
		var a model.Assignment
		err := rows.Scan(
			&a.ID, &a.CourseID, &a.LecturerID,
			&a.Title, &a.Description, &a.Deadline,
			&a.CreatedAt, &a.UpdatedAt, &a.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, a)
	}
	return assignments, nil
}

func (r *assignmentRepository) Update(a *model.Assignment) error {
	query := `
		UPDATE assignments
		SET course_id = $1, lecturer_id = $2, title = $3, description = $4, deadline = $5, updated_at = NOW()
		WHERE id = $6 AND deleted_at IS NULL
	`
	result, err := r.db.Exec(query, a.CourseID, a.LecturerID, a.Title, a.Description, a.Deadline, a.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("assignment not found")
	}
	return nil
}

func (r *assignmentRepository) Delete(id int) error {
	query := `
		UPDATE assignments SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL
	`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("assignment not found")
	}
	return nil
}
