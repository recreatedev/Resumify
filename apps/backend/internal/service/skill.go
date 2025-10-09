package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/recreatedev/Resumify/internal/errs"
	"github.com/recreatedev/Resumify/internal/model/skill"
	"github.com/recreatedev/Resumify/internal/repository"
	"github.com/recreatedev/Resumify/internal/server"
)

type SkillService struct {
	server     *server.Server
	skillRepo  *repository.SkillRepository
	resumeRepo *repository.ResumeRepository
}

func NewSkillService(s *server.Server, repos *repository.Repositories) *SkillService {
	return &SkillService{
		server:     s,
		skillRepo:  repos.Skill,
		resumeRepo: repos.Resume,
	}
}

// CreateSkill creates a new skill entry
func (s *SkillService) CreateSkill(ctx context.Context, userID string, payload *skill.CreateSkillRequest) (*skill.SkillResponse, error) {
	// Verify resume belongs to user
	_, err := s.resumeRepo.GetResumeByID(ctx, userID, payload.ResumeID)
	if err != nil {
		if err.Error() == "failed to collect row from table:resumes" {
			return nil, errs.NewNotFoundError("resume not found", false, nil)
		}
		return nil, fmt.Errorf("failed to verify resume ownership: %w", err)
	}

	// Business logic: Validate skill level
	if payload.Level != nil {
		validLevels := []string{"Beginner", "Intermediate", "Advanced", "Expert"}
		isValidLevel := false
		for _, level := range validLevels {
			if *payload.Level == level {
				isValidLevel = true
				break
			}
		}
		if !isValidLevel {
			return nil, errs.NewBadRequestError(
				"invalid skill level. Must be one of: Beginner, Intermediate, Advanced, Expert",
				false, nil, nil, nil,
			)
		}
	}

	// Business logic: Check for duplicate skill entries
	existingSkills, err := s.skillRepo.GetSkillsByResumeID(ctx, userID, payload.ResumeID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing skills: %w", err)
	}

	// Check for duplicates based on name
	for _, existing := range existingSkills {
		if *existing.Name == *payload.Name {
			return nil, errs.NewBadRequestError(
				"skill with same name already exists",
				false, nil, nil, nil,
			)
		}
	}

	// Set default order index if not provided
	if payload.OrderIndex == 0 {
		payload.OrderIndex = len(existingSkills) + 1
	}

	// Set default category if not provided
	if payload.Category == nil || *payload.Category == "" {
		defaultCategory := "Other"
		payload.Category = &defaultCategory
	}

	// Create skill in repository
	skillItem, err := s.skillRepo.CreateSkill(ctx, userID, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create skill: %w", err)
	}

	// Convert to response DTO
	response := s.convertToSkillResponse(skillItem)

	return response, nil
}

// GetSkillByID retrieves a skill entry by ID
func (s *SkillService) GetSkillByID(ctx context.Context, userID string, skillID uuid.UUID) (*skill.SkillResponse, error) {
	skillItem, err := s.skillRepo.GetSkillByID(ctx, userID, skillID)
	if err != nil {
		if err.Error() == "failed to collect row from table:skills" {
			return nil, errs.NewNotFoundError("skill not found", false, nil)
		}
		return nil, fmt.Errorf("failed to get skill: %w", err)
	}

	return s.convertToSkillResponse(skillItem), nil
}

// GetSkillsByResumeID retrieves all skill entries for a resume
func (s *SkillService) GetSkillsByResumeID(ctx context.Context, userID string, resumeID uuid.UUID) ([]skill.SkillResponse, error) {
	// Verify resume belongs to user
	_, err := s.resumeRepo.GetResumeByID(ctx, userID, resumeID)
	if err != nil {
		if err.Error() == "failed to collect row from table:resumes" {
			return nil, errs.NewNotFoundError("resume not found", false, nil)
		}
		return nil, fmt.Errorf("failed to verify resume ownership: %w", err)
	}

	skillItems, err := s.skillRepo.GetSkillsByResumeID(ctx, userID, resumeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get skill entries: %w", err)
	}

	// Convert to response DTOs
	responses := make([]skill.SkillResponse, len(skillItems))
	for i, item := range skillItems {
		responses[i] = *s.convertToSkillResponse(&item)
	}

	return responses, nil
}

// GetSkillsByCategory retrieves skills grouped by category for a resume
func (s *SkillService) GetSkillsByCategory(ctx context.Context, userID string, resumeID uuid.UUID) ([]skill.SkillsByCategoryResponse, error) {
	// Verify resume belongs to user
	_, err := s.resumeRepo.GetResumeByID(ctx, userID, resumeID)
	if err != nil {
		if err.Error() == "failed to collect row from table:resumes" {
			return nil, errs.NewNotFoundError("resume not found", false, nil)
		}
		return nil, fmt.Errorf("failed to verify resume ownership: %w", err)
	}

	skillsByCategory, err := s.skillRepo.GetSkillsByCategory(ctx, userID, resumeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get skills by category: %w", err)
	}

	return skillsByCategory, nil
}

// UpdateSkill updates a skill entry
func (s *SkillService) UpdateSkill(ctx context.Context, userID string, skillID uuid.UUID, payload *skill.UpdateSkillRequest) (*skill.SkillResponse, error) {
	// Check if skill exists and belongs to user
	existingSkill, err := s.skillRepo.GetSkillByID(ctx, userID, skillID)
	if err != nil {
		if err.Error() == "failed to collect row from table:skills" {
			return nil, errs.NewNotFoundError("skill not found", false, nil)
		}
		return nil, fmt.Errorf("failed to get existing skill: %w", err)
	}

	// Business logic: Validate skill level if provided
	if payload.Level != nil {
		validLevels := []string{"beginner", "intermediate", "advanced", "expert"}
		isValidLevel := false
		for _, level := range validLevels {
			if *payload.Level == level {
				isValidLevel = true
				break
			}
		}
		if !isValidLevel {
			return nil, errs.NewBadRequestError(
				"invalid skill level. Must be one of: beginner, intermediate, advanced, expert",
				false, nil, nil, nil,
			)
		}
	}

	// Business logic: Check for duplicate skill names (excluding current skill)
	if payload.Name != nil && *payload.Name != *existingSkill.Name {
		skills, err := s.skillRepo.GetSkillsByResumeID(ctx, userID, existingSkill.ResumeID)
		if err != nil {
			return nil, fmt.Errorf("failed to check existing skills: %w", err)
		}

		for _, sk := range skills {
			if sk.ID != skillID && *sk.Name == *payload.Name {
				return nil, errs.NewBadRequestError(
					"skill with same name already exists",
					false, nil, nil, nil,
				)
			}
		}
	}

	// Update skill in repository
	updatedSkill, err := s.skillRepo.UpdateSkill(ctx, userID, skillID, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to update skill: %w", err)
	}

	return s.convertToSkillResponse(updatedSkill), nil
}

// BulkUpdateSkillOrder updates the order of multiple skill entries
func (s *SkillService) BulkUpdateSkillOrder(ctx context.Context, userID string, payload *skill.BulkUpdateSkillsRequest) error {
	// Validate that all skill entries belong to the user
	for _, skillUpdate := range payload.Skills {
		skillID, err := uuid.Parse(skillUpdate.ID)
		if err != nil {
			return errs.NewBadRequestError("invalid skill ID", false, nil, nil, nil)
		}
		_, err = s.skillRepo.GetSkillByID(ctx, userID, skillID)
		if err != nil {
			if err.Error() == "failed to collect row from table:skills" {
				return errs.NewNotFoundError("skill not found", false, nil)
			}
			return fmt.Errorf("failed to verify skill ownership: %w", err)
		}
	}

	// Update order in repository
	err := s.skillRepo.BulkUpdateSkillOrder(ctx, userID, payload)
	if err != nil {
		return fmt.Errorf("failed to update skill order: %w", err)
	}

	return nil
}

// DeleteSkill deletes a skill entry
func (s *SkillService) DeleteSkill(ctx context.Context, userID string, skillID uuid.UUID) error {
	// Check if skill exists and belongs to user
	_, err := s.skillRepo.GetSkillByID(ctx, userID, skillID)
	if err != nil {
		if err.Error() == "failed to collect row from table:skills" {
			return errs.NewNotFoundError("skill not found", false, nil)
		}
		return fmt.Errorf("failed to get existing skill: %w", err)
	}

	// Delete skill
	err = s.skillRepo.DeleteSkill(ctx, userID, skillID)
	if err != nil {
		return fmt.Errorf("failed to delete skill: %w", err)
	}

	return nil
}

// Helper methods

func (s *SkillService) convertToSkillResponse(skillItem *skill.Skill) *skill.SkillResponse {
	response := &skill.SkillResponse{
		ID:         skillItem.ID.String(),
		ResumeID:   skillItem.ResumeID,
		Name:       skillItem.Name,
		Level:      skillItem.Level,
		Category:   skillItem.Category,
		OrderIndex: skillItem.OrderIndex,
	}

	return response
}
