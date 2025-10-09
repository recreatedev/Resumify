package handler

import (
	"github.com/recreatedev/Resumify/internal/server"
	"github.com/recreatedev/Resumify/internal/service"
)

type Handlers struct {
	Health        *HealthHandler
	Resume        *ResumeHandler
	Education     *EducationHandler
	Experience    *ExperienceHandler
	Project       *ProjectHandler
	Skill         *SkillHandler
	Certification *CertificationHandler
	Section       *SectionHandler
	OpenAPI       *OpenAPIHandler
}

func NewHandlers(s *server.Server, services *service.Services) *Handlers {
	return &Handlers{
		Health:        NewHealthHandler(s),
		Resume:        NewResumeHandler(s, services),
		Education:     NewEducationHandler(s, services.Education),
		Experience:    NewExperienceHandler(s, services.Experience),
		Project:       NewProjectHandler(s, services.Project),
		Skill:         NewSkillHandler(s, services.Skill),
		Certification: NewCertificationHandler(s, services.Certification),
		Section:       NewSectionHandler(s, services.Section),
		OpenAPI:       NewOpenAPIHandler(s),
	}
}
