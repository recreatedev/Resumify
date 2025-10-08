package certification

import (
	"time"

	"github.com/google/uuid"
	"github.com/sriniously/go-resumify/internal/model"
)

// Certification represents certification entries
type Certification struct {
	model.BaseWithId
	ResumeID      uuid.UUID  `json:"resumeId" db:"resume_id"`
	Name          *string    `json:"name" db:"name"`
	Organization  *string    `json:"organization" db:"organization"`
	IssueDate     *time.Time `json:"issueDate" db:"issue_date"`
	ExpiryDate    *time.Time `json:"expiryDate" db:"expiry_date"`
	CredentialID  *string    `json:"credentialId" db:"credential_id"`
	CredentialURL *string    `json:"credentialUrl" db:"credential_url"`
	OrderIndex    int        `json:"orderIndex" db:"order_index"`
}
