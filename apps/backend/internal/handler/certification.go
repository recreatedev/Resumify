package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/recreatedev/Resumify/internal/middleware"
	"github.com/recreatedev/Resumify/internal/model/certification"
	"github.com/recreatedev/Resumify/internal/server"
	"github.com/recreatedev/Resumify/internal/service"
)

type CertificationHandler struct {
	Handler
	certificationService *service.CertificationService
}

func NewCertificationHandler(s *server.Server, certificationService *service.CertificationService) *CertificationHandler {
	return &CertificationHandler{
		Handler:              NewHandler(s),
		certificationService: certificationService,
	}
}

func (h *CertificationHandler) CreateCertification(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *certification.CreateCertificationRequest) (*certification.CertificationResponse, error) {
			userID := middleware.GetUserID(c)
			return h.certificationService.CreateCertification(c.Request().Context(), userID, payload)
		},
		http.StatusCreated,
		&certification.CreateCertificationRequest{},
	)(c)
}

func (h *CertificationHandler) GetCertificationByID(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, req *GetCertificationByIDRequest) (*certification.CertificationResponse, error) {
			userID := middleware.GetUserID(c)
			certificationID, err := req.ParseID()
			if err != nil {
				return nil, err
			}
			return h.certificationService.GetCertificationByID(c.Request().Context(), userID, certificationID)
		},
		http.StatusOK,
		&GetCertificationByIDRequest{},
	)(c)
}

func (h *CertificationHandler) GetCertificationsByResumeID(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, req *GetCertificationsByResumeIDRequest) ([]certification.CertificationResponse, error) {
			userID := middleware.GetUserID(c)
			resumeID, err := req.ParseResumeID()
			if err != nil {
				return nil, err
			}
			return h.certificationService.GetCertificationsByResumeID(c.Request().Context(), userID, resumeID)
		},
		http.StatusOK,
		&GetCertificationsByResumeIDRequest{},
	)(c)
}

func (h *CertificationHandler) UpdateCertification(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, req *UpdateCertificationRequest) (*certification.CertificationResponse, error) {
			userID := middleware.GetUserID(c)
			certificationID, err := req.ParseID()
			if err != nil {
				return nil, err
			}
			return h.certificationService.UpdateCertification(c.Request().Context(), userID, certificationID, req.UpdateCertificationRequest)
		},
		http.StatusOK,
		&UpdateCertificationRequest{},
	)(c)
}

func (h *CertificationHandler) BulkUpdateCertificationOrder(c echo.Context) error {
	return HandleNoContent(
		h.Handler,
		func(c echo.Context, payload *certification.BulkUpdateCertificationsRequest) error {
			userID := middleware.GetUserID(c)
			return h.certificationService.BulkUpdateCertificationOrder(c.Request().Context(), userID, payload)
		},
		http.StatusNoContent,
		&certification.BulkUpdateCertificationsRequest{},
	)(c)
}

func (h *CertificationHandler) DeleteCertification(c echo.Context) error {
	return HandleNoContent(
		h.Handler,
		func(c echo.Context, req *DeleteCertificationRequest) error {
			userID := middleware.GetUserID(c)
			certificationID, err := req.ParseID()
			if err != nil {
				return err
			}
			return h.certificationService.DeleteCertification(c.Request().Context(), userID, certificationID)
		},
		http.StatusNoContent,
		&DeleteCertificationRequest{},
	)(c)
}

// Request DTOs

type GetCertificationByIDRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}

func (r *GetCertificationByIDRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *GetCertificationByIDRequest) ParseID() (uuid.UUID, error) {
	return uuid.Parse(r.ID)
}

type GetCertificationsByResumeIDRequest struct {
	ResumeID string `param:"resumeId" validate:"required,uuid"`
}

func (r *GetCertificationsByResumeIDRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *GetCertificationsByResumeIDRequest) ParseResumeID() (uuid.UUID, error) {
	return uuid.Parse(r.ResumeID)
}

type UpdateCertificationRequest struct {
	ID string `param:"id" validate:"required,uuid"`
	*certification.UpdateCertificationRequest
}

func (r *UpdateCertificationRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err
	}
	return r.UpdateCertificationRequest.Validate()
}

func (r *UpdateCertificationRequest) ParseID() (uuid.UUID, error) {
	return uuid.Parse(r.ID)
}

type DeleteCertificationRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}

func (r *DeleteCertificationRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *DeleteCertificationRequest) ParseID() (uuid.UUID, error) {
	return uuid.Parse(r.ID)
}
