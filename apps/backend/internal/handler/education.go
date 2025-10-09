package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/recreatedev/Resumify/internal/middleware"
	"github.com/recreatedev/Resumify/internal/model/education"
	"github.com/recreatedev/Resumify/internal/server"
	"github.com/recreatedev/Resumify/internal/service"
)

type EducationHandler struct {
	Handler
	educationService *service.EducationService
}

func NewEducationHandler(s *server.Server, educationService *service.EducationService) *EducationHandler {
	return &EducationHandler{
		Handler:          NewHandler(s),
		educationService: educationService,
	}
}

func (h *EducationHandler) CreateEducation(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *education.CreateEducationRequest) (*education.EducationResponse, error) {
			userID := middleware.GetUserID(c)
			return h.educationService.CreateEducation(c.Request().Context(), userID, payload)
		},
		http.StatusCreated,
		&education.CreateEducationRequest{},
	)(c)
}

func (h *EducationHandler) GetEducationByID(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, req *GetEducationByIDRequest) (*education.EducationResponse, error) {
			userID := middleware.GetUserID(c)
			educationID, err := req.ParseID()
			if err != nil {
				return nil, err
			}
			return h.educationService.GetEducationByID(c.Request().Context(), userID, educationID)
		},
		http.StatusOK,
		&GetEducationByIDRequest{},
	)(c)
}

func (h *EducationHandler) GetEducationByResumeID(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, req *GetEducationByResumeIDRequest) ([]education.EducationResponse, error) {
			userID := middleware.GetUserID(c)
			resumeID, err := req.ParseResumeID()
			if err != nil {
				return nil, err
			}
			return h.educationService.GetEducationByResumeID(c.Request().Context(), userID, resumeID)
		},
		http.StatusOK,
		&GetEducationByResumeIDRequest{},
	)(c)
}

func (h *EducationHandler) UpdateEducation(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, req *UpdateEducationRequest) (*education.EducationResponse, error) {
			userID := middleware.GetUserID(c)
			educationID, err := req.ParseID()
			if err != nil {
				return nil, err
			}
			return h.educationService.UpdateEducation(c.Request().Context(), userID, educationID, req.UpdateEducationRequest)
		},
		http.StatusOK,
		&UpdateEducationRequest{},
	)(c)
}

func (h *EducationHandler) BulkUpdateEducationOrder(c echo.Context) error {
	return HandleNoContent(
		h.Handler,
		func(c echo.Context, payload *education.BulkUpdateEducationRequest) error {
			userID := middleware.GetUserID(c)
			return h.educationService.BulkUpdateEducationOrder(c.Request().Context(), userID, payload)
		},
		http.StatusNoContent,
		&education.BulkUpdateEducationRequest{},
	)(c)
}

func (h *EducationHandler) DeleteEducation(c echo.Context) error {
	return HandleNoContent(
		h.Handler,
		func(c echo.Context, req *DeleteEducationRequest) error {
			userID := middleware.GetUserID(c)
			educationID, err := req.ParseID()
			if err != nil {
				return err
			}
			return h.educationService.DeleteEducation(c.Request().Context(), userID, educationID)
		},
		http.StatusNoContent,
		&DeleteEducationRequest{},
	)(c)
}

// Request DTOs

type GetEducationByIDRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}

func (r *GetEducationByIDRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *GetEducationByIDRequest) ParseID() (uuid.UUID, error) {
	return uuid.Parse(r.ID)
}

type GetEducationByResumeIDRequest struct {
	ResumeID string `param:"resumeId" validate:"required,uuid"`
}

func (r *GetEducationByResumeIDRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *GetEducationByResumeIDRequest) ParseResumeID() (uuid.UUID, error) {
	return uuid.Parse(r.ResumeID)
}

type UpdateEducationRequest struct {
	ID string `param:"id" validate:"required,uuid"`
	*education.UpdateEducationRequest
}

func (r *UpdateEducationRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err
	}
	return r.UpdateEducationRequest.Validate()
}

func (r *UpdateEducationRequest) ParseID() (uuid.UUID, error) {
	return uuid.Parse(r.ID)
}

type DeleteEducationRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}

func (r *DeleteEducationRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *DeleteEducationRequest) ParseID() (uuid.UUID, error) {
	return uuid.Parse(r.ID)
}
