package skill

import (
	"github.com/google/uuid"
	"github.com/sriniously/go-resumify/internal/model"
)

// Skill represents skill entries
type Skill struct {
	model.BaseWithId
	ResumeID   uuid.UUID `json:"resumeId" db:"resume_id"`
	Name       *string   `json:"name" db:"name"`
	Level      *string   `json:"level" db:"level"`
	Category   *string   `json:"category" db:"category"`
	OrderIndex int       `json:"orderIndex" db:"order_index"`
}
