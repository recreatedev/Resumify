package composite

import (
	"github.com/recreatedev/Resumify/internal/model/certification"
	"github.com/recreatedev/Resumify/internal/model/education"
	"github.com/recreatedev/Resumify/internal/model/experience"
	"github.com/recreatedev/Resumify/internal/model/project"
	"github.com/recreatedev/Resumify/internal/model/resume"
	"github.com/recreatedev/Resumify/internal/model/section"
	"github.com/recreatedev/Resumify/internal/model/skill"
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
