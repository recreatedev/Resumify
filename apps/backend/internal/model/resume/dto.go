package resume

import "github.com/go-playground/validator/v10"

// CreateResumeRequest represents the request to create a new resume
type CreateResumeRequest struct {
	Title string `json:"title" validate:"required,min=1,max=100"`
	Theme string `json:"theme" validate:"omitempty,oneof=default modern classic professional"`
}

// UpdateResumeRequest represents the request to update an existing resume
type UpdateResumeRequest struct {
	Title *string `json:"title" validate:"omitempty,min=1,max=100"`
	Theme *string `json:"theme" validate:"omitempty,oneof=default modern classic professional"`
}

// ResumeResponse represents the response for resume data
type ResumeResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"userId"`
	Title     string `json:"title"`
	Theme     string `json:"theme"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// ResumeSummaryResponse represents a summary of resume data (for lists)
type ResumeSummaryResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Theme     string `json:"theme"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// Validate implements the Validatable interface for CreateResumeRequest
func (r *CreateResumeRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// Validate implements the Validatable interface for UpdateResumeRequest
func (r *UpdateResumeRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
