package section

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// CreateSectionRequest represents the request to create a new resume section
type CreateSectionRequest struct {
	ResumeID    uuid.UUID `json:"resumeId" validate:"required"`
	Name        string    `json:"name" validate:"required,min=1,max=50"`
	DisplayName *string   `json:"displayName" validate:"omitempty,max=100"`
	IsVisible   bool      `json:"isVisible"`
	OrderIndex  int       `json:"orderIndex" validate:"min=0"`
}

// UpdateSectionRequest represents the request to update an existing resume section
type UpdateSectionRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=50"`
	DisplayName *string `json:"displayName" validate:"omitempty,max=100"`
	IsVisible   *bool   `json:"isVisible"`
	OrderIndex  *int    `json:"orderIndex" validate:"omitempty,min=0"`
}

// SectionResponse represents the response for resume section data
type SectionResponse struct {
	ID          string    `json:"id"`
	ResumeID    uuid.UUID `json:"resumeId"`
	Name        string    `json:"name"`
	DisplayName *string   `json:"displayName"`
	IsVisible   bool      `json:"isVisible"`
	OrderIndex  int       `json:"orderIndex"`
	CreatedAt   string    `json:"createdAt"`
	UpdatedAt   string    `json:"updatedAt"`
}

// BulkUpdateSectionsRequest represents the request to update multiple sections order
type BulkUpdateSectionsRequest struct {
	Sections []SectionOrderUpdate `json:"sections" validate:"required,min=1"`
}

// SectionOrderUpdate represents a single section order update
type SectionOrderUpdate struct {
	ID         string `json:"id" validate:"required"`
	OrderIndex int    `json:"orderIndex" validate:"min=0"`
}

// Validate implements the Validatable interface for CreateSectionRequest
func (r *CreateSectionRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// Validate implements the Validatable interface for UpdateSectionRequest
func (r *UpdateSectionRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// Validate implements the Validatable interface for BulkUpdateSectionsRequest
func (r *BulkUpdateSectionsRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
