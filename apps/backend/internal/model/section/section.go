package section

import (
	"github.com/google/uuid"
	"github.com/recreatedev/Resumify/internal/model"
)

// ResumeSection represents a section within a resume
type ResumeSection struct {
	model.Base
	ResumeID    uuid.UUID `json:"resumeId" db:"resume_id"`
	Name        string    `json:"name" db:"name"`
	DisplayName *string   `json:"displayName" db:"display_name"`
	IsVisible   bool      `json:"isVisible" db:"is_visible"`
	OrderIndex  int       `json:"orderIndex" db:"order_index"`
}
