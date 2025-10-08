package education

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// CreateEducationRequest represents the request to create a new education entry
type CreateEducationRequest struct {
	ResumeID     uuid.UUID  `json:"resumeId" validate:"required"`
	Institution  *string    `json:"institution" validate:"omitempty,max=200"`
	Degree       *string    `json:"degree" validate:"omitempty,max=100"`
	FieldOfStudy *string    `json:"fieldOfStudy" validate:"omitempty,max=100"`
	StartDate    *time.Time `json:"startDate"`
	EndDate      *time.Time `json:"endDate"`
	Grade        *string    `json:"grade" validate:"omitempty,max=50"`
	Description  *string    `json:"description" validate:"omitempty,max=1000"`
	OrderIndex   int        `json:"orderIndex" validate:"min=0"`
}

// UpdateEducationRequest represents the request to update an existing education entry
type UpdateEducationRequest struct {
	Institution  *string    `json:"institution" validate:"omitempty,max=200"`
	Degree       *string    `json:"degree" validate:"omitempty,max=100"`
	FieldOfStudy *string    `json:"fieldOfStudy" validate:"omitempty,max=100"`
	StartDate    *time.Time `json:"startDate"`
	EndDate      *time.Time `json:"endDate"`
	Grade        *string    `json:"grade" validate:"omitempty,max=50"`
	Description  *string    `json:"description" validate:"omitempty,max=1000"`
	OrderIndex   *int       `json:"orderIndex" validate:"omitempty,min=0"`
}

// EducationResponse represents the response for education data
type EducationResponse struct {
	ID           string     `json:"id"`
	ResumeID     uuid.UUID  `json:"resumeId"`
	Institution  *string    `json:"institution"`
	Degree       *string    `json:"degree"`
	FieldOfStudy *string    `json:"fieldOfStudy"`
	StartDate    *time.Time `json:"startDate"`
	EndDate      *time.Time `json:"endDate"`
	Grade        *string    `json:"grade"`
	Description  *string    `json:"description"`
	OrderIndex   int        `json:"orderIndex"`
	CreatedAt    string     `json:"createdAt"`
	UpdatedAt    string     `json:"updatedAt"`
}

// BulkUpdateEducationRequest represents the request to update multiple education entries order
type BulkUpdateEducationRequest struct {
	Education []EducationOrderUpdate `json:"education" validate:"required,min=1"`
}

// EducationOrderUpdate represents a single education order update
type EducationOrderUpdate struct {
	ID         string `json:"id" validate:"required"`
	OrderIndex int    `json:"orderIndex" validate:"min=0"`
}

// Validate implements the Validatable interface for CreateEducationRequest
func (r *CreateEducationRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// Validate implements the Validatable interface for UpdateEducationRequest
func (r *UpdateEducationRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// Validate implements the Validatable interface for BulkUpdateEducationRequest
func (r *BulkUpdateEducationRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
