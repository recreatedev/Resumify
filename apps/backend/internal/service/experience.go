package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/recreatedev/Resumify/internal/errs"
	"github.com/recreatedev/Resumify/internal/model/experience"
	"github.com/recreatedev/Resumify/internal/repository"
	"github.com/recreatedev/Resumify/internal/server"
)

type ExperienceService struct {
	server         *server.Server
	experienceRepo *repository.ExperienceRepository
	resumeRepo     *repository.ResumeRepository
}

func NewExperienceService(s *server.Server, repos *repository.Repositories) *ExperienceService {
	return &ExperienceService{
		server:         s,
		experienceRepo: repos.Experience,
		resumeRepo:     repos.Resume,
	}
}

// CreateExperience creates a new experience entry
func (s *ExperienceService) CreateExperience(ctx context.Context, userID string, payload *experience.CreateExperienceRequest) (*experience.ExperienceResponse, error) {
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

	// Business logic: Check for duplicate experience entries
	existingExperiences, err := s.experienceRepo.GetExperienceByResumeID(ctx, userID, payload.ResumeID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing experience: %w", err)
	}

	// Check for duplicates based on company and position
	for _, existing := range existingExperiences {
		if existing.Company == payload.Company && existing.Position == payload.Position {
			return nil, errs.NewBadRequestError(
				"experience entry with same company and position already exists",
				false, nil, nil, nil,
			)
		}
	}

	// Set default order index if not provided
	if payload.OrderIndex == 0 {
		payload.OrderIndex = len(existingExperiences) + 1
	}

	// Create experience in repository
	experienceItem, err := s.experienceRepo.CreateExperience(ctx, userID, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create experience: %w", err)
	}

	// Convert to response DTO
	response := s.convertToExperienceResponse(experienceItem)

	return response, nil
}

// GetExperienceByID retrieves an experience entry by ID
func (s *ExperienceService) GetExperienceByID(ctx context.Context, userID string, experienceID uuid.UUID) (*experience.ExperienceResponse, error) {
	experienceItem, err := s.experienceRepo.GetExperienceByID(ctx, userID, experienceID)
	if err != nil {
		if err.Error() == "failed to collect row from table:experience" {
			return nil, errs.NewNotFoundError("experience not found", false, nil)
		}
		return nil, fmt.Errorf("failed to get experience: %w", err)
	}

	return s.convertToExperienceResponse(experienceItem), nil
}

// GetExperienceByResumeID retrieves all experience entries for a resume
func (s *ExperienceService) GetExperienceByResumeID(ctx context.Context, userID string, resumeID uuid.UUID) ([]experience.ExperienceResponse, error) {
	// Verify resume belongs to user
	_, err := s.resumeRepo.GetResumeByID(ctx, userID, resumeID)
	if err != nil {
		if err.Error() == "failed to collect row from table:resumes" {
			return nil, errs.NewNotFoundError("resume not found", false, nil)
		}
		return nil, fmt.Errorf("failed to verify resume ownership: %w", err)
	}

	experienceItems, err := s.experienceRepo.GetExperienceByResumeID(ctx, userID, resumeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get experience entries: %w", err)
	}

	// Convert to response DTOs
	responses := make([]experience.ExperienceResponse, len(experienceItems))
	for i, item := range experienceItems {
		responses[i] = *s.convertToExperienceResponse(&item)
	}

	return responses, nil
}

// UpdateExperience updates an experience entry
func (s *ExperienceService) UpdateExperience(ctx context.Context, userID string, experienceID uuid.UUID, payload *experience.UpdateExperienceRequest) (*experience.ExperienceResponse, error) {
	// Check if experience exists and belongs to user
	existingExperience, err := s.experienceRepo.GetExperienceByID(ctx, userID, experienceID)
	if err != nil {
		if err.Error() == "failed to collect row from table:experience" {
			return nil, errs.NewNotFoundError("experience not found", false, nil)
		}
		return nil, fmt.Errorf("failed to get existing experience: %w", err)
	}

	// Business logic: Validate date ranges
	if payload.StartDate != nil && payload.EndDate != nil {
		if payload.StartDate.After(*payload.EndDate) {
			return nil, errs.NewBadRequestError(
				"start date cannot be after end date",
				false, nil, nil, nil,
			)
		}
	} else if payload.StartDate != nil && existingExperience.EndDate != nil {
		if payload.StartDate.After(*existingExperience.EndDate) {
			return nil, errs.NewBadRequestError(
				"start date cannot be after end date",
				false, nil, nil, nil,
			)
		}
	} else if payload.EndDate != nil && existingExperience.StartDate != nil {
		if existingExperience.StartDate.After(*payload.EndDate) {
			return nil, errs.NewBadRequestError(
				"start date cannot be after end date",
				false, nil, nil, nil,
			)
		}
	}

	// Update experience in repository
	updatedExperience, err := s.experienceRepo.UpdateExperience(ctx, userID, experienceID, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to update experience: %w", err)
	}

	return s.convertToExperienceResponse(updatedExperience), nil
}

// BulkUpdateExperienceOrder updates the order of multiple experience entries
func (s *ExperienceService) BulkUpdateExperienceOrder(ctx context.Context, userID string, payload *experience.BulkUpdateExperienceRequest) error {
	// Validate that all experience entries belong to the user
	for _, expUpdate := range payload.Experience {
		experienceID, err := uuid.Parse(expUpdate.ID)
		if err != nil {
			return errs.NewBadRequestError("invalid experience ID", false, nil, nil, nil)
		}
		_, err = s.experienceRepo.GetExperienceByID(ctx, userID, experienceID)
		if err != nil {
			if err.Error() == "failed to collect row from table:experience" {
				return errs.NewNotFoundError("experience not found", false, nil)
			}
			return fmt.Errorf("failed to verify experience ownership: %w", err)
		}
	}

	// Update order in repository
	err := s.experienceRepo.BulkUpdateExperienceOrder(ctx, userID, payload)
	if err != nil {
		return fmt.Errorf("failed to update experience order: %w", err)
	}

	return nil
}

// DeleteExperience deletes an experience entry
func (s *ExperienceService) DeleteExperience(ctx context.Context, userID string, experienceID uuid.UUID) error {
	// Check if experience exists and belongs to user
	_, err := s.experienceRepo.GetExperienceByID(ctx, userID, experienceID)
	if err != nil {
		if err.Error() == "failed to collect row from table:experience" {
			return errs.NewNotFoundError("experience not found", false, nil)
		}
		return fmt.Errorf("failed to get existing experience: %w", err)
	}

	// Delete experience
	err = s.experienceRepo.DeleteExperience(ctx, userID, experienceID)
	if err != nil {
		return fmt.Errorf("failed to delete experience: %w", err)
	}

	return nil
}

// Helper methods

func (s *ExperienceService) convertToExperienceResponse(experienceItem *experience.Experience) *experience.ExperienceResponse {
	response := &experience.ExperienceResponse{
		ID:          experienceItem.ID.String(),
		ResumeID:    experienceItem.ResumeID,
		Company:     experienceItem.Company,
		Position:    experienceItem.Position,
		Location:    experienceItem.Location,
		Description: experienceItem.Description,
		OrderIndex:  experienceItem.OrderIndex,
		CreatedAt:   experienceItem.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   experienceItem.UpdatedAt.Format(time.RFC3339),
	}

	// Handle optional date fields
	if experienceItem.StartDate != nil {
		response.StartDate = experienceItem.StartDate
	}
	if experienceItem.EndDate != nil {
		response.EndDate = experienceItem.EndDate
	}

	return response
}
