package certification

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// CreateCertificationRequest represents the request to create a new certification entry
type CreateCertificationRequest struct {
	ResumeID      uuid.UUID  `json:"resumeId" validate:"required"`
	Name          *string    `json:"name" validate:"omitempty,max=200"`
	Organization  *string    `json:"organization" validate:"omitempty,max=200"`
	IssueDate     *time.Time `json:"issueDate"`
	ExpiryDate    *time.Time `json:"expiryDate"`
	CredentialID  *string    `json:"credentialId" validate:"omitempty,max=100"`
	CredentialURL *string    `json:"credentialUrl" validate:"omitempty,url"`
	OrderIndex    int        `json:"orderIndex" validate:"min=0"`
}

// UpdateCertificationRequest represents the request to update an existing certification entry
type UpdateCertificationRequest struct {
	Name          *string    `json:"name" validate:"omitempty,max=200"`
	Organization  *string    `json:"organization" validate:"omitempty,max=200"`
	IssueDate     *time.Time `json:"issueDate"`
	ExpiryDate    *time.Time `json:"expiryDate"`
	CredentialID  *string    `json:"credentialId" validate:"omitempty,max=100"`
	CredentialURL *string    `json:"credentialUrl" validate:"omitempty,url"`
	OrderIndex    *int       `json:"orderIndex" validate:"omitempty,min=0"`
}

// CertificationResponse represents the response for certification data
type CertificationResponse struct {
	ID            string     `json:"id"`
	ResumeID      uuid.UUID  `json:"resumeId"`
	Name          *string    `json:"name"`
	Organization  *string    `json:"organization"`
	IssueDate     *time.Time `json:"issueDate"`
	ExpiryDate    *time.Time `json:"expiryDate"`
	CredentialID  *string    `json:"credentialId"`
	CredentialURL *string    `json:"credentialUrl"`
	OrderIndex    int        `json:"orderIndex"`
}

// BulkUpdateCertificationsRequest represents the request to update multiple certification entries order
type BulkUpdateCertificationsRequest struct {
	Certifications []CertificationOrderUpdate `json:"certifications" validate:"required,min=1"`
}

// CertificationOrderUpdate represents a single certification order update
type CertificationOrderUpdate struct {
	ID         string `json:"id" validate:"required"`
	OrderIndex int    `json:"orderIndex" validate:"min=0"`
}

// Validate implements the Validatable interface for CreateCertificationRequest
func (r *CreateCertificationRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// Validate implements the Validatable interface for UpdateCertificationRequest
func (r *UpdateCertificationRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// Validate implements the Validatable interface for BulkUpdateCertificationsRequest
func (r *BulkUpdateCertificationsRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
