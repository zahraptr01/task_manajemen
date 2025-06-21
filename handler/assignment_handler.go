package handler

import (
	"go-19/model"
	"go-19/service"
	"html/template"
	"net/http"
	"strconv"
)

// AssignmentHandler is responsible for assignment logic (list, submit, display form).
type AssignmentHandler struct {
	Service  service.Service
	Template *template.Template
}

// Constructor for AssignmentHandler
func NewAssignmentHandler(server service.Service, template *template.Template) AssignmentHandler {
	return AssignmentHandler{
		Service:  server,
		Template: template,
	}
}

// Display assignment list for students
func (assignmentHandler *AssignmentHandler) ListAssignments(w http.ResponseWriter, r *http.Request) {
	// Retrieve user_id from login cookie
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Convert user_id cookie to int
	userID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Retrieve all assignments from service
	assignments, err := assignmentHandler.Service.AssignmentService.GetAllAssignments()
	if err != nil {
		http.Error(w, "Failed to fetch assignments: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Retrieve student data to display name on page
	student, err := assignmentHandler.Service.UserService.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Failed to fetch student data", http.StatusInternalServerError)
		return
	}

	status := r.URL.Query().Get("status")

	data := struct {
		StudentName string
		Assignments []model.Assignment
		Status      string
	}{
		StudentName: student.Name,
		Assignments: assignments,
		Status:      status,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := assignmentHandler.Template.ExecuteTemplate(w, "assignment_list", data); err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
	}
}

// Handle assignment submit process by students
func (assignmentHandler *AssignmentHandler) SubmitAssignment(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse form for file upload (max 10MB)
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Gagal membaca form: "+err.Error(), http.StatusBadRequest)
			return
		}

		assignmentID, err := strconv.Atoi(r.FormValue("assignment_id"))
		if err != nil {
			http.Error(w, "Invalid assignment ID", http.StatusBadRequest)
			return
		}

		studentID, err := strconv.Atoi(r.FormValue("student_id"))
		if err != nil {
			http.Error(w, "Invalid student ID", http.StatusBadRequest)
			return
		}

		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "File tidak valid: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		_, err = assignmentHandler.Service.AssignmentService.SubmitAssignment(studentID, assignmentID, file, fileHeader)
		if err != nil {
			http.Error(w, "Gagal submit: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect and send status=submitted parameter
		http.Redirect(w, r, "/student/home?status=submitted", http.StatusSeeOther)
	}
}

// Display assignment submit form based on assignment IDs
func (h *AssignmentHandler) ShowSubmitForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/student/home", http.StatusSeeOther)
		return
	}

	assignmentIDStr := r.URL.Query().Get("assignment_id")
	assignmentID, err := strconv.Atoi(assignmentIDStr)
	if err != nil {
		http.Error(w, "Invalid assignment ID", http.StatusBadRequest)
		return
	}

	assignment, err := h.Service.AssignmentService.GetAssignmentByID(assignmentID)
	if err != nil {
		http.Error(w, "Assignment not found", http.StatusNotFound)
		return
	}

	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	userID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.Service.UserService.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	data := struct {
		Assignment  model.Assignment
		StudentID   int
		StudentName string
	}{
		Assignment:  *assignment,
		StudentID:   user.ID,
		StudentName: user.Name,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.Template.ExecuteTemplate(w, "submit_form", data)
}
