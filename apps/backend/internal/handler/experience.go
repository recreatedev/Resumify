package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/recreatedev/Resumify/internal/middleware"
	"github.com/recreatedev/Resumify/internal/model/experience"
	"github.com/recreatedev/Resumify/internal/server"
	"github.com/recreatedev/Resumify/internal/service"
)

type ExperienceHandler struct {
	Handler
	experienceService *service.ExperienceService
}

func NewExperienceHandler(s *server.Server, experienceService *service.ExperienceService) *ExperienceHandler {
	return &ExperienceHandler{
		Handler:           NewHandler(s),
		experienceService: experienceService,
	}
}

func (h *ExperienceHandler) CreateExperience(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *experience.CreateExperienceRequest) (*experience.ExperienceResponse, error) {
			userID := middleware.GetUserID(c)
			return h.experienceService.CreateExperience(c.Request().Context(), userID, payload)
		},
		http.StatusCreated,
		&experience.CreateExperienceRequest{},
	)(c)
}

func (h *ExperienceHandler) GetExperienceByID(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, req *GetExperienceByIDRequest) (*experience.ExperienceResponse, error) {
			userID := middleware.GetUserID(c)
			experienceID, err := req.ParseID()
			if err != nil {
				return nil, err
			}
			return h.experienceService.GetExperienceByID(c.Request().Context(), userID, experienceID)
		},
		http.StatusOK,
		&GetExperienceByIDRequest{},
	)(c)
}

func (h *ExperienceHandler) GetExperienceByResumeID(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, req *GetExperienceByResumeIDRequest) ([]experience.ExperienceResponse, error) {
			userID := middleware.GetUserID(c)
			resumeID, err := req.ParseResumeID()
			if err != nil {
				return nil, err
			}
			return h.experienceService.GetExperienceByResumeID(c.Request().Context(), userID, resumeID)
		},
		http.StatusOK,
		&GetExperienceByResumeIDRequest{},
	)(c)
}

func (h *ExperienceHandler) UpdateExperience(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, req *UpdateExperienceRequest) (*experience.ExperienceResponse, error) {
			userID := middleware.GetUserID(c)
			experienceID, err := req.ParseID()
			if err != nil {
				return nil, err
			}
			return h.experienceService.UpdateExperience(c.Request().Context(), userID, experienceID, req.UpdateExperienceRequest)
		},
		http.StatusOK,
		&UpdateExperienceRequest{},
	)(c)
}

func (h *ExperienceHandler) BulkUpdateExperienceOrder(c echo.Context) error {
	return HandleNoContent(
		h.Handler,
		func(c echo.Context, payload *experience.BulkUpdateExperienceRequest) error {
			userID := middleware.GetUserID(c)
			return h.experienceService.BulkUpdateExperienceOrder(c.Request().Context(), userID, payload)
		},
		http.StatusNoContent,
		&experience.BulkUpdateExperienceRequest{},
	)(c)
}

func (h *ExperienceHandler) DeleteExperience(c echo.Context) error {
	return HandleNoContent(
		h.Handler,
		func(c echo.Context, req *DeleteExperienceRequest) error {
			userID := middleware.GetUserID(c)
			experienceID, err := req.ParseID()
			if err != nil {
				return err
			}
			return h.experienceService.DeleteExperience(c.Request().Context(), userID, experienceID)
		},
		http.StatusNoContent,
		&DeleteExperienceRequest{},
	)(c)
}

// Request DTOs

type GetExperienceByIDRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}

func (r *GetExperienceByIDRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *GetExperienceByIDRequest) ParseID() (uuid.UUID, error) {
	return uuid.Parse(r.ID)
}

type GetExperienceByResumeIDRequest struct {
	ResumeID string `param:"resumeId" validate:"required,uuid"`
}

func (r *GetExperienceByResumeIDRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *GetExperienceByResumeIDRequest) ParseResumeID() (uuid.UUID, error) {
	return uuid.Parse(r.ResumeID)
}

type UpdateExperienceRequest struct {
	ID string `param:"id" validate:"required,uuid"`
	*experience.UpdateExperienceRequest
}

func (r *UpdateExperienceRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err
	}
	return r.UpdateExperienceRequest.Validate()
}

func (r *UpdateExperienceRequest) ParseID() (uuid.UUID, error) {
	return uuid.Parse(r.ID)
}

type DeleteExperienceRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}

func (r *DeleteExperienceRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *DeleteExperienceRequest) ParseID() (uuid.UUID, error) {
	return uuid.Parse(r.ID)
}
