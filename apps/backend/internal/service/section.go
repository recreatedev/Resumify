package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/recreatedev/Resumify/internal/errs"
	"github.com/recreatedev/Resumify/internal/model/section"
	"github.com/recreatedev/Resumify/internal/repository"
	"github.com/recreatedev/Resumify/internal/server"
)

type SectionService struct {
	server      *server.Server
	sectionRepo *repository.ResumeSectionRepository
	resumeRepo  *repository.ResumeRepository
}

func NewSectionService(s *server.Server, repos *repository.Repositories) *SectionService {
	return &SectionService{
		server:      s,
		sectionRepo: repos.Section,
		resumeRepo:  repos.Resume,
	}
}

// CreateSection creates a new resume section
func (s *SectionService) CreateSection(ctx context.Context, userID string, payload *section.CreateSectionRequest) (*section.SectionResponse, error) {
	// Verify resume belongs to user
	_, err := s.resumeRepo.GetResumeByID(ctx, userID, payload.ResumeID)
	if err != nil {
		if err.Error() == "failed to collect row from table:resumes" {
			return nil, errs.NewNotFoundError("resume not found", false, nil)
		}
		return nil, fmt.Errorf("failed to verify resume ownership: %w", err)
	}

	// Business logic: Validate section name
	validSections := []string{"education", "experience", "projects", "skills", "certifications", "summary", "contact"}
	isValidSection := false
	for _, validSection := range validSections {
		if payload.Name == validSection {
			isValidSection = true
			break
		}
	}
	if !isValidSection {
		return nil, errs.NewBadRequestError(
			"invalid section name. Must be one of: education, experience, projects, skills, certifications, summary, contact",
			false, nil, nil, nil,
		)
	}

	// Business logic: Check for duplicate section names
	existingSections, err := s.sectionRepo.GetSectionsByResumeID(ctx, userID, payload.ResumeID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing sections: %w", err)
	}

	// Check for duplicates based on name
	for _, existing := range existingSections {
		if existing.Name == payload.Name {
			return nil, errs.NewBadRequestError(
				"section with same name already exists",
				false, nil, nil, nil,
			)
		}
	}

	// Set default order index if not provided
	if payload.OrderIndex == 0 {
		payload.OrderIndex = len(existingSections) + 1
	}

	// Set default display name if not provided
	if payload.DisplayName == nil || *payload.DisplayName == "" {
		displayName := s.getDefaultDisplayName(payload.Name)
		payload.DisplayName = &displayName
	}

	// Create section in repository
	sectionItem, err := s.sectionRepo.CreateSection(ctx, userID, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create section: %w", err)
	}

	// Convert to response DTO
	response := s.convertToSectionResponse(sectionItem)

	return response, nil
}

// GetSectionByID retrieves a section by ID
func (s *SectionService) GetSectionByID(ctx context.Context, userID string, sectionID uuid.UUID) (*section.SectionResponse, error) {
	sectionItem, err := s.sectionRepo.GetSectionByID(ctx, userID, sectionID)
	if err != nil {
		if err.Error() == "failed to collect row from table:resume_sections" {
			return nil, errs.NewNotFoundError("section not found", false, nil)
		}
		return nil, fmt.Errorf("failed to get section: %w", err)
	}

	return s.convertToSectionResponse(sectionItem), nil
}

// GetSectionsByResumeID retrieves all sections for a resume
func (s *SectionService) GetSectionsByResumeID(ctx context.Context, userID string, resumeID uuid.UUID) ([]section.SectionResponse, error) {
	// Verify resume belongs to user
	_, err := s.resumeRepo.GetResumeByID(ctx, userID, resumeID)
	if err != nil {
		if err.Error() == "failed to collect row from table:resumes" {
			return nil, errs.NewNotFoundError("resume not found", false, nil)
		}
		return nil, fmt.Errorf("failed to verify resume ownership: %w", err)
	}

	sectionItems, err := s.sectionRepo.GetSectionsByResumeID(ctx, userID, resumeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sections: %w", err)
	}

	// Convert to response DTOs
	responses := make([]section.SectionResponse, len(sectionItems))
	for i, item := range sectionItems {
		responses[i] = *s.convertToSectionResponse(&item)
	}

	return responses, nil
}

// UpdateSection updates a section
func (s *SectionService) UpdateSection(ctx context.Context, userID string, sectionID uuid.UUID, payload *section.UpdateSectionRequest) (*section.SectionResponse, error) {
	// Check if section exists and belongs to user
	existingSection, err := s.sectionRepo.GetSectionByID(ctx, userID, sectionID)
	if err != nil {
		if err.Error() == "failed to collect row from table:resume_sections" {
			return nil, errs.NewNotFoundError("section not found", false, nil)
		}
		return nil, fmt.Errorf("failed to get existing section: %w", err)
	}

	// Business logic: Validate section name if provided
	if payload.Name != nil {
		validSections := []string{"education", "experience", "projects", "skills", "certifications", "summary", "contact"}
		isValidSection := false
		for _, validSection := range validSections {
			if *payload.Name == validSection {
				isValidSection = true
				break
			}
		}
		if !isValidSection {
			return nil, errs.NewBadRequestError(
				"invalid section name. Must be one of: education, experience, projects, skills, certifications, summary, contact",
				false, nil, nil, nil,
			)
		}

		// Check for duplicate section names (excluding current section)
		if *payload.Name != existingSection.Name {
			sections, err := s.sectionRepo.GetSectionsByResumeID(ctx, userID, existingSection.ResumeID)
			if err != nil {
				return nil, fmt.Errorf("failed to check existing sections: %w", err)
			}

			for _, sec := range sections {
				if sec.ID != sectionID && sec.Name == *payload.Name {
					return nil, errs.NewBadRequestError(
						"section with same name already exists",
						false, nil, nil, nil,
					)
				}
			}
		}
	}

	// Update section in repository
	updatedSection, err := s.sectionRepo.UpdateSection(ctx, userID, sectionID, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to update section: %w", err)
	}

	return s.convertToSectionResponse(updatedSection), nil
}

// BulkUpdateSectionOrder updates the order of multiple sections
func (s *SectionService) BulkUpdateSectionOrder(ctx context.Context, userID string, payload *section.BulkUpdateSectionsRequest) error {
	// Validate that all sections belong to the user
	for _, sectionUpdate := range payload.Sections {
		sectionID, err := uuid.Parse(sectionUpdate.ID)
		if err != nil {
			return errs.NewBadRequestError("invalid section ID", false, nil, nil, nil)
		}
		_, err = s.sectionRepo.GetSectionByID(ctx, userID, sectionID)
		if err != nil {
			if err.Error() == "failed to collect row from table:resume_sections" {
				return errs.NewNotFoundError("section not found", false, nil)
			}
			return fmt.Errorf("failed to verify section ownership: %w", err)
		}
	}

	// Update order in repository
	err := s.sectionRepo.BulkUpdateSectionOrder(ctx, userID, payload)
	if err != nil {
		return fmt.Errorf("failed to update section order: %w", err)
	}

	return nil
}

// DeleteSection deletes a section
func (s *SectionService) DeleteSection(ctx context.Context, userID string, sectionID uuid.UUID) error {
	// Check if section exists and belongs to user
	_, err := s.sectionRepo.GetSectionByID(ctx, userID, sectionID)
	if err != nil {
		if err.Error() == "failed to collect row from table:resume_sections" {
			return errs.NewNotFoundError("section not found", false, nil)
		}
		return fmt.Errorf("failed to get existing section: %w", err)
	}

	// Delete section
	err = s.sectionRepo.DeleteSection(ctx, userID, sectionID)
	if err != nil {
		return fmt.Errorf("failed to delete section: %w", err)
	}

	return nil
}

// Helper methods

func (s *SectionService) convertToSectionResponse(sectionItem *section.ResumeSection) *section.SectionResponse {
	response := &section.SectionResponse{
		ID:          sectionItem.ID.String(),
		ResumeID:    sectionItem.ResumeID,
		Name:        sectionItem.Name,
		DisplayName: sectionItem.DisplayName,
		IsVisible:   sectionItem.IsVisible,
		OrderIndex:  sectionItem.OrderIndex,
		CreatedAt:   sectionItem.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   sectionItem.UpdatedAt.Format(time.RFC3339),
	}

	return response
}

func (s *SectionService) getDefaultDisplayName(sectionName string) string {
	displayNames := map[string]string{
		"education":      "Education",
		"experience":     "Work Experience",
		"projects":       "Projects",
		"skills":         "Skills",
		"certifications": "Certifications",
		"summary":        "Summary",
		"contact":        "Contact Information",
	}

	if displayName, exists := displayNames[sectionName]; exists {
		return displayName
	}
	return sectionName
}
