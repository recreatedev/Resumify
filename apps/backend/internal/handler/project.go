package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/recreatedev/Resumify/internal/middleware"
	"github.com/recreatedev/Resumify/internal/model/project"
	"github.com/recreatedev/Resumify/internal/server"
	"github.com/recreatedev/Resumify/internal/service"
)

type ProjectHandler struct {
	Handler
	projectService *service.ProjectService
}

func NewProjectHandler(s *server.Server, projectService *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		Handler:        NewHandler(s),
		projectService: projectService,
	}
}

func (h *ProjectHandler) CreateProject(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *project.CreateProjectRequest) (*project.ProjectResponse, error) {
			userID := middleware.GetUserID(c)
			return h.projectService.CreateProject(c.Request().Context(), userID, payload)
		},
		http.StatusCreated,
		&project.CreateProjectRequest{},
	)(c)
}

func (h *ProjectHandler) GetProjectByID(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, req *GetProjectByIDRequest) (*project.ProjectResponse, error) {
			userID := middleware.GetUserID(c)
			projectID, err := req.ParseID()
			if err != nil {
				return nil, err
			}
			return h.projectService.GetProjectByID(c.Request().Context(), userID, projectID)
		},
		http.StatusOK,
		&GetProjectByIDRequest{},
	)(c)
}

func (h *ProjectHandler) GetProjectsByResumeID(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, req *GetProjectsByResumeIDRequest) ([]project.ProjectResponse, error) {
			userID := middleware.GetUserID(c)
			resumeID, err := req.ParseResumeID()
			if err != nil {
				return nil, err
			}
			return h.projectService.GetProjectsByResumeID(c.Request().Context(), userID, resumeID)
		},
		http.StatusOK,
		&GetProjectsByResumeIDRequest{},
	)(c)
}

func (h *ProjectHandler) UpdateProject(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, req *UpdateProjectRequest) (*project.ProjectResponse, error) {
			userID := middleware.GetUserID(c)
			projectID, err := req.ParseID()
			if err != nil {
				return nil, err
			}
			return h.projectService.UpdateProject(c.Request().Context(), userID, projectID, req.UpdateProjectRequest)
		},
		http.StatusOK,
		&UpdateProjectRequest{},
	)(c)
}

func (h *ProjectHandler) BulkUpdateProjectOrder(c echo.Context) error {
	return HandleNoContent(
		h.Handler,
		func(c echo.Context, payload *project.BulkUpdateProjectsRequest) error {
			userID := middleware.GetUserID(c)
			return h.projectService.BulkUpdateProjectOrder(c.Request().Context(), userID, payload)
		},
		http.StatusNoContent,
		&project.BulkUpdateProjectsRequest{},
	)(c)
}

func (h *ProjectHandler) DeleteProject(c echo.Context) error {
	return HandleNoContent(
		h.Handler,
		func(c echo.Context, req *DeleteProjectRequest) error {
			userID := middleware.GetUserID(c)
			projectID, err := req.ParseID()
			if err != nil {
				return err
			}
			return h.projectService.DeleteProject(c.Request().Context(), userID, projectID)
		},
		http.StatusNoContent,
		&DeleteProjectRequest{},
	)(c)
}

// Request DTOs

type GetProjectByIDRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}

func (r *GetProjectByIDRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *GetProjectByIDRequest) ParseID() (uuid.UUID, error) {
	return uuid.Parse(r.ID)
}

type GetProjectsByResumeIDRequest struct {
	ResumeID string `param:"resumeId" validate:"required,uuid"`
}

func (r *GetProjectsByResumeIDRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *GetProjectsByResumeIDRequest) ParseResumeID() (uuid.UUID, error) {
	return uuid.Parse(r.ResumeID)
}

type UpdateProjectRequest struct {
	ID string `param:"id" validate:"required,uuid"`
	*project.UpdateProjectRequest
}

func (r *UpdateProjectRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err
	}
	return r.UpdateProjectRequest.Validate()
}

func (r *UpdateProjectRequest) ParseID() (uuid.UUID, error) {
	return uuid.Parse(r.ID)
}

type DeleteProjectRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}

func (r *DeleteProjectRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *DeleteProjectRequest) ParseID() (uuid.UUID, error) {
	return uuid.Parse(r.ID)
}
