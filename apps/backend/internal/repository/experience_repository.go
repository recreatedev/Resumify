package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sriniously/go-resumify/internal/model/experience"
	"github.com/sriniously/go-resumify/internal/server"
)

type ExperienceRepository struct {
	server *server.Server
}

func NewExperienceRepository(server *server.Server) *ExperienceRepository {
	return &ExperienceRepository{server: server}
}

func (r *ExperienceRepository) CreateExperience(ctx context.Context, userID string, payload *experience.CreateExperienceRequest) (*experience.Experience, error) {
	stmt := `
		INSERT INTO
			experience (
				resume_id,
				company,
				position,
				start_date,
				end_date,
				location,
				description,
				order_index
			)
		VALUES
			(
				@resume_id,
				@company,
				@position,
				@start_date,
				@end_date,
				@location,
				@description,
				@order_index
			)
		RETURNING
		*
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"resume_id":   payload.ResumeID,
		"company":     payload.Company,
		"position":    payload.Position,
		"start_date":  payload.StartDate,
		"end_date":    payload.EndDate,
		"location":    payload.Location,
		"description": payload.Description,
		"order_index": payload.OrderIndex,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute create experience query for resume_id=%s: %w", payload.ResumeID.String(), err)
	}

	experienceItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[experience.Experience])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:experience for resume_id=%s: %w", payload.ResumeID.String(), err)
	}

	return &experienceItem, nil
}

func (r *ExperienceRepository) GetExperienceByID(ctx context.Context, userID string, experienceID uuid.UUID) (*experience.Experience, error) {
	stmt := `
		SELECT
			e.*
		FROM
			experience e
		JOIN resumes r ON e.resume_id = r.id
		WHERE
			e.id=@id
			AND r.user_id=@user_id
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"id":      experienceID,
		"user_id": userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get experience by id query for experience_id=%s user_id=%s: %w", experienceID.String(), userID, err)
	}

	experienceItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[experience.Experience])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:experience for experience_id=%s user_id=%s: %w", experienceID.String(), userID, err)
	}

	return &experienceItem, nil
}

func (r *ExperienceRepository) GetExperienceByResumeID(ctx context.Context, userID string, resumeID uuid.UUID) ([]experience.Experience, error) {
	stmt := `
		SELECT
			e.*
		FROM
			experience e
		JOIN resumes r ON e.resume_id = r.id
		WHERE
			e.resume_id=@resume_id
			AND r.user_id=@user_id
		ORDER BY e.order_index ASC
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"resume_id": resumeID,
		"user_id":   userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get experience by resume query for resume_id=%s user_id=%s: %w", resumeID.String(), userID, err)
	}

	experienceItems, err := pgx.CollectRows(rows, pgx.RowToStructByName[experience.Experience])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []experience.Experience{}, nil
		}
		return nil, fmt.Errorf("failed to collect rows from table:experience for resume_id=%s user_id=%s: %w", resumeID.String(), userID, err)
	}

	return experienceItems, nil
}

func (r *ExperienceRepository) UpdateExperience(ctx context.Context, userID string, experienceID uuid.UUID, payload *experience.UpdateExperienceRequest) (*experience.Experience, error) {
	stmt := `UPDATE experience SET `
	args := pgx.NamedArgs{
		"id": experienceID,
	}
	setClauses := []string{}

	if payload.Company != nil {
		setClauses = append(setClauses, "company = @company")
		args["company"] = *payload.Company
	}
	if payload.Position != nil {
		setClauses = append(setClauses, "position = @position")
		args["position"] = *payload.Position
	}
	if payload.StartDate != nil {
		setClauses = append(setClauses, "start_date = @start_date")
		args["start_date"] = *payload.StartDate
	}
	if payload.EndDate != nil {
		setClauses = append(setClauses, "end_date = @end_date")
		args["end_date"] = *payload.EndDate
	}
	if payload.Location != nil {
		setClauses = append(setClauses, "location = @location")
		args["location"] = *payload.Location
	}
	if payload.Description != nil {
		setClauses = append(setClauses, "description = @description")
		args["description"] = *payload.Description
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
		return nil, fmt.Errorf("failed to execute update experience query for experience_id=%s user_id=%s: %w", experienceID.String(), userID, err)
	}

	experienceItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[experience.Experience])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:experience for experience_id=%s user_id=%s: %w", experienceID.String(), userID, err)
	}

	return &experienceItem, nil
}

func (r *ExperienceRepository) BulkUpdateExperienceOrder(ctx context.Context, userID string, payload *experience.BulkUpdateExperienceRequest) error {
	tx, err := r.server.DB.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	for _, experienceUpdate := range payload.Experience {
		_, err := tx.Exec(ctx, `
			UPDATE experience 
			SET order_index = @order_index
			WHERE id = @id 
			AND resume_id IN (SELECT id FROM resumes WHERE user_id = @user_id)
		`, pgx.NamedArgs{
			"id":          experienceUpdate.ID,
			"order_index": experienceUpdate.OrderIndex,
			"user_id":     userID,
		})
		if err != nil {
			return fmt.Errorf("failed to update experience order for experience_id=%s: %w", experienceUpdate.ID, err)
		}
	}

	return tx.Commit(ctx)
}

func (r *ExperienceRepository) DeleteExperience(ctx context.Context, userID string, experienceID uuid.UUID) error {
	result, err := r.server.DB.Pool.Exec(ctx, `
		DELETE FROM experience
		WHERE id = @id 
		AND resume_id IN (SELECT id FROM resumes WHERE user_id = @user_id)
	`, pgx.NamedArgs{
		"id":      experienceID,
		"user_id": userID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete experience: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("experience not found")
	}

	return nil
}
