package skill

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// CreateSkillRequest represents the request to create a new skill entry
type CreateSkillRequest struct {
	ResumeID   uuid.UUID `json:"resumeId" validate:"required"`
	Name       *string   `json:"name" validate:"omitempty,max=100"`
	Level      *string   `json:"level" validate:"omitempty,oneof=Beginner Intermediate Advanced Expert"`
	Category   *string   `json:"category" validate:"omitempty,max=50"`
	OrderIndex int       `json:"orderIndex" validate:"min=0"`
}

// UpdateSkillRequest represents the request to update an existing skill entry
type UpdateSkillRequest struct {
	Name       *string `json:"name" validate:"omitempty,max=100"`
	Level      *string `json:"level" validate:"omitempty,oneof=Beginner Intermediate Advanced Expert"`
	Category   *string `json:"category" validate:"omitempty,max=50"`
	OrderIndex *int    `json:"orderIndex" validate:"omitempty,min=0"`
}

// SkillResponse represents the response for skill data
type SkillResponse struct {
	ID         string    `json:"id"`
	ResumeID   uuid.UUID `json:"resumeId"`
	Name       *string   `json:"name"`
	Level      *string   `json:"level"`
	Category   *string   `json:"category"`
	OrderIndex int       `json:"orderIndex"`
}

// BulkUpdateSkillsRequest represents the request to update multiple skill entries order
type BulkUpdateSkillsRequest struct {
	Skills []SkillOrderUpdate `json:"skills" validate:"required,min=1"`
}

// SkillOrderUpdate represents a single skill order update
type SkillOrderUpdate struct {
	ID         string `json:"id" validate:"required"`
	OrderIndex int    `json:"orderIndex" validate:"min=0"`
}

// SkillsByCategoryResponse represents skills grouped by category
type SkillsByCategoryResponse struct {
	Category string          `json:"category"`
	Skills   []SkillResponse `json:"skills"`
}

// Validate implements the Validatable interface for CreateSkillRequest
func (r *CreateSkillRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// Validate implements the Validatable interface for UpdateSkillRequest
func (r *UpdateSkillRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// Validate implements the Validatable interface for BulkUpdateSkillsRequest
func (r *BulkUpdateSkillsRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
