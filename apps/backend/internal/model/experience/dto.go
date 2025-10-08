package experience

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// CreateExperienceRequest represents the request to create a new experience entry
type CreateExperienceRequest struct {
	ResumeID    uuid.UUID  `json:"resumeId" validate:"required"`
	Company     *string    `json:"company" validate:"omitempty,max=200"`
	Position    *string    `json:"position" validate:"omitempty,max=200"`
	StartDate   *time.Time `json:"startDate"`
	EndDate     *time.Time `json:"endDate"`
	Location    *string    `json:"location" validate:"omitempty,max=200"`
	Description *string    `json:"description" validate:"omitempty,max=2000"`
	OrderIndex  int        `json:"orderIndex" validate:"min=0"`
}

// UpdateExperienceRequest represents the request to update an existing experience entry
type UpdateExperienceRequest struct {
	Company     *string    `json:"company" validate:"omitempty,max=200"`
	Position    *string    `json:"position" validate:"omitempty,max=200"`
	StartDate   *time.Time `json:"startDate"`
	EndDate     *time.Time `json:"endDate"`
	Location    *string    `json:"location" validate:"omitempty,max=200"`
	Description *string    `json:"description" validate:"omitempty,max=2000"`
	OrderIndex  *int       `json:"orderIndex" validate:"omitempty,min=0"`
}

// ExperienceResponse represents the response for experience data
type ExperienceResponse struct {
	ID          string     `json:"id"`
	ResumeID    uuid.UUID  `json:"resumeId"`
	Company     *string    `json:"company"`
	Position    *string    `json:"position"`
	StartDate   *time.Time `json:"startDate"`
	EndDate     *time.Time `json:"endDate"`
	Location    *string    `json:"location"`
	Description *string    `json:"description"`
	OrderIndex  int        `json:"orderIndex"`
	CreatedAt   string     `json:"createdAt"`
	UpdatedAt   string     `json:"updatedAt"`
}

// BulkUpdateExperienceRequest represents the request to update multiple experience entries order
type BulkUpdateExperienceRequest struct {
	Experience []ExperienceOrderUpdate `json:"experience" validate:"required,min=1"`
}

// ExperienceOrderUpdate represents a single experience order update
type ExperienceOrderUpdate struct {
	ID         string `json:"id" validate:"required"`
	OrderIndex int    `json:"orderIndex" validate:"min=0"`
}

// Validate implements the Validatable interface for CreateExperienceRequest
func (r *CreateExperienceRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// Validate implements the Validatable interface for UpdateExperienceRequest
func (r *UpdateExperienceRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// Validate implements the Validatable interface for BulkUpdateExperienceRequest
func (r *BulkUpdateExperienceRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
