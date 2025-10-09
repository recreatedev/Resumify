package service

import (
	"context"
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/recreatedev/Resumify/internal/errs"
	"github.com/recreatedev/Resumify/internal/model/certification"
	"github.com/recreatedev/Resumify/internal/repository"
	"github.com/recreatedev/Resumify/internal/server"
)

type CertificationService struct {
	server            *server.Server
	certificationRepo *repository.CertificationRepository
	resumeRepo        *repository.ResumeRepository
}

func NewCertificationService(s *server.Server, repos *repository.Repositories) *CertificationService {
	return &CertificationService{
		server:            s,
		certificationRepo: repos.Certification,
		resumeRepo:        repos.Resume,
	}
}

// CreateCertification creates a new certification entry
func (s *CertificationService) CreateCertification(ctx context.Context, userID string, payload *certification.CreateCertificationRequest) (*certification.CertificationResponse, error) {
	// Verify resume belongs to user
	_, err := s.resumeRepo.GetResumeByID(ctx, userID, payload.ResumeID)
	if err != nil {
		if err.Error() == "failed to collect row from table:resumes" {
			return nil, errs.NewNotFoundError("resume not found", false, nil)
		}
		return nil, fmt.Errorf("failed to verify resume ownership: %w", err)
	}

	// Business logic: Validate date ranges
	if payload.IssueDate != nil && payload.ExpiryDate != nil {
		if payload.IssueDate.After(*payload.ExpiryDate) {
			return nil, errs.NewBadRequestError(
				"issue date cannot be after expiry date",
				false, nil, nil, nil,
			)
		}
	}

	// Business logic: Validate credential URL if provided
	if payload.CredentialURL != nil && *payload.CredentialURL != "" {
		if _, err := url.Parse(*payload.CredentialURL); err != nil {
			return nil, errs.NewBadRequestError(
				"invalid credential URL",
				false, nil, nil, nil,
			)
		}
	}

	// Business logic: Check for duplicate certification entries
	existingCertifications, err := s.certificationRepo.GetCertificationsByResumeID(ctx, userID, payload.ResumeID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing certifications: %w", err)
	}

	// Check for duplicates based on name and organization
	for _, existing := range existingCertifications {
		if *existing.Name == *payload.Name && *existing.Organization == *payload.Organization {
			return nil, errs.NewBadRequestError(
				"certification with same name and organization already exists",
				false, nil, nil, nil,
			)
		}
	}

	// Set default order index if not provided
	if payload.OrderIndex == 0 {
		payload.OrderIndex = len(existingCertifications) + 1
	}

	// Create certification in repository
	certificationItem, err := s.certificationRepo.CreateCertification(ctx, userID, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create certification: %w", err)
	}

	// Convert to response DTO
	response := s.convertToCertificationResponse(certificationItem)

	return response, nil
}

// GetCertificationByID retrieves a certification entry by ID
func (s *CertificationService) GetCertificationByID(ctx context.Context, userID string, certificationID uuid.UUID) (*certification.CertificationResponse, error) {
	certificationItem, err := s.certificationRepo.GetCertificationByID(ctx, userID, certificationID)
	if err != nil {
		if err.Error() == "failed to collect row from table:certifications" {
			return nil, errs.NewNotFoundError("certification not found", false, nil)
		}
		return nil, fmt.Errorf("failed to get certification: %w", err)
	}

	return s.convertToCertificationResponse(certificationItem), nil
}

// GetCertificationsByResumeID retrieves all certification entries for a resume
func (s *CertificationService) GetCertificationsByResumeID(ctx context.Context, userID string, resumeID uuid.UUID) ([]certification.CertificationResponse, error) {
	// Verify resume belongs to user
	_, err := s.resumeRepo.GetResumeByID(ctx, userID, resumeID)
	if err != nil {
		if err.Error() == "failed to collect row from table:resumes" {
			return nil, errs.NewNotFoundError("resume not found", false, nil)
		}
		return nil, fmt.Errorf("failed to verify resume ownership: %w", err)
	}

	certificationItems, err := s.certificationRepo.GetCertificationsByResumeID(ctx, userID, resumeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get certification entries: %w", err)
	}

	// Convert to response DTOs
	responses := make([]certification.CertificationResponse, len(certificationItems))
	for i, item := range certificationItems {
		responses[i] = *s.convertToCertificationResponse(&item)
	}

	return responses, nil
}

// UpdateCertification updates a certification entry
func (s *CertificationService) UpdateCertification(ctx context.Context, userID string, certificationID uuid.UUID, payload *certification.UpdateCertificationRequest) (*certification.CertificationResponse, error) {
	// Check if certification exists and belongs to user
	existingCertification, err := s.certificationRepo.GetCertificationByID(ctx, userID, certificationID)
	if err != nil {
		if err.Error() == "failed to collect row from table:certifications" {
			return nil, errs.NewNotFoundError("certification not found", false, nil)
		}
		return nil, fmt.Errorf("failed to get existing certification: %w", err)
	}

	// Business logic: Validate date ranges
	if payload.IssueDate != nil && payload.ExpiryDate != nil {
		if payload.IssueDate.After(*payload.ExpiryDate) {
			return nil, errs.NewBadRequestError(
				"issue date cannot be after expiry date",
				false, nil, nil, nil,
			)
		}
	} else if payload.IssueDate != nil && existingCertification.ExpiryDate != nil {
		if payload.IssueDate.After(*existingCertification.ExpiryDate) {
			return nil, errs.NewBadRequestError(
				"issue date cannot be after expiry date",
				false, nil, nil, nil,
			)
		}
	} else if payload.ExpiryDate != nil && existingCertification.IssueDate != nil {
		if existingCertification.IssueDate.After(*payload.ExpiryDate) {
			return nil, errs.NewBadRequestError(
				"issue date cannot be after expiry date",
				false, nil, nil, nil,
			)
		}
	}

	// Business logic: Validate credential URL if provided
	if payload.CredentialURL != nil && *payload.CredentialURL != "" {
		if _, err := url.Parse(*payload.CredentialURL); err != nil {
			return nil, errs.NewBadRequestError(
				"invalid credential URL",
				false, nil, nil, nil,
			)
		}
	}

	// Business logic: Check for duplicate certification names (excluding current certification)
	if (payload.Name != nil && *payload.Name != *existingCertification.Name) ||
		(payload.Organization != nil && *payload.Organization != *existingCertification.Organization) {

		certifications, err := s.certificationRepo.GetCertificationsByResumeID(ctx, userID, existingCertification.ResumeID)
		if err != nil {
			return nil, fmt.Errorf("failed to check existing certifications: %w", err)
		}

		name := *existingCertification.Name
		organization := *existingCertification.Organization

		if payload.Name != nil {
			name = *payload.Name
		}
		if payload.Organization != nil {
			organization = *payload.Organization
		}

		for _, cert := range certifications {
			if cert.ID != certificationID && *cert.Name == name && *cert.Organization == organization {
				return nil, errs.NewBadRequestError(
					"certification with same name and organization already exists",
					false, nil, nil, nil,
				)
			}
		}
	}

	// Update certification in repository
	updatedCertification, err := s.certificationRepo.UpdateCertification(ctx, userID, certificationID, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to update certification: %w", err)
	}

	return s.convertToCertificationResponse(updatedCertification), nil
}

// BulkUpdateCertificationOrder updates the order of multiple certification entries
func (s *CertificationService) BulkUpdateCertificationOrder(ctx context.Context, userID string, payload *certification.BulkUpdateCertificationsRequest) error {
	// Validate that all certification entries belong to the user
	for _, certUpdate := range payload.Certifications {
		certificationID, err := uuid.Parse(certUpdate.ID)
		if err != nil {
			return errs.NewBadRequestError("invalid certification ID", false, nil, nil, nil)
		}
		_, err = s.certificationRepo.GetCertificationByID(ctx, userID, certificationID)
		if err != nil {
			if err.Error() == "failed to collect row from table:certifications" {
				return errs.NewNotFoundError("certification not found", false, nil)
			}
			return fmt.Errorf("failed to verify certification ownership: %w", err)
		}
	}

	// Update order in repository
	err := s.certificationRepo.BulkUpdateCertificationOrder(ctx, userID, payload)
	if err != nil {
		return fmt.Errorf("failed to update certification order: %w", err)
	}

	return nil
}

// DeleteCertification deletes a certification entry
func (s *CertificationService) DeleteCertification(ctx context.Context, userID string, certificationID uuid.UUID) error {
	// Check if certification exists and belongs to user
	_, err := s.certificationRepo.GetCertificationByID(ctx, userID, certificationID)
	if err != nil {
		if err.Error() == "failed to collect row from table:certifications" {
			return errs.NewNotFoundError("certification not found", false, nil)
		}
		return fmt.Errorf("failed to get existing certification: %w", err)
	}

	// Delete certification
	err = s.certificationRepo.DeleteCertification(ctx, userID, certificationID)
	if err != nil {
		return fmt.Errorf("failed to delete certification: %w", err)
	}

	return nil
}

// Helper methods

func (s *CertificationService) convertToCertificationResponse(certificationItem *certification.Certification) *certification.CertificationResponse {
	response := &certification.CertificationResponse{
		ID:            certificationItem.ID.String(),
		ResumeID:      certificationItem.ResumeID,
		Name:          certificationItem.Name,
		Organization:  certificationItem.Organization,
		CredentialID:  certificationItem.CredentialID,
		CredentialURL: certificationItem.CredentialURL,
		OrderIndex:    certificationItem.OrderIndex,
	}

	// Handle optional date fields
	if certificationItem.IssueDate != nil {
		response.IssueDate = certificationItem.IssueDate
	}
	if certificationItem.ExpiryDate != nil {
		response.ExpiryDate = certificationItem.ExpiryDate
	}

	return response
}
