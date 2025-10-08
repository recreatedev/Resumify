package education

import (
	"time"

	"github.com/google/uuid"
	"github.com/recreatedev/Resumify/internal/model"
)

// Education represents education entries
type Education struct {
	model.Base
	ResumeID     uuid.UUID  `json:"resumeId" db:"resume_id"`
	Institution  *string    `json:"institution" db:"institution"`
	Degree       *string    `json:"degree" db:"degree"`
	FieldOfStudy *string    `json:"fieldOfStudy" db:"field_of_study"`
	StartDate    *time.Time `json:"startDate" db:"start_date"`
	EndDate      *time.Time `json:"endDate" db:"end_date"`
	Grade        *string    `json:"grade" db:"grade"`
	Description  *string    `json:"description" db:"description"`
	OrderIndex   int        `json:"orderIndex" db:"order_index"`
}
