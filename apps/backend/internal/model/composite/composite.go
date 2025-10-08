package composite

import (
	"github.com/sriniously/go-resumify/internal/model/certification"
	"github.com/sriniously/go-resumify/internal/model/education"
	"github.com/sriniously/go-resumify/internal/model/experience"
	"github.com/sriniously/go-resumify/internal/model/project"
	"github.com/sriniously/go-resumify/internal/model/resume"
	"github.com/sriniously/go-resumify/internal/model/section"
	"github.com/sriniously/go-resumify/internal/model/skill"
)

// ResumeWithSections represents a complete resume with all its sections and items
type ResumeWithSections struct {
	Resume         resume.Resume                 `json:"resume"`
	Sections       []section.ResumeSection       `json:"sections"`
	Education      []education.Education         `json:"education"`
	Experience     []experience.Experience       `json:"experience"`
	Projects       []project.Project             `json:"projects"`
	Skills         []skill.Skill                 `json:"skills"`
	Certifications []certification.Certification `json:"certifications"`
}
