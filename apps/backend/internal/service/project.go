package service

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/recreatedev/Resumify/internal/errs"
	"github.com/recreatedev/Resumify/internal/model/project"
	"github.com/recreatedev/Resumify/internal/repository"
	"github.com/recreatedev/Resumify/internal/server"
)

type ProjectService struct {
	server      *server.Server
	projectRepo *repository.ProjectRepository
	resumeRepo  *repository.ResumeRepository
}

func NewProjectService(s *server.Server, repos *repository.Repositories) *ProjectService {
	return &ProjectService{
		server:      s,
		projectRepo: repos.Project,
		resumeRepo:  repos.Resume,
	}
}

// CreateProject creates a new project entry
func (s *ProjectService) CreateProject(ctx context.Context, userID string, payload *project.CreateProjectRequest) (*project.ProjectResponse, error) {
	// Verify resume belongs to user
	_, err := s.resumeRepo.GetResumeByID(ctx, userID, payload.ResumeID)
	if err != nil {
		if err.Error() == "failed to collect row from table:resumes" {
			return nil, errs.NewNotFoundError("resume not found", false, nil)
		}
		return nil, fmt.Errorf("failed to verify resume ownership: %w", err)
	}

	// Business logic: Validate URL if provided
	if payload.Link != nil && *payload.Link != "" {
		if _, err := url.Parse(*payload.Link); err != nil {
			return nil, errs.NewBadRequestError(
				"invalid project URL",
				false, nil, nil, nil,
			)
		}
	}

	// Business logic: Check for duplicate project entries
	existingProjects, err := s.projectRepo.GetProjectsByResumeID(ctx, userID, payload.ResumeID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing projects: %w", err)
	}

	// Check for duplicates based on name
	for _, existing := range existingProjects {
		if existing.Name == payload.Name {
			return nil, errs.NewBadRequestError(
				"project with same name already exists",
				false, nil, nil, nil,
			)
		}
	}

	// Set default order index if not provided
	if payload.OrderIndex == 0 {
		payload.OrderIndex = len(existingProjects) + 1
	}

	// Create project in repository
	projectItem, err := s.projectRepo.CreateProject(ctx, userID, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	// Convert to response DTO
	response := s.convertToProjectResponse(projectItem)

	return response, nil
}

// GetProjectByID retrieves a project entry by ID
func (s *ProjectService) GetProjectByID(ctx context.Context, userID string, projectID uuid.UUID) (*project.ProjectResponse, error) {
	projectItem, err := s.projectRepo.GetProjectByID(ctx, userID, projectID)
	if err != nil {
		if err.Error() == "failed to collect row from table:projects" {
			return nil, errs.NewNotFoundError("project not found", false, nil)
		}
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return s.convertToProjectResponse(projectItem), nil
}

// GetProjectsByResumeID retrieves all project entries for a resume
func (s *ProjectService) GetProjectsByResumeID(ctx context.Context, userID string, resumeID uuid.UUID) ([]project.ProjectResponse, error) {
	// Verify resume belongs to user
	_, err := s.resumeRepo.GetResumeByID(ctx, userID, resumeID)
	if err != nil {
		if err.Error() == "failed to collect row from table:resumes" {
			return nil, errs.NewNotFoundError("resume not found", false, nil)
		}
		return nil, fmt.Errorf("failed to verify resume ownership: %w", err)
	}

	projectItems, err := s.projectRepo.GetProjectsByResumeID(ctx, userID, resumeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project entries: %w", err)
	}

	// Convert to response DTOs
	responses := make([]project.ProjectResponse, len(projectItems))
	for i, item := range projectItems {
		responses[i] = *s.convertToProjectResponse(&item)
	}

	return responses, nil
}

// UpdateProject updates a project entry
func (s *ProjectService) UpdateProject(ctx context.Context, userID string, projectID uuid.UUID, payload *project.UpdateProjectRequest) (*project.ProjectResponse, error) {
	// Check if project exists and belongs to user
	existingProject, err := s.projectRepo.GetProjectByID(ctx, userID, projectID)
	if err != nil {
		if err.Error() == "failed to collect row from table:projects" {
			return nil, errs.NewNotFoundError("project not found", false, nil)
		}
		return nil, fmt.Errorf("failed to get existing project: %w", err)
	}

	// Business logic: Validate URL if provided
	if payload.Link != nil && *payload.Link != "" {
		if _, err := url.Parse(*payload.Link); err != nil {
			return nil, errs.NewBadRequestError(
				"invalid project URL",
				false, nil, nil, nil,
			)
		}
	}

	// Business logic: Check for duplicate project names (excluding current project)
	if payload.Name != nil && *payload.Name != *existingProject.Name {
		projects, err := s.projectRepo.GetProjectsByResumeID(ctx, userID, existingProject.ResumeID)
		if err != nil {
			return nil, fmt.Errorf("failed to check existing projects: %w", err)
		}

		for _, proj := range projects {
			if proj.ID != projectID && *proj.Name == *payload.Name {
				return nil, errs.NewBadRequestError(
					"project with same name already exists",
					false, nil, nil, nil,
				)
			}
		}
	}

	// Update project in repository
	updatedProject, err := s.projectRepo.UpdateProject(ctx, userID, projectID, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	return s.convertToProjectResponse(updatedProject), nil
}

// BulkUpdateProjectOrder updates the order of multiple project entries
func (s *ProjectService) BulkUpdateProjectOrder(ctx context.Context, userID string, payload *project.BulkUpdateProjectsRequest) error {
	// Validate that all project entries belong to the user
	for _, projUpdate := range payload.Projects {
		projectID, err := uuid.Parse(projUpdate.ID)
		if err != nil {
			return errs.NewBadRequestError("invalid project ID", false, nil, nil, nil)
		}
		_, err = s.projectRepo.GetProjectByID(ctx, userID, projectID)
		if err != nil {
			if err.Error() == "failed to collect row from table:projects" {
				return errs.NewNotFoundError("project not found", false, nil)
			}
			return fmt.Errorf("failed to verify project ownership: %w", err)
		}
	}

	// Update order in repository
	err := s.projectRepo.BulkUpdateProjectOrder(ctx, userID, payload)
	if err != nil {
		return fmt.Errorf("failed to update project order: %w", err)
	}

	return nil
}

// DeleteProject deletes a project entry
func (s *ProjectService) DeleteProject(ctx context.Context, userID string, projectID uuid.UUID) error {
	// Check if project exists and belongs to user
	_, err := s.projectRepo.GetProjectByID(ctx, userID, projectID)
	if err != nil {
		if err.Error() == "failed to collect row from table:projects" {
			return errs.NewNotFoundError("project not found", false, nil)
		}
		return fmt.Errorf("failed to get existing project: %w", err)
	}

	// Delete project
	err = s.projectRepo.DeleteProject(ctx, userID, projectID)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	return nil
}

// Helper methods

func (s *ProjectService) convertToProjectResponse(projectItem *project.Project) *project.ProjectResponse {
	response := &project.ProjectResponse{
		ID:           projectItem.ID.String(),
		ResumeID:     projectItem.ResumeID,
		Name:         projectItem.Name,
		Role:         projectItem.Role,
		Description:  projectItem.Description,
		Link:         projectItem.Link,
		Technologies: projectItem.Technologies,
		OrderIndex:   projectItem.OrderIndex,
		CreatedAt:    projectItem.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    projectItem.UpdatedAt.Format(time.RFC3339),
	}

	return response
}
