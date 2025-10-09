package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/recreatedev/Resumify/internal/middleware"
	"github.com/recreatedev/Resumify/internal/model/skill"
	"github.com/recreatedev/Resumify/internal/server"
	"github.com/recreatedev/Resumify/internal/service"
)

type SkillHandler struct {
	Handler
	skillService *service.SkillService
}

func NewSkillHandler(s *server.Server, skillService *service.SkillService) *SkillHandler {
	return &SkillHandler{
		Handler:      NewHandler(s),
		skillService: skillService,
	}
}

func (h *SkillHandler) CreateSkill(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *skill.CreateSkillRequest) (*skill.SkillResponse, error) {
			userID := middleware.GetUserID(c)
			return h.skillService.CreateSkill(c.Request().Context(), userID, payload)
		},
		http.StatusCreated,
		&skill.CreateSkillRequest{},
	)(c)
}

func (h *SkillHandler) GetSkillByID(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, req *GetSkillByIDRequest) (*skill.SkillResponse, error) {
			userID := middleware.GetUserID(c)
			skillID, err := req.ParseID()
			if err != nil {
				return nil, err
			}
			return h.skillService.GetSkillByID(c.Request().Context(), userID, skillID)
		},
		http.StatusOK,
		&GetSkillByIDRequest{},
	)(c)
}

func (h *SkillHandler) GetSkillsByResumeID(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, req *GetSkillsByResumeIDRequest) ([]skill.SkillResponse, error) {
			userID := middleware.GetUserID(c)
			resumeID, err := req.ParseResumeID()
			if err != nil {
				return nil, err
			}
			return h.skillService.GetSkillsByResumeID(c.Request().Context(), userID, resumeID)
		},
		http.StatusOK,
		&GetSkillsByResumeIDRequest{},
	)(c)
}

func (h *SkillHandler) GetSkillsByCategory(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, req *GetSkillsByCategoryRequest) ([]skill.SkillsByCategoryResponse, error) {
			userID := middleware.GetUserID(c)
			resumeID, err := req.ParseResumeID()
			if err != nil {
				return nil, err
			}
			return h.skillService.GetSkillsByCategory(c.Request().Context(), userID, resumeID)
		},
		http.StatusOK,
		&GetSkillsByCategoryRequest{},
	)(c)
}

func (h *SkillHandler) UpdateSkill(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, req *UpdateSkillRequest) (*skill.SkillResponse, error) {
			userID := middleware.GetUserID(c)
			skillID, err := req.ParseID()
			if err != nil {
				return nil, err
			}
			return h.skillService.UpdateSkill(c.Request().Context(), userID, skillID, req.UpdateSkillRequest)
		},
		http.StatusOK,
		&UpdateSkillRequest{},
	)(c)
}

func (h *SkillHandler) BulkUpdateSkillOrder(c echo.Context) error {
	return HandleNoContent(
		h.Handler,
		func(c echo.Context, payload *skill.BulkUpdateSkillsRequest) error {
			userID := middleware.GetUserID(c)
			return h.skillService.BulkUpdateSkillOrder(c.Request().Context(), userID, payload)
		},
		http.StatusNoContent,
		&skill.BulkUpdateSkillsRequest{},
	)(c)
}

func (h *SkillHandler) DeleteSkill(c echo.Context) error {
	return HandleNoContent(
		h.Handler,
		func(c echo.Context, req *DeleteSkillRequest) error {
			userID := middleware.GetUserID(c)
			skillID, err := req.ParseID()
			if err != nil {
				return err
			}
			return h.skillService.DeleteSkill(c.Request().Context(), userID, skillID)
		},
		http.StatusNoContent,
		&DeleteSkillRequest{},
	)(c)
}

// Request DTOs

type GetSkillByIDRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}

func (r *GetSkillByIDRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *GetSkillByIDRequest) ParseID() (uuid.UUID, error) {
	return uuid.Parse(r.ID)
}

type GetSkillsByResumeIDRequest struct {
	ResumeID string `param:"resumeId" validate:"required,uuid"`
}

func (r *GetSkillsByResumeIDRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *GetSkillsByResumeIDRequest) ParseResumeID() (uuid.UUID, error) {
	return uuid.Parse(r.ResumeID)
}

type GetSkillsByCategoryRequest struct {
	ResumeID string `param:"resumeId" validate:"required,uuid"`
}

func (r *GetSkillsByCategoryRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *GetSkillsByCategoryRequest) ParseResumeID() (uuid.UUID, error) {
	return uuid.Parse(r.ResumeID)
}

type UpdateSkillRequest struct {
	ID string `param:"id" validate:"required,uuid"`
	*skill.UpdateSkillRequest
}

func (r *UpdateSkillRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err
	}
	return r.UpdateSkillRequest.Validate()
}

func (r *UpdateSkillRequest) ParseID() (uuid.UUID, error) {
	return uuid.Parse(r.ID)
}

type DeleteSkillRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}

func (r *DeleteSkillRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *DeleteSkillRequest) ParseID() (uuid.UUID, error) {
	return uuid.Parse(r.ID)
}
