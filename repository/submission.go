package repository

import (
	"database/sql"
	"fmt"
	"go-19/model"
)

type SubmissionRepo interface {
	CountByStudentAndAssignment(studentID, assignmentID int) (int64, error)
	Create(submission *model.Submission) error
	GetAllWithStudentAndAssignment() ([]model.Submission, error)
	FindByStudentAndAssignment(studentID, assignmentID int) (*model.Submission, error)
	UpdateGrade(sub *model.Submission) error
	DeleteByStudentAndAssignment(studentID, assignmentID int) error
}

type submissionRepo struct {
	db *sql.DB
}

func NewSubmissionRepo(db *sql.DB) SubmissionRepo {
	return &submissionRepo{db}
}

func (r *submissionRepo) CountByStudentAndAssignment(studentID, assignmentID int) (int64, error) {
	var count int64
	err := r.db.QueryRow("SELECT COUNT(*) FROM submissions WHERE student_id=$1 AND assignment_id=$2", studentID, assignmentID).Scan(&count)
	return count, err
}

func (r *submissionRepo) Create(sub *model.Submission) error {
	_, err := r.db.Exec("INSERT INTO submissions (assignment_id, student_id, submitted_at, file_url, status) VALUES ($1, $2, $3, $4, $5)",
		sub.AssignmentID, sub.StudentID, sub.SubmittedAt, sub.FileURL, sub.Status)
	return err
}

func (r *submissionRepo) GetAllWithStudentAndAssignment() ([]model.Submission, error) {
	query := `
		SELECT s.id, s.assignment_id, s.student_id, u.name as student_name,
		       a.title as assignment_title, s.file_url, s.status, s.grade
		FROM submissions s
		JOIN users u ON s.student_id = u.id
		JOIN assignments a ON s.assignment_id = a.id
		ORDER BY s.submitted_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submissions []model.Submission
	for rows.Next() {
		var s model.Submission
		err := rows.Scan(&s.ID, &s.AssignmentID, &s.StudentID, &s.StudentName, &s.AssignmentTitle, &s.FileURL, &s.Status, &s.Grade)
		if err != nil {
			return nil, err
		}
		submissions = append(submissions, s)
	}
	fmt.Printf("data %+v", submissions)
	return submissions, nil
}

func (r *submissionRepo) FindByStudentAndAssignment(studentID, assignmentID int) (*model.Submission, error) {
	query := `SELECT id, assignment_id, student_id, submitted_at, file_url, status, grade 
			  FROM submissions 
			  WHERE student_id = $1 AND assignment_id = $2 LIMIT 1`

	row := r.db.QueryRow(query, studentID, assignmentID)

	var sub model.Submission
	err := row.Scan(
		&sub.ID,
		&sub.AssignmentID,
		&sub.StudentID,
		&sub.SubmittedAt,
		&sub.FileURL,
		&sub.Status,
		&sub.Grade,
	)

	if err != nil {
		return nil, err
	}

	return &sub, nil
}

func (r *submissionRepo) UpdateGrade(sub *model.Submission) error {
	query := `UPDATE submissions SET grade = $1 WHERE student_id = $2 AND assignment_id = $3`
	_, err := r.db.Exec(query, sub.Grade, sub.StudentID, sub.AssignmentID)
	return err
}

func (r *submissionRepo) DeleteByStudentAndAssignment(studentID, assignmentID int) error {
	query := `DELETE FROM submissions WHERE student_id = $1 AND assignment_id = $2`
	_, err := r.db.Exec(query, studentID, assignmentID)
	return err
}
