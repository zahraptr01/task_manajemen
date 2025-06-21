package handler

import (
	"go-19/model"
	"go-19/service"
	"html/template"
	"net/http"
	"strconv"
)

type SubmissionHandler struct {
	SubmissionService service.SubmissionService
	UserService       service.UserService
	AssignmentService service.AssignmentService
	Template          *template.Template
}

func NewSubmissionHandler(submissionService service.SubmissionService, userService service.UserService, assignmentService service.AssignmentService, tmpl *template.Template) *SubmissionHandler {
	return &SubmissionHandler{
		SubmissionService: submissionService,
		UserService:       userService,
		AssignmentService: assignmentService,
		Template:          tmpl,
	}
}

func (h *SubmissionHandler) Home(w http.ResponseWriter, r *http.Request) {
	submissions, err := h.SubmissionService.GetAllSubmissions()
	if err != nil {
		http.Error(w, "Gagal mengambil data submission", http.StatusInternalServerError)
		return
	}

	status := r.URL.Query().Get("status")

	data := struct {
		Submissions []model.Submission
		Status      string
	}{
		Submissions: submissions,
		Status:      status,
	}

	h.Template.ExecuteTemplate(w, "lecturer_home", data)
}

func (h *SubmissionHandler) ShowGradeForm(w http.ResponseWriter, r *http.Request) {
	studentIDStr := r.URL.Query().Get("student_id")
	assignmentIDStr := r.URL.Query().Get("assignment_id")

	studentID, err := strconv.Atoi(studentIDStr)
	if err != nil {
		http.Error(w, "Invalid student_id", http.StatusBadRequest)
		return
	}

	assignmentID, err := strconv.Atoi(assignmentIDStr)
	if err != nil {
		http.Error(w, "Invalid assignment_id", http.StatusBadRequest)
		return
	}

	student, err := h.UserService.GetUserByID(studentID)
	if err != nil {
		http.Error(w, "Student not found", http.StatusInternalServerError)
		return
	}

	assignment, err := h.AssignmentService.GetAssignmentByID(assignmentID)
	if err != nil {
		http.Error(w, "Assignment not found", http.StatusInternalServerError)
		return
	}

	data := struct {
		StudentID       int
		AssignmentID    int
		StudentName     string
		AssignmentTitle string
	}{
		StudentID:       student.ID,
		AssignmentID:    assignment.ID,
		StudentName:     student.Name,
		AssignmentTitle: assignment.Title,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.Template.ExecuteTemplate(w, "grade_form", data)
}

func (h *SubmissionHandler) GradeSubmission(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Gagal parsing form", http.StatusBadRequest)
		return
	}

	studentID, err := strconv.Atoi(r.FormValue("student_id"))
	if err != nil {
		http.Error(w, "Invalid student_id", http.StatusBadRequest)
		return
	}

	assignmentID, err := strconv.Atoi(r.FormValue("assignment_id"))
	if err != nil {
		http.Error(w, "Invalid assignment_id", http.StatusBadRequest)
		return
	}

	gradeStr := r.FormValue("grade")
	grade, err := strconv.ParseFloat(gradeStr, 64)
	if err != nil {
		http.Error(w, "Invalid grade", http.StatusBadRequest)
		return
	}

	err = h.SubmissionService.GradeSubmission(studentID, assignmentID, grade)
	if err != nil {
		http.Error(w, "Gagal memberi nilai: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/lecturer/home", http.StatusSeeOther)
}

func (h *SubmissionHandler) ResetSubmission(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Gagal membaca form", http.StatusBadRequest)
		return
	}

	studentID, err := strconv.Atoi(r.FormValue("student_id"))
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	assignmentID, err := strconv.Atoi(r.FormValue("assignment_id"))
	if err != nil {
		http.Error(w, "Invalid assignment ID", http.StatusBadRequest)
		return
	}

	err = h.SubmissionService.DeleteSubmission(studentID, assignmentID)
	if err != nil {
		http.Error(w, "Gagal menghapus submission: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/lecturer/home?status=reset", http.StatusSeeOther)
}
