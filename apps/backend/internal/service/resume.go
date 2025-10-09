package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/recreatedev/Resumify/internal/errs"
	"github.com/recreatedev/Resumify/internal/lib/email"
	"github.com/recreatedev/Resumify/internal/model"
	"github.com/recreatedev/Resumify/internal/model/resume"
	"github.com/recreatedev/Resumify/internal/repository"
	"github.com/recreatedev/Resumify/internal/server"
)

type ResumeService struct {
	server         *server.Server
	resumeRepo     *repository.ResumeRepository
	sectionRepo    *repository.ResumeSectionRepository
	educationRepo  *repository.EducationRepository
	experienceRepo *repository.ExperienceRepository
	projectRepo    *repository.ProjectRepository
	skillRepo      *repository.SkillRepository
	certRepo       *repository.CertificationRepository
	emailClient    *email.Client
}

func NewResumeService(s *server.Server, repos *repository.Repositories) *ResumeService {
	return &ResumeService{
		server:         s,
		resumeRepo:     repos.Resume,
		sectionRepo:    repos.Section,
		educationRepo:  repos.Education,
		experienceRepo: repos.Experience,
		projectRepo:    repos.Project,
		skillRepo:      repos.Skill,
		certRepo:       repos.Certification,
		emailClient:    nil, // TODO: Initialize email client when available
	}
}

// CreateResume creates a new resume with business logic validation
func (s *ResumeService) CreateResume(ctx context.Context, userID string, payload *resume.CreateResumeRequest) (*resume.ResumeResponse, error) {
	// Business logic: Check if user has reached maximum resume limit
	maxResumes := 10 // Configurable business rule
	existingResumes, err := s.resumeRepo.GetResumes(ctx, userID, 1, maxResumes+1)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing resumes: %w", err)
	}

	if existingResumes.Total >= maxResumes {
		return nil, errs.NewBadRequestError(
			fmt.Sprintf("maximum number of resumes (%d) reached", maxResumes),
			false, nil, nil, nil,
		)
	}

	// Set default theme if not provided
	if payload.Theme == "" {
		payload.Theme = "default"
	}

	// Create resume in repository
	resumeItem, err := s.resumeRepo.CreateResume(ctx, userID, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create resume: %w", err)
	}

	// Convert to response DTO
	response := s.convertToResumeResponse(resumeItem)

	// TODO: Create default sections for new resume
	// TODO: Send welcome email for first resume
	// TODO: Log resume creation event

	return response, nil
}

// GetResumeByID retrieves a resume by ID with authorization check
func (s *ResumeService) GetResumeByID(ctx context.Context, userID string, resumeID uuid.UUID) (*resume.ResumeResponse, error) {
	// Repository already handles authorization by checking user_id
	resumeItem, err := s.resumeRepo.GetResumeByID(ctx, userID, resumeID)
	if err != nil {
		if err.Error() == "failed to collect row from table:resumes" {
			return nil, errs.NewNotFoundError("resume not found", false, nil)
		}
		return nil, fmt.Errorf("failed to get resume: %w", err)
	}

	return s.convertToResumeResponse(resumeItem), nil
}

// GetResumes retrieves paginated list of user's resumes
func (s *ResumeService) GetResumes(ctx context.Context, userID string, page, limit int) (*model.PaginatedResponse[resume.ResumeSummaryResponse], error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20 // Default limit
	}

	// Get resumes from repository
	resumes, err := s.resumeRepo.GetResumes(ctx, userID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get resumes: %w", err)
	}

	// Convert to summary responses
	summaryResponses := make([]resume.ResumeSummaryResponse, len(resumes.Data))
	for i, resumeItem := range resumes.Data {
		summaryResponses[i] = s.convertToResumeSummaryResponse(&resumeItem)
	}

	return &model.PaginatedResponse[resume.ResumeSummaryResponse]{
		Data:       summaryResponses,
		Page:       resumes.Page,
		Limit:      resumes.Limit,
		Total:      resumes.Total,
		TotalPages: resumes.TotalPages,
	}, nil
}

// UpdateResume updates a resume with business logic validation
func (s *ResumeService) UpdateResume(ctx context.Context, userID string, resumeID uuid.UUID, payload *resume.UpdateResumeRequest) (*resume.ResumeResponse, error) {
	// Check if resume exists and belongs to user
	existingResume, err := s.resumeRepo.GetResumeByID(ctx, userID, resumeID)
	if err != nil {
		if err.Error() == "failed to collect row from table:resumes" {
			return nil, errs.NewNotFoundError("resume not found", false, nil)
		}
		return nil, fmt.Errorf("failed to get existing resume: %w", err)
	}

	// Business logic: Validate theme transition
	if payload.Theme != nil && *payload.Theme != existingResume.Theme {
		if !s.isValidThemeTransition(existingResume.Theme, *payload.Theme) {
			return nil, errs.NewBadRequestError(
				"invalid theme transition",
				false, nil, nil, nil,
			)
		}
	}

	// Update resume in repository
	updatedResume, err := s.resumeRepo.UpdateResume(ctx, userID, resumeID, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to update resume: %w", err)
	}

	// Convert to response DTO
	response := s.convertToResumeResponse(updatedResume)

	// TODO: Log resume update event
	// TODO: Send notification if significant changes

	return response, nil
}

// DeleteResume deletes a resume and all related data
func (s *ResumeService) DeleteResume(ctx context.Context, userID string, resumeID uuid.UUID) error {
	// Check if resume exists and belongs to user
	_, err := s.resumeRepo.GetResumeByID(ctx, userID, resumeID)
	if err != nil {
		if err.Error() == "failed to collect row from table:resumes" {
			return errs.NewNotFoundError("resume not found", false, nil)
		}
		return fmt.Errorf("failed to get existing resume: %w", err)
	}

	// TODO: Implement cascade deletion in a transaction
	// For now, rely on database CASCADE DELETE constraints

	// Delete resume (cascade will handle related data)
	err = s.resumeRepo.DeleteResume(ctx, userID, resumeID)
	if err != nil {
		return fmt.Errorf("failed to delete resume: %w", err)
	}

	// TODO: Log resume deletion event
	// TODO: Send deletion confirmation email

	return nil
}

// DuplicateResume creates a copy of an existing resume
func (s *ResumeService) DuplicateResume(ctx context.Context, userID string, resumeID uuid.UUID) (*resume.ResumeResponse, error) {
	// Get original resume
	originalResume, err := s.resumeRepo.GetResumeByID(ctx, userID, resumeID)
	if err != nil {
		if err.Error() == "failed to collect row from table:resumes" {
			return nil, errs.NewNotFoundError("resume not found", false, nil)
		}
		return nil, fmt.Errorf("failed to get original resume: %w", err)
	}

	// Create duplicate with modified title
	duplicateTitle := fmt.Sprintf("%s (Copy)", originalResume.Title)
	createPayload := &resume.CreateResumeRequest{
		Title: duplicateTitle,
		Theme: originalResume.Theme,
	}

	duplicateResume, err := s.resumeRepo.CreateResume(ctx, userID, createPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to create duplicate resume: %w", err)
	}

	// TODO: Copy all sections and related data
	// TODO: Log duplication event

	return s.convertToResumeResponse(duplicateResume), nil
}

// GetResumeWithSections retrieves a resume with all its sections and data
func (s *ResumeService) GetResumeWithSections(ctx context.Context, userID string, resumeID uuid.UUID) (*ResumeWithSectionsResponse, error) {
	// Get resume
	resumeItem, err := s.resumeRepo.GetResumeByID(ctx, userID, resumeID)
	if err != nil {
		if err.Error() == "failed to collect row from table:resumes" {
			return nil, errs.NewNotFoundError("resume not found", false, nil)
		}
		return nil, fmt.Errorf("failed to get resume: %w", err)
	}

	// TODO: Get all sections and their data
	// For now, return basic resume with placeholder for sections

	response := &ResumeWithSectionsResponse{
		ResumeResponse: *s.convertToResumeResponse(resumeItem),
		Sections:       []SectionData{}, // TODO: Implement section retrieval
	}

	return response, nil
}

// Helper methods

func (s *ResumeService) convertToResumeResponse(resumeItem *resume.Resume) *resume.ResumeResponse {
	return &resume.ResumeResponse{
		ID:        resumeItem.ID.String(),
		UserID:    resumeItem.UserID,
		Title:     resumeItem.Title,
		Theme:     resumeItem.Theme,
		CreatedAt: resumeItem.CreatedAt.Format(time.RFC3339),
		UpdatedAt: resumeItem.UpdatedAt.Format(time.RFC3339),
	}
}

func (s *ResumeService) convertToResumeSummaryResponse(resumeItem *resume.Resume) resume.ResumeSummaryResponse {
	return resume.ResumeSummaryResponse{
		ID:        resumeItem.ID.String(),
		Title:     resumeItem.Title,
		Theme:     resumeItem.Theme,
		CreatedAt: resumeItem.CreatedAt.Format(time.RFC3339),
		UpdatedAt: resumeItem.UpdatedAt.Format(time.RFC3339),
	}
}

func (s *ResumeService) isValidThemeTransition(fromTheme, toTheme string) bool {
	// Business rule: Define valid theme transitions
	validTransitions := map[string][]string{
		"default":      {"modern", "classic", "professional"},
		"modern":       {"default", "professional"},
		"classic":      {"default", "professional"},
		"professional": {"default", "modern", "classic"},
	}

	allowedThemes, exists := validTransitions[fromTheme]
	if !exists {
		return false
	}

	for _, theme := range allowedThemes {
		if theme == toTheme {
			return true
		}
	}
	return false
}

// Response DTOs for complex operations

type ResumeWithSectionsResponse struct {
	resume.ResumeResponse
	Sections []SectionData `json:"sections"`
}

type SectionData struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	DisplayName string      `json:"displayName"`
	IsVisible   bool        `json:"isVisible"`
	OrderIndex  int         `json:"orderIndex"`
	Data        interface{} `json:"data"` // Will contain section-specific data
}
