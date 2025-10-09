package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/recreatedev/Resumify/internal/handler"
	"github.com/recreatedev/Resumify/internal/middleware"
	"github.com/recreatedev/Resumify/internal/server"
	"github.com/recreatedev/Resumify/internal/service"
)

func RegisterRoutes(e *echo.Echo, h *handler.Handlers, s *server.Server, services *service.Services) {
	v1 := e.Group("/api/v1")

	// Apply authentication middleware to all v1 routes
	authMiddleware := middleware.NewAuthMiddleware(s)
	v1.Use(authMiddleware.RequireAuth)

	// Resume routes
	registerResumeRoutes(v1, h)

	// Education routes
	registerEducationRoutes(v1, h)

	// Experience routes
	registerExperienceRoutes(v1, h)

	// Project routes
	registerProjectRoutes(v1, h)

	// Skill routes
	registerSkillRoutes(v1, h)

	// Certification routes
	registerCertificationRoutes(v1, h)

	// Section routes
	registerSectionRoutes(v1, h)
}

func registerResumeRoutes(g *echo.Group, h *handler.Handlers) {
	resumes := g.Group("/resumes")

	// Resume CRUD operations
	resumes.POST("", h.Resume.CreateResume)
	resumes.GET("", h.Resume.GetResumes)
	resumes.GET("/:id", h.Resume.GetResumeByID)
	resumes.PUT("/:id", h.Resume.UpdateResume)
	resumes.DELETE("/:id", h.Resume.DeleteResume)

	// Resume operations
	resumes.POST("/:id/duplicate", h.Resume.DuplicateResume)
	resumes.GET("/:id/sections", h.Resume.GetResumeWithSections)
}

func registerEducationRoutes(g *echo.Group, h *handler.Handlers) {
	educations := g.Group("/educations")

	// Education CRUD operations
	educations.POST("", h.Education.CreateEducation)
	educations.GET("/:id", h.Education.GetEducationByID)
	educations.PUT("/:id", h.Education.UpdateEducation)
	educations.DELETE("/:id", h.Education.DeleteEducation)

	// Education bulk operations
	educations.PUT("/order", h.Education.BulkUpdateEducationOrder)

	// Resume-specific education routes
	resumes := g.Group("/resumes")
	resumes.GET("/:resumeId/educations", h.Education.GetEducationByResumeID)
}

func registerExperienceRoutes(g *echo.Group, h *handler.Handlers) {
	experiences := g.Group("/experiences")

	// Experience CRUD operations
	experiences.POST("", h.Experience.CreateExperience)
	experiences.GET("/:id", h.Experience.GetExperienceByID)
	experiences.PUT("/:id", h.Experience.UpdateExperience)
	experiences.DELETE("/:id", h.Experience.DeleteExperience)

	// Experience bulk operations
	experiences.PUT("/order", h.Experience.BulkUpdateExperienceOrder)

	// Resume-specific experience routes
	resumes := g.Group("/resumes")
	resumes.GET("/:resumeId/experiences", h.Experience.GetExperienceByResumeID)
}

func registerProjectRoutes(g *echo.Group, h *handler.Handlers) {
	projects := g.Group("/projects")

	// Project CRUD operations
	projects.POST("", h.Project.CreateProject)
	projects.GET("/:id", h.Project.GetProjectByID)
	projects.PUT("/:id", h.Project.UpdateProject)
	projects.DELETE("/:id", h.Project.DeleteProject)

	// Project bulk operations
	projects.PUT("/order", h.Project.BulkUpdateProjectOrder)

	// Resume-specific project routes
	resumes := g.Group("/resumes")
	resumes.GET("/:resumeId/projects", h.Project.GetProjectsByResumeID)
}

func registerSkillRoutes(g *echo.Group, h *handler.Handlers) {
	skills := g.Group("/skills")

	// Skill CRUD operations
	skills.POST("", h.Skill.CreateSkill)
	skills.GET("/:id", h.Skill.GetSkillByID)
	skills.PUT("/:id", h.Skill.UpdateSkill)
	skills.DELETE("/:id", h.Skill.DeleteSkill)

	// Skill bulk operations
	skills.PUT("/order", h.Skill.BulkUpdateSkillOrder)

	// Resume-specific skill routes
	resumes := g.Group("/resumes")
	resumes.GET("/:resumeId/skills", h.Skill.GetSkillsByResumeID)
	resumes.GET("/:resumeId/skills/category", h.Skill.GetSkillsByCategory)
}

func registerCertificationRoutes(g *echo.Group, h *handler.Handlers) {
	certifications := g.Group("/certifications")

	// Certification CRUD operations
	certifications.POST("", h.Certification.CreateCertification)
	certifications.GET("/:id", h.Certification.GetCertificationByID)
	certifications.PUT("/:id", h.Certification.UpdateCertification)
	certifications.DELETE("/:id", h.Certification.DeleteCertification)

	// Certification bulk operations
	certifications.PUT("/order", h.Certification.BulkUpdateCertificationOrder)

	// Resume-specific certification routes
	resumes := g.Group("/resumes")
	resumes.GET("/:resumeId/certifications", h.Certification.GetCertificationsByResumeID)
}

func registerSectionRoutes(g *echo.Group, h *handler.Handlers) {
	sections := g.Group("/sections")

	// Section CRUD operations
	sections.POST("", h.Section.CreateSection)
	sections.GET("/:id", h.Section.GetSectionByID)
	sections.PUT("/:id", h.Section.UpdateSection)
	sections.DELETE("/:id", h.Section.DeleteSection)

	// Section bulk operations
	sections.PUT("/order", h.Section.BulkUpdateSectionOrder)

	// Resume-specific section routes
	resumes := g.Group("/resumes")
	resumes.GET("/:resumeId/sections", h.Section.GetSectionsByResumeID)
}
