package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/recreatedev/Resumify/internal/model/education"
	"github.com/recreatedev/Resumify/internal/server"
)

type EducationRepository struct {
	server *server.Server
}

func NewEducationRepository(server *server.Server) *EducationRepository {
	return &EducationRepository{server: server}
}

func (r *EducationRepository) CreateEducation(ctx context.Context, userID string, payload *education.CreateEducationRequest) (*education.Education, error) {
	stmt := `
		INSERT INTO
			education (
				resume_id,
				institution,
				degree,
				field_of_study,
				start_date,
				end_date,
				grade,
				description,
				order_index
			)
		VALUES
			(
				@resume_id,
				@institution,
				@degree,
				@field_of_study,
				@start_date,
				@end_date,
				@grade,
				@description,
				@order_index
			)
		RETURNING
		*
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"resume_id":      payload.ResumeID,
		"institution":    payload.Institution,
		"degree":         payload.Degree,
		"field_of_study": payload.FieldOfStudy,
		"start_date":     payload.StartDate,
		"end_date":       payload.EndDate,
		"grade":          payload.Grade,
		"description":    payload.Description,
		"order_index":    payload.OrderIndex,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute create education query for resume_id=%s: %w", payload.ResumeID.String(), err)
	}

	educationItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[education.Education])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:education for resume_id=%s: %w", payload.ResumeID.String(), err)
	}

	return &educationItem, nil
}

func (r *EducationRepository) GetEducationByID(ctx context.Context, userID string, educationID uuid.UUID) (*education.Education, error) {
	stmt := `
		SELECT
			e.*
		FROM
			education e
		JOIN resumes r ON e.resume_id = r.id
		WHERE
			e.id=@id
			AND r.user_id=@user_id
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"id":      educationID,
		"user_id": userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get education by id query for education_id=%s user_id=%s: %w", educationID.String(), userID, err)
	}

	educationItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[education.Education])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:education for education_id=%s user_id=%s: %w", educationID.String(), userID, err)
	}

	return &educationItem, nil
}

func (r *EducationRepository) GetEducationByResumeID(ctx context.Context, userID string, resumeID uuid.UUID) ([]education.Education, error) {
	stmt := `
		SELECT
			e.*
		FROM
			education e
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
		return nil, fmt.Errorf("failed to execute get education by resume query for resume_id=%s user_id=%s: %w", resumeID.String(), userID, err)
	}

	educationItems, err := pgx.CollectRows(rows, pgx.RowToStructByName[education.Education])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []education.Education{}, nil
		}
		return nil, fmt.Errorf("failed to collect rows from table:education for resume_id=%s user_id=%s: %w", resumeID.String(), userID, err)
	}

	return educationItems, nil
}

func (r *EducationRepository) UpdateEducation(ctx context.Context, userID string, educationID uuid.UUID, payload *education.UpdateEducationRequest) (*education.Education, error) {
	stmt := `UPDATE education SET `
	args := pgx.NamedArgs{
		"id": educationID,
	}
	setClauses := []string{}

	if payload.Institution != nil {
		setClauses = append(setClauses, "institution = @institution")
		args["institution"] = *payload.Institution
	}
	if payload.Degree != nil {
		setClauses = append(setClauses, "degree = @degree")
		args["degree"] = *payload.Degree
	}
	if payload.FieldOfStudy != nil {
		setClauses = append(setClauses, "field_of_study = @field_of_study")
		args["field_of_study"] = *payload.FieldOfStudy
	}
	if payload.StartDate != nil {
		setClauses = append(setClauses, "start_date = @start_date")
		args["start_date"] = *payload.StartDate
	}
	if payload.EndDate != nil {
		setClauses = append(setClauses, "end_date = @end_date")
		args["end_date"] = *payload.EndDate
	}
	if payload.Grade != nil {
		setClauses = append(setClauses, "grade = @grade")
		args["grade"] = *payload.Grade
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
		return nil, fmt.Errorf("failed to execute update education query for education_id=%s user_id=%s: %w", educationID.String(), userID, err)
	}

	educationItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[education.Education])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:education for education_id=%s user_id=%s: %w", educationID.String(), userID, err)
	}

	return &educationItem, nil
}

func (r *EducationRepository) BulkUpdateEducationOrder(ctx context.Context, userID string, payload *education.BulkUpdateEducationRequest) error {
	tx, err := r.server.DB.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	for _, educationUpdate := range payload.Education {
		_, err := tx.Exec(ctx, `
			UPDATE education 
			SET order_index = @order_index
			WHERE id = @id 
			AND resume_id IN (SELECT id FROM resumes WHERE user_id = @user_id)
		`, pgx.NamedArgs{
			"id":          educationUpdate.ID,
			"order_index": educationUpdate.OrderIndex,
			"user_id":     userID,
		})
		if err != nil {
			return fmt.Errorf("failed to update education order for education_id=%s: %w", educationUpdate.ID, err)
		}
	}

	return tx.Commit(ctx)
}

func (r *EducationRepository) DeleteEducation(ctx context.Context, userID string, educationID uuid.UUID) error {
	result, err := r.server.DB.Pool.Exec(ctx, `
		DELETE FROM education
		WHERE id = @id 
		AND resume_id IN (SELECT id FROM resumes WHERE user_id = @user_id)
	`, pgx.NamedArgs{
		"id":      educationID,
		"user_id": userID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete education: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("education not found")
	}

	return nil
}
