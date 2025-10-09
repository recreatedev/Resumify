package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/recreatedev/Resumify/internal/middleware"
	"github.com/recreatedev/Resumify/internal/model/section"
	"github.com/recreatedev/Resumify/internal/server"
	"github.com/recreatedev/Resumify/internal/service"
)

type SectionHandler struct {
	Handler
	sectionService *service.SectionService
}

func NewSectionHandler(s *server.Server, sectionService *service.SectionService) *SectionHandler {
	return &SectionHandler{
		Handler:        NewHandler(s),
		sectionService: sectionService,
	}
}

func (h *SectionHandler) CreateSection(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *section.CreateSectionRequest) (*section.SectionResponse, error) {
			userID := middleware.GetUserID(c)
			return h.sectionService.CreateSection(c.Request().Context(), userID, payload)
		},
		http.StatusCreated,
		&section.CreateSectionRequest{},
	)(c)
}

func (h *SectionHandler) GetSectionByID(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, req *GetSectionByIDRequest) (*section.SectionResponse, error) {
			userID := middleware.GetUserID(c)
			sectionID, err := req.ParseID()
			if err != nil {
				return nil, err
			}
			return h.sectionService.GetSectionByID(c.Request().Context(), userID, sectionID)
		},
		http.StatusOK,
		&GetSectionByIDRequest{},
	)(c)
}

func (h *SectionHandler) GetSectionsByResumeID(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, req *GetSectionsByResumeIDRequest) ([]section.SectionResponse, error) {
			userID := middleware.GetUserID(c)
			resumeID, err := req.ParseResumeID()
			if err != nil {
				return nil, err
			}
			return h.sectionService.GetSectionsByResumeID(c.Request().Context(), userID, resumeID)
		},
		http.StatusOK,
		&GetSectionsByResumeIDRequest{},
	)(c)
}

func (h *SectionHandler) UpdateSection(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, req *UpdateSectionRequest) (*section.SectionResponse, error) {
			userID := middleware.GetUserID(c)
			sectionID, err := req.ParseID()
			if err != nil {
				return nil, err
			}
			return h.sectionService.UpdateSection(c.Request().Context(), userID, sectionID, req.UpdateSectionRequest)
		},
		http.StatusOK,
		&UpdateSectionRequest{},
	)(c)
}

func (h *SectionHandler) BulkUpdateSectionOrder(c echo.Context) error {
	return HandleNoContent(
		h.Handler,
		func(c echo.Context, payload *section.BulkUpdateSectionsRequest) error {
			userID := middleware.GetUserID(c)
			return h.sectionService.BulkUpdateSectionOrder(c.Request().Context(), userID, payload)
		},
		http.StatusNoContent,
		&section.BulkUpdateSectionsRequest{},
	)(c)
}

func (h *SectionHandler) DeleteSection(c echo.Context) error {
	return HandleNoContent(
		h.Handler,
		func(c echo.Context, req *DeleteSectionRequest) error {
			userID := middleware.GetUserID(c)
			sectionID, err := req.ParseID()
			if err != nil {
				return err
			}
			return h.sectionService.DeleteSection(c.Request().Context(), userID, sectionID)
		},
		http.StatusNoContent,
		&DeleteSectionRequest{},
	)(c)
}

// Request DTOs

type GetSectionByIDRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}

func (r *GetSectionByIDRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *GetSectionByIDRequest) ParseID() (uuid.UUID, error) {
	return uuid.Parse(r.ID)
}

type GetSectionsByResumeIDRequest struct {
	ResumeID string `param:"resumeId" validate:"required,uuid"`
}

func (r *GetSectionsByResumeIDRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *GetSectionsByResumeIDRequest) ParseResumeID() (uuid.UUID, error) {
	return uuid.Parse(r.ResumeID)
}

type UpdateSectionRequest struct {
	ID string `param:"id" validate:"required,uuid"`
	*section.UpdateSectionRequest
}

func (r *UpdateSectionRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err
	}
	return r.UpdateSectionRequest.Validate()
}

func (r *UpdateSectionRequest) ParseID() (uuid.UUID, error) {
	return uuid.Parse(r.ID)
}

type DeleteSectionRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}

func (r *DeleteSectionRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *DeleteSectionRequest) ParseID() (uuid.UUID, error) {
	return uuid.Parse(r.ID)
}
