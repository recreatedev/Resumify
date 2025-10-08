package project

import (
	"github.com/google/uuid"
	"github.com/recreatedev/Resumify/internal/model"
)

// Project represents project entries
type Project struct {
	model.Base
	ResumeID     uuid.UUID `json:"resumeId" db:"resume_id"`
	Name         *string   `json:"name" db:"name"`
	Role         *string   `json:"role" db:"role"`
	Description  *string   `json:"description" db:"description"`
	Link         *string   `json:"link" db:"link"`
	Technologies []string  `json:"technologies" db:"technologies"`
	OrderIndex   int       `json:"orderIndex" db:"order_index"`
}
