package project

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// CreateProjectRequest represents the request to create a new project entry
type CreateProjectRequest struct {
	ResumeID     uuid.UUID `json:"resumeId" validate:"required"`
	Name         *string   `json:"name" validate:"omitempty,max=200"`
	Role         *string   `json:"role" validate:"omitempty,max=200"`
	Description  *string   `json:"description" validate:"omitempty,max=2000"`
	Link         *string   `json:"link" validate:"omitempty,url"`
	Technologies []string  `json:"technologies" validate:"omitempty,max=20"`
	OrderIndex   int       `json:"orderIndex" validate:"min=0"`
}

// UpdateProjectRequest represents the request to update an existing project entry
type UpdateProjectRequest struct {
	Name         *string  `json:"name" validate:"omitempty,max=200"`
	Role         *string  `json:"role" validate:"omitempty,max=200"`
	Description  *string  `json:"description" validate:"omitempty,max=2000"`
	Link         *string  `json:"link" validate:"omitempty,url"`
	Technologies []string `json:"technologies" validate:"omitempty,max=20"`
	OrderIndex   *int     `json:"orderIndex" validate:"omitempty,min=0"`
}

// ProjectResponse represents the response for project data
type ProjectResponse struct {
	ID           string    `json:"id"`
	ResumeID     uuid.UUID `json:"resumeId"`
	Name         *string   `json:"name"`
	Role         *string   `json:"role"`
	Description  *string   `json:"description"`
	Link         *string   `json:"link"`
	Technologies []string  `json:"technologies"`
	OrderIndex   int       `json:"orderIndex"`
	CreatedAt    string    `json:"createdAt"`
	UpdatedAt    string    `json:"updatedAt"`
}

// BulkUpdateProjectsRequest represents the request to update multiple project entries order
type BulkUpdateProjectsRequest struct {
	Projects []ProjectOrderUpdate `json:"projects" validate:"required,min=1"`
}

// ProjectOrderUpdate represents a single project order update
type ProjectOrderUpdate struct {
	ID         string `json:"id" validate:"required"`
	OrderIndex int    `json:"orderIndex" validate:"min=0"`
}

// Validate implements the Validatable interface for CreateProjectRequest
func (r *CreateProjectRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// Validate implements the Validatable interface for UpdateProjectRequest
func (r *UpdateProjectRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// Validate implements the Validatable interface for BulkUpdateProjectsRequest
func (r *BulkUpdateProjectsRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
