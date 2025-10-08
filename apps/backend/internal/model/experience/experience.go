package experience

import (
	"time"

	"github.com/google/uuid"
	"github.com/sriniously/go-resumify/internal/model"
)

// Experience represents work experience entries
type Experience struct {
	model.Base
	ResumeID    uuid.UUID  `json:"resumeId" db:"resume_id"`
	Company     *string    `json:"company" db:"company"`
	Position    *string    `json:"position" db:"position"`
	StartDate   *time.Time `json:"startDate" db:"start_date"`
	EndDate     *time.Time `json:"endDate" db:"end_date"`
	Location    *string    `json:"location" db:"location"`
	Description *string    `json:"description" db:"description"`
	OrderIndex  int        `json:"orderIndex" db:"order_index"`
}
