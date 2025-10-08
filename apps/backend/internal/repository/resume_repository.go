package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/recreatedev/Resumify/internal/model"
	"github.com/recreatedev/Resumify/internal/model/resume"
	"github.com/recreatedev/Resumify/internal/server"
)

type ResumeRepository struct {
	server *server.Server
}

func NewResumeRepository(server *server.Server) *ResumeRepository {
	return &ResumeRepository{server: server}
}

func (r *ResumeRepository) CreateResume(ctx context.Context, userID string, payload *resume.CreateResumeRequest) (*resume.Resume, error) {
	stmt := `
		INSERT INTO
			resumes (
				user_id,
				title,
				theme
			)
		VALUES
			(
				@user_id,
				@title,
				@theme
			)
		RETURNING
		*
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"user_id": userID,
		"title":   payload.Title,
		"theme":   payload.Theme,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute create resume query for user_id=%s title=%s: %w", userID, payload.Title, err)
	}

	resumeItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[resume.Resume])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:resumes for user_id=%s title=%s: %w", userID, payload.Title, err)
	}

	return &resumeItem, nil
}

func (r *ResumeRepository) GetResumeByID(ctx context.Context, userID string, resumeID uuid.UUID) (*resume.Resume, error) {
	stmt := `
		SELECT
			*
		FROM
			resumes
		WHERE
			id=@id
			AND user_id=@user_id
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"id":      resumeID,
		"user_id": userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get resume by id query for resume_id=%s user_id=%s: %w", resumeID.String(), userID, err)
	}

	resumeItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[resume.Resume])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:resumes for resume_id=%s user_id=%s: %w", resumeID.String(), userID, err)
	}

	return &resumeItem, nil
}

func (r *ResumeRepository) GetResumes(ctx context.Context, userID string, page, limit int) (*model.PaginatedResponse[resume.Resume], error) {
	stmt := `
		SELECT
			*
		FROM
			resumes
		WHERE
			user_id=@user_id
		ORDER BY created_at DESC
		LIMIT @limit OFFSET @offset
	`

	args := pgx.NamedArgs{
		"user_id": userID,
		"limit":   limit,
		"offset":  (page - 1) * limit,
	}

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get resumes query for user_id=%s: %w", userID, err)
	}

	resumes, err := pgx.CollectRows(rows, pgx.RowToStructByName[resume.Resume])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.PaginatedResponse[resume.Resume]{
				Data:       []resume.Resume{},
				Page:       page,
				Limit:      limit,
				Total:      0,
				TotalPages: 0,
			}, nil
		}
		return nil, fmt.Errorf("failed to collect rows from table:resumes for user_id=%s: %w", userID, err)
	}

	// Get total count
	countStmt := `
		SELECT
			COUNT(*)
		FROM
			resumes
		WHERE
			user_id=@user_id
	`

	var total int
	err = r.server.DB.Pool.QueryRow(ctx, countStmt, pgx.NamedArgs{"user_id": userID}).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count of resumes for user_id=%s: %w", userID, err)
	}

	return &model.PaginatedResponse[resume.Resume]{
		Data:       resumes,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: (total + limit - 1) / limit,
	}, nil
}

func (r *ResumeRepository) UpdateResume(ctx context.Context, userID string, resumeID uuid.UUID, payload *resume.UpdateResumeRequest) (*resume.Resume, error) {
	stmt := `UPDATE resumes SET `
	args := pgx.NamedArgs{
		"id":      resumeID,
		"user_id": userID,
	}
	setClauses := []string{}

	if payload.Title != nil {
		setClauses = append(setClauses, "title = @title")
		args["title"] = *payload.Title
	}
	if payload.Theme != nil {
		setClauses = append(setClauses, "theme = @theme")
		args["theme"] = *payload.Theme
	}

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	stmt += strings.Join(setClauses, ", ")
	stmt += ` WHERE id = @id AND user_id = @user_id RETURNING *`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute update resume query for resume_id=%s user_id=%s: %w", resumeID.String(), userID, err)
	}

	resumeItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[resume.Resume])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:resumes for resume_id=%s user_id=%s: %w", resumeID.String(), userID, err)
	}

	return &resumeItem, nil
}

func (r *ResumeRepository) DeleteResume(ctx context.Context, userID string, resumeID uuid.UUID) error {
	result, err := r.server.DB.Pool.Exec(ctx, `
		DELETE FROM resumes
		WHERE id = @id AND user_id = @user_id
	`, pgx.NamedArgs{
		"id":      resumeID,
		"user_id": userID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete resume: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("resume not found")
	}

	return nil
}
