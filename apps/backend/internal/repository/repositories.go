package repository

import "github.com/recreatedev/Resumify/internal/server"

type Repositories struct {
	Resume        *ResumeRepository
	Section       *ResumeSectionRepository
	Education     *EducationRepository
	Experience    *ExperienceRepository
	Project       *ProjectRepository
	Skill         *SkillRepository
	Certification *CertificationRepository
}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{
		Resume:        NewResumeRepository(s),
		Section:       NewResumeSectionRepository(s),
		Education:     NewEducationRepository(s),
		Experience:    NewExperienceRepository(s),
		Project:       NewProjectRepository(s),
		Skill:         NewSkillRepository(s),
		Certification: NewCertificationRepository(s),
	}
}
