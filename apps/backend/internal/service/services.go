package service

import (
	"github.com/recreatedev/Resumify/internal/lib/job"
	"github.com/recreatedev/Resumify/internal/repository"
	"github.com/recreatedev/Resumify/internal/server"
)

type Services struct {
	Auth          *AuthService
	Resume        *ResumeService
	Education     *EducationService
	Experience    *ExperienceService
	Project       *ProjectService
	Skill         *SkillService
	Certification *CertificationService
	Section       *SectionService
	Job           *job.JobService
}

func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error) {
	authService := NewAuthService(s)
	resumeService := NewResumeService(s, repos)
	educationService := NewEducationService(s, repos)
	experienceService := NewExperienceService(s, repos)
	projectService := NewProjectService(s, repos)
	skillService := NewSkillService(s, repos)
	certificationService := NewCertificationService(s, repos)
	sectionService := NewSectionService(s, repos)

	return &Services{
		Job:           s.Job,
		Auth:          authService,
		Resume:        resumeService,
		Education:     educationService,
		Experience:    experienceService,
		Project:       projectService,
		Skill:         skillService,
		Certification: certificationService,
		Section:       sectionService,
	}, nil
}
