package main

import (
	"fmt"
	"go-19/database"
	"go-19/handler"
	"go-19/middleware"
	"go-19/repository"
	"go-19/service"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	middlechi "github.com/go-chi/chi/v5/middleware"
)

func main() {
	if err := database.InitDB(); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	repo := repository.NewRepository(database.DB)

	userService := service.NewUserService(repo)
	assignmentService := service.NewAssignmentService(repo)
	submissionService := service.NewSubmissionService(repo)

	svc := service.Service{
		UserService:       userService,
		AssignmentService: assignmentService,
		SubmissionService: submissionService,
	}

	tmpl := template.Must(template.ParseGlob("view/*.html"))

	authHandler := handler.NewAuthHandler(tmpl, userService)
	submissionHandler := handler.NewSubmissionHandler(submissionService, userService, assignmentService, tmpl)
	assignmentHandler := handler.NewAssignmentHandler(svc, tmpl)

	r := chi.NewRouter()
	r.Use(middlechi.Logger)

	// ✅ Static file handling
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// ✅ Public routes
	r.Get("/", authHandler.Login)
	r.Post("/login", authHandler.DoLogin)
	r.Get("/register", authHandler.Register)
	r.Post("/register", authHandler.DoRegister)

	// ✅ Student routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)
		r.Get("/student/home", assignmentHandler.ListAssignments)
		r.Get("/student/submit", assignmentHandler.ShowSubmitForm)
		r.Post("/student/submit", assignmentHandler.SubmitAssignment)
	})

	// ✅ Lecturer routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)
		r.Get("/lecturer/home", submissionHandler.Home)
		r.Get("/lecturer/grade-form", submissionHandler.ShowGradeForm)
		r.Post("/lecturer/grade", submissionHandler.GradeSubmission)
		r.Post("/lecturer/reset", submissionHandler.ResetSubmission)

	})

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
