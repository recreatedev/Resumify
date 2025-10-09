package handler

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/recreatedev/Resumify/internal/middleware"
	"github.com/recreatedev/Resumify/internal/model/resume"
	"github.com/recreatedev/Resumify/internal/server"
	"github.com/recreatedev/Resumify/internal/service"
)

type ResumeHandler struct {
	Handler
	service *service.ResumeService
}

func NewResumeHandler(s *server.Server, services *service.Services) *ResumeHandler {
	return &ResumeHandler{
		Handler: NewHandler(s),
		service: services.Resume,
	}
}

// CreateResume creates a new resume
func (h *ResumeHandler) CreateResume(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, req *resume.CreateResumeRequest) (*resume.ResumeResponse, error) {
			userID := middleware.GetUserID(c)
			return h.service.CreateResume(c.Request().Context(), userID, req)
		},
		http.StatusCreated,
		&resume.CreateResumeRequest{},
	)(c)
}

// GetResumeByID retrieves a resume by ID
func (h *ResumeHandler) GetResumeByID(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, req *GetResumeByIDRequest) (*resume.ResumeResponse, error) {
			userID := middleware.GetUserID(c)
			resumeID, err := req.ParseID()
			if err != nil {
				return nil, err
			}
			return h.service.GetResumeByID(c.Request().Context(), userID, resumeID)
		},
		http.StatusOK,
		&GetResumeByIDRequest{},
	)(c)
}

// GetResumes retrieves paginated list of resumes
func (h *ResumeHandler) GetResumes(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, req *GetResumesRequest) (*PaginatedResumesResponse, error) {
			userID := middleware.GetUserID(c)

			// Parse pagination parameters
			page := 1
			limit := 20

			if req.Page != "" {
				if p, err := strconv.Atoi(req.Page); err == nil && p > 0 {
					page = p
				}
			}

			if req.Limit != "" {
				if l, err := strconv.Atoi(req.Limit); err == nil && l > 0 && l <= 100 {
					limit = l
				}
			}

			result, err := h.service.GetResumes(c.Request().Context(), userID, page, limit)
			if err != nil {
				return nil, err
			}

			return &PaginatedResumesResponse{
				Data:       result.Data,
				Page:       result.Page,
				Limit:      result.Limit,
				Total:      result.Total,
				TotalPages: result.TotalPages,
			}, nil
		},
		http.StatusOK,
		&GetResumesRequest{},
	)(c)
}

// UpdateResume updates a resume
func (h *ResumeHandler) UpdateResume(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, req *UpdateResumeRequest) (*resume.ResumeResponse, error) {
			userID := middleware.GetUserID(c)
			resumeID, err := req.ParseID()
			if err != nil {
				return nil, err
			}
			return h.service.UpdateResume(c.Request().Context(), userID, resumeID, req.UpdateResumeRequest)
		},
		http.StatusOK,
		&UpdateResumeRequest{},
	)(c)
}

// DeleteResume deletes a resume
func (h *ResumeHandler) DeleteResume(c echo.Context) error {
	return HandleNoContent(
		h.Handler,
		func(c echo.Context, req *DeleteResumeRequest) error {
			userID := middleware.GetUserID(c)
			resumeID, err := req.ParseID()
			if err != nil {
				return err
			}
			return h.service.DeleteResume(c.Request().Context(), userID, resumeID)
		},
		http.StatusNoContent,
		&DeleteResumeRequest{},
	)(c)
}

// DuplicateResume creates a copy of an existing resume
func (h *ResumeHandler) DuplicateResume(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, req *DuplicateResumeRequest) (*resume.ResumeResponse, error) {
			userID := middleware.GetUserID(c)
			resumeID, err := req.ParseID()
			if err != nil {
				return nil, err
			}
			return h.service.DuplicateResume(c.Request().Context(), userID, resumeID)
		},
		http.StatusCreated,
		&DuplicateResumeRequest{},
	)(c)
}

// GetResumeWithSections retrieves a resume with all its sections
func (h *ResumeHandler) GetResumeWithSections(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, req *GetResumeByIDRequest) (*service.ResumeWithSectionsResponse, error) {
			userID := middleware.GetUserID(c)
			resumeID, err := req.ParseID()
			if err != nil {
				return nil, err
			}
			return h.service.GetResumeWithSections(c.Request().Context(), userID, resumeID)
		},
		http.StatusOK,
		&GetResumeByIDRequest{},
	)(c)
}

// Request DTOs

type GetResumeByIDRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}

func (r *GetResumeByIDRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *GetResumeByIDRequest) ParseID() (uuid.UUID, error) {
	return uuid.Parse(r.ID)
}

type GetResumesRequest struct {
	Page  string `query:"page" validate:"omitempty,numeric"`
	Limit string `query:"limit" validate:"omitempty,numeric"`
}

func (r *GetResumesRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

type UpdateResumeRequest struct {
	ID string `param:"id" validate:"required,uuid"`
	*resume.UpdateResumeRequest
}

func (r *UpdateResumeRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err
	}
	return r.UpdateResumeRequest.Validate()
}

func (r *UpdateResumeRequest) ParseID() (uuid.UUID, error) {
	return uuid.Parse(r.ID)
}

type DeleteResumeRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}

func (r *DeleteResumeRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *DeleteResumeRequest) ParseID() (uuid.UUID, error) {
	return uuid.Parse(r.ID)
}

type DuplicateResumeRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}

func (r *DuplicateResumeRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *DuplicateResumeRequest) ParseID() (uuid.UUID, error) {
	return uuid.Parse(r.ID)
}

// Response DTOs

type PaginatedResumesResponse struct {
	Data       []resume.ResumeSummaryResponse `json:"data"`
	Page       int                            `json:"page"`
	Limit      int                            `json:"limit"`
	Total      int                            `json:"total"`
	TotalPages int                            `json:"totalPages"`
}
