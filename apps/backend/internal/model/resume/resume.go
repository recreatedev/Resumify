package resume

import (
	"github.com/sriniously/go-resumify/internal/model"
)

// Resume represents a user's resume
type Resume struct {
	model.Base
	UserID string `json:"userId" db:"user_id"`
	Title  string `json:"title" db:"title"`
	Theme  string `json:"theme" db:"theme"`
}
