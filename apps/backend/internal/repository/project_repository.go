package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/recreatedev/Resumify/internal/model/project"
	"github.com/recreatedev/Resumify/internal/server"
)

type ProjectRepository struct {
	server *server.Server
}

func NewProjectRepository(server *server.Server) *ProjectRepository {
	return &ProjectRepository{server: server}
}

func (r *ProjectRepository) CreateProject(ctx context.Context, userID string, payload *project.CreateProjectRequest) (*project.Project, error) {
	stmt := `
		INSERT INTO
			projects (
				resume_id,
				name,
				role,
				description,
				link,
				technologies,
				order_index
			)
		VALUES
			(
				@resume_id,
				@name,
				@role,
				@description,
				@link,
				@technologies,
				@order_index
			)
		RETURNING
		*
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"resume_id":    payload.ResumeID,
		"name":         payload.Name,
		"role":         payload.Role,
		"description":  payload.Description,
		"link":         payload.Link,
		"technologies": payload.Technologies,
		"order_index":  payload.OrderIndex,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute create project query for resume_id=%s: %w", payload.ResumeID.String(), err)
	}

	projectItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[project.Project])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:projects for resume_id=%s: %w", payload.ResumeID.String(), err)
	}

	return &projectItem, nil
}

func (r *ProjectRepository) GetProjectByID(ctx context.Context, userID string, projectID uuid.UUID) (*project.Project, error) {
	stmt := `
		SELECT
			p.*
		FROM
			projects p
		JOIN resumes r ON p.resume_id = r.id
		WHERE
			p.id=@id
			AND r.user_id=@user_id
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"id":      projectID,
		"user_id": userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get project by id query for project_id=%s user_id=%s: %w", projectID.String(), userID, err)
	}

	projectItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[project.Project])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:projects for project_id=%s user_id=%s: %w", projectID.String(), userID, err)
	}

	return &projectItem, nil
}

func (r *ProjectRepository) GetProjectsByResumeID(ctx context.Context, userID string, resumeID uuid.UUID) ([]project.Project, error) {
	stmt := `
		SELECT
			p.*
		FROM
			projects p
		JOIN resumes r ON p.resume_id = r.id
		WHERE
			p.resume_id=@resume_id
			AND r.user_id=@user_id
		ORDER BY p.order_index ASC
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"resume_id": resumeID,
		"user_id":   userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get projects by resume query for resume_id=%s user_id=%s: %w", resumeID.String(), userID, err)
	}

	projectItems, err := pgx.CollectRows(rows, pgx.RowToStructByName[project.Project])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []project.Project{}, nil
		}
		return nil, fmt.Errorf("failed to collect rows from table:projects for resume_id=%s user_id=%s: %w", resumeID.String(), userID, err)
	}

	return projectItems, nil
}

func (r *ProjectRepository) UpdateProject(ctx context.Context, userID string, projectID uuid.UUID, payload *project.UpdateProjectRequest) (*project.Project, error) {
	stmt := `UPDATE projects SET `
	args := pgx.NamedArgs{
		"id": projectID,
	}
	setClauses := []string{}

	if payload.Name != nil {
		setClauses = append(setClauses, "name = @name")
		args["name"] = *payload.Name
	}
	if payload.Role != nil {
		setClauses = append(setClauses, "role = @role")
		args["role"] = *payload.Role
	}
	if payload.Description != nil {
		setClauses = append(setClauses, "description = @description")
		args["description"] = *payload.Description
	}
	if payload.Link != nil {
		setClauses = append(setClauses, "link = @link")
		args["link"] = *payload.Link
	}
	if payload.Technologies != nil {
		setClauses = append(setClauses, "technologies = @technologies")
		args["technologies"] = payload.Technologies
	}
	if payload.OrderIndex != nil {
		setClauses = append(setClauses, "order_index = @order_index")
		args["order_index"] = *payload.OrderIndex
	}

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	stmt += strings.Join(setClauses, ", ")
	stmt += ` WHERE id = @id AND resume_id IN (SELECT id FROM resumes WHERE user_id = @user_id) RETURNING *`

	args["user_id"] = userID

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute update project query for project_id=%s user_id=%s: %w", projectID.String(), userID, err)
	}

	projectItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[project.Project])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:projects for project_id=%s user_id=%s: %w", projectID.String(), userID, err)
	}

	return &projectItem, nil
}

func (r *ProjectRepository) BulkUpdateProjectOrder(ctx context.Context, userID string, payload *project.BulkUpdateProjectsRequest) error {
	tx, err := r.server.DB.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	for _, projectUpdate := range payload.Projects {
		_, err := tx.Exec(ctx, `
			UPDATE projects 
			SET order_index = @order_index
			WHERE id = @id 
			AND resume_id IN (SELECT id FROM resumes WHERE user_id = @user_id)
		`, pgx.NamedArgs{
			"id":          projectUpdate.ID,
			"order_index": projectUpdate.OrderIndex,
			"user_id":     userID,
		})
		if err != nil {
			return fmt.Errorf("failed to update project order for project_id=%s: %w", projectUpdate.ID, err)
		}
	}

	return tx.Commit(ctx)
}

func (r *ProjectRepository) DeleteProject(ctx context.Context, userID string, projectID uuid.UUID) error {
	result, err := r.server.DB.Pool.Exec(ctx, `
		DELETE FROM projects
		WHERE id = @id 
		AND resume_id IN (SELECT id FROM resumes WHERE user_id = @user_id)
	`, pgx.NamedArgs{
		"id":      projectID,
		"user_id": userID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("project not found")
	}

	return nil
}
