package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/recreatedev/Resumify/internal/errs"
	"github.com/recreatedev/Resumify/internal/model/education"
	"github.com/recreatedev/Resumify/internal/repository"
	"github.com/recreatedev/Resumify/internal/server"
)

type EducationService struct {
	server        *server.Server
	educationRepo *repository.EducationRepository
	resumeRepo    *repository.ResumeRepository
}

func NewEducationService(s *server.Server, repos *repository.Repositories) *EducationService {
	return &EducationService{
		server:        s,
		educationRepo: repos.Education,
		resumeRepo:    repos.Resume,
	}
}

// CreateEducation creates a new education entry
func (s *EducationService) CreateEducation(ctx context.Context, userID string, payload *education.CreateEducationRequest) (*education.EducationResponse, error) {
	// Verify resume belongs to user
	_, err := s.resumeRepo.GetResumeByID(ctx, userID, payload.ResumeID)
	if err != nil {
		if err.Error() == "failed to collect row from table:resumes" {
			return nil, errs.NewNotFoundError("resume not found", false, nil)
		}
		return nil, fmt.Errorf("failed to verify resume ownership: %w", err)
	}

	// Business logic: Validate date ranges
	if payload.StartDate != nil && payload.EndDate != nil {
		if payload.StartDate.After(*payload.EndDate) {
			return nil, errs.NewBadRequestError(
				"start date cannot be after end date",
				false, nil, nil, nil,
			)
		}
	}

	// Business logic: Check for duplicate education entries
	existingEducations, err := s.educationRepo.GetEducationByResumeID(ctx, userID, payload.ResumeID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing education: %w", err)
	}

	// Check for duplicates based on institution and degree
	for _, existing := range existingEducations {
		if existing.Institution == payload.Institution && existing.Degree == payload.Degree {
			return nil, errs.NewBadRequestError(
				"education entry with same institution and degree already exists",
				false, nil, nil, nil,
			)
		}
	}

	// Set default order index if not provided
	if payload.OrderIndex == 0 {
		payload.OrderIndex = len(existingEducations) + 1
	}

	// Create education in repository
	educationItem, err := s.educationRepo.CreateEducation(ctx, userID, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create education: %w", err)
	}

	// Convert to response DTO
	response := s.convertToEducationResponse(educationItem)

	return response, nil
}

// GetEducationByID retrieves an education entry by ID
func (s *EducationService) GetEducationByID(ctx context.Context, userID string, educationID uuid.UUID) (*education.EducationResponse, error) {
	educationItem, err := s.educationRepo.GetEducationByID(ctx, userID, educationID)
	if err != nil {
		if err.Error() == "failed to collect row from table:education" {
			return nil, errs.NewNotFoundError("education not found", false, nil)
		}
		return nil, fmt.Errorf("failed to get education: %w", err)
	}

	return s.convertToEducationResponse(educationItem), nil
}

// GetEducationByResumeID retrieves all education entries for a resume
func (s *EducationService) GetEducationByResumeID(ctx context.Context, userID string, resumeID uuid.UUID) ([]education.EducationResponse, error) {
	// Verify resume belongs to user
	_, err := s.resumeRepo.GetResumeByID(ctx, userID, resumeID)
	if err != nil {
		if err.Error() == "failed to collect row from table:resumes" {
			return nil, errs.NewNotFoundError("resume not found", false, nil)
		}
		return nil, fmt.Errorf("failed to verify resume ownership: %w", err)
	}

	educationItems, err := s.educationRepo.GetEducationByResumeID(ctx, userID, resumeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get education entries: %w", err)
	}

	// Convert to response DTOs
	responses := make([]education.EducationResponse, len(educationItems))
	for i, item := range educationItems {
		responses[i] = *s.convertToEducationResponse(&item)
	}

	return responses, nil
}

// UpdateEducation updates an education entry
func (s *EducationService) UpdateEducation(ctx context.Context, userID string, educationID uuid.UUID, payload *education.UpdateEducationRequest) (*education.EducationResponse, error) {
	// Check if education exists and belongs to user
	existingEducation, err := s.educationRepo.GetEducationByID(ctx, userID, educationID)
	if err != nil {
		if err.Error() == "failed to collect row from table:education" {
			return nil, errs.NewNotFoundError("education not found", false, nil)
		}
		return nil, fmt.Errorf("failed to get existing education: %w", err)
	}

	// Business logic: Validate date ranges
	if payload.StartDate != nil && payload.EndDate != nil {
		if payload.StartDate.After(*payload.EndDate) {
			return nil, errs.NewBadRequestError(
				"start date cannot be after end date",
				false, nil, nil, nil,
			)
		}
	} else if payload.StartDate != nil && existingEducation.EndDate != nil {
		if payload.StartDate.After(*existingEducation.EndDate) {
			return nil, errs.NewBadRequestError(
				"start date cannot be after end date",
				false, nil, nil, nil,
			)
		}
	} else if payload.EndDate != nil && existingEducation.StartDate != nil {
		if existingEducation.StartDate.After(*payload.EndDate) {
			return nil, errs.NewBadRequestError(
				"start date cannot be after end date",
				false, nil, nil, nil,
			)
		}
	}

	// Update education in repository
	updatedEducation, err := s.educationRepo.UpdateEducation(ctx, userID, educationID, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to update education: %w", err)
	}

	return s.convertToEducationResponse(updatedEducation), nil
}

// BulkUpdateEducationOrder updates the order of multiple education entries
func (s *EducationService) BulkUpdateEducationOrder(ctx context.Context, userID string, payload *education.BulkUpdateEducationRequest) error {
	// Validate that all education entries belong to the user
	for _, eduUpdate := range payload.Education {
		educationID, err := uuid.Parse(eduUpdate.ID)
		if err != nil {
			return errs.NewBadRequestError("invalid education ID", false, nil, nil, nil)
		}
		_, err = s.educationRepo.GetEducationByID(ctx, userID, educationID)
		if err != nil {
			if err.Error() == "failed to collect row from table:education" {
				return errs.NewNotFoundError("education not found", false, nil)
			}
			return fmt.Errorf("failed to verify education ownership: %w", err)
		}
	}

	// Update order in repository
	err := s.educationRepo.BulkUpdateEducationOrder(ctx, userID, payload)
	if err != nil {
		return fmt.Errorf("failed to update education order: %w", err)
	}

	return nil
}

// DeleteEducation deletes an education entry
func (s *EducationService) DeleteEducation(ctx context.Context, userID string, educationID uuid.UUID) error {
	// Check if education exists and belongs to user
	_, err := s.educationRepo.GetEducationByID(ctx, userID, educationID)
	if err != nil {
		if err.Error() == "failed to collect row from table:education" {
			return errs.NewNotFoundError("education not found", false, nil)
		}
		return fmt.Errorf("failed to get existing education: %w", err)
	}

	// Delete education
	err = s.educationRepo.DeleteEducation(ctx, userID, educationID)
	if err != nil {
		return fmt.Errorf("failed to delete education: %w", err)
	}

	return nil
}

// Helper methods

func (s *EducationService) convertToEducationResponse(educationItem *education.Education) *education.EducationResponse {
	response := &education.EducationResponse{
		ID:           educationItem.ID.String(),
		ResumeID:     educationItem.ResumeID,
		Institution:  educationItem.Institution,
		Degree:       educationItem.Degree,
		FieldOfStudy: educationItem.FieldOfStudy,
		Grade:        educationItem.Grade,
		Description:  educationItem.Description,
		OrderIndex:   educationItem.OrderIndex,
		CreatedAt:    educationItem.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    educationItem.UpdatedAt.Format(time.RFC3339),
	}

	// Handle optional date fields
	if educationItem.StartDate != nil {
		response.StartDate = educationItem.StartDate
	}
	if educationItem.EndDate != nil {
		response.EndDate = educationItem.EndDate
	}

	return response
}
