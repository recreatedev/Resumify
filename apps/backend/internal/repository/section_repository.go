package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sriniously/go-resumify/internal/model/section"
	"github.com/sriniously/go-resumify/internal/server"
)

type ResumeSectionRepository struct {
	server *server.Server
}

func NewResumeSectionRepository(server *server.Server) *ResumeSectionRepository {
	return &ResumeSectionRepository{server: server}
}

func (r *ResumeSectionRepository) CreateSection(ctx context.Context, userID string, payload *section.CreateSectionRequest) (*section.ResumeSection, error) {
	stmt := `
		INSERT INTO
			resume_sections (
				resume_id,
				name,
				display_name,
				is_visible,
				order_index
			)
		VALUES
			(
				@resume_id,
				@name,
				@display_name,
				@is_visible,
				@order_index
			)
		RETURNING
		*
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"resume_id":    payload.ResumeID,
		"name":         payload.Name,
		"display_name": payload.DisplayName,
		"is_visible":   payload.IsVisible,
		"order_index":  payload.OrderIndex,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute create section query for resume_id=%s name=%s: %w", payload.ResumeID.String(), payload.Name, err)
	}

	sectionItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[section.ResumeSection])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:resume_sections for resume_id=%s name=%s: %w", payload.ResumeID.String(), payload.Name, err)
	}

	return &sectionItem, nil
}

func (r *ResumeSectionRepository) GetSectionByID(ctx context.Context, userID string, sectionID uuid.UUID) (*section.ResumeSection, error) {
	stmt := `
		SELECT
			rs.*
		FROM
			resume_sections rs
		JOIN resumes r ON rs.resume_id = r.id
		WHERE
			rs.id=@id
			AND r.user_id=@user_id
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"id":      sectionID,
		"user_id": userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get section by id query for section_id=%s user_id=%s: %w", sectionID.String(), userID, err)
	}

	sectionItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[section.ResumeSection])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:resume_sections for section_id=%s user_id=%s: %w", sectionID.String(), userID, err)
	}

	return &sectionItem, nil
}

func (r *ResumeSectionRepository) GetSectionsByResumeID(ctx context.Context, userID string, resumeID uuid.UUID) ([]section.ResumeSection, error) {
	stmt := `
		SELECT
			rs.*
		FROM
			resume_sections rs
		JOIN resumes r ON rs.resume_id = r.id
		WHERE
			rs.resume_id=@resume_id
			AND r.user_id=@user_id
		ORDER BY rs.order_index ASC
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"resume_id": resumeID,
		"user_id":   userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get sections by resume query for resume_id=%s user_id=%s: %w", resumeID.String(), userID, err)
	}

	sections, err := pgx.CollectRows(rows, pgx.RowToStructByName[section.ResumeSection])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []section.ResumeSection{}, nil
		}
		return nil, fmt.Errorf("failed to collect rows from table:resume_sections for resume_id=%s user_id=%s: %w", resumeID.String(), userID, err)
	}

	return sections, nil
}

func (r *ResumeSectionRepository) UpdateSection(ctx context.Context, userID string, sectionID uuid.UUID, payload *section.UpdateSectionRequest) (*section.ResumeSection, error) {
	stmt := `UPDATE resume_sections SET `
	args := pgx.NamedArgs{
		"id": sectionID,
	}
	setClauses := []string{}

	if payload.Name != nil {
		setClauses = append(setClauses, "name = @name")
		args["name"] = *payload.Name
	}
	if payload.DisplayName != nil {
		setClauses = append(setClauses, "display_name = @display_name")
		args["display_name"] = *payload.DisplayName
	}
	if payload.IsVisible != nil {
		setClauses = append(setClauses, "is_visible = @is_visible")
		args["is_visible"] = *payload.IsVisible
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
		return nil, fmt.Errorf("failed to execute update section query for section_id=%s user_id=%s: %w", sectionID.String(), userID, err)
	}

	sectionItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[section.ResumeSection])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:resume_sections for section_id=%s user_id=%s: %w", sectionID.String(), userID, err)
	}

	return &sectionItem, nil
}

func (r *ResumeSectionRepository) BulkUpdateSectionOrder(ctx context.Context, userID string, payload *section.BulkUpdateSectionsRequest) error {
	tx, err := r.server.DB.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	for _, sectionUpdate := range payload.Sections {
		_, err := tx.Exec(ctx, `
			UPDATE resume_sections 
			SET order_index = @order_index
			WHERE id = @id 
			AND resume_id IN (SELECT id FROM resumes WHERE user_id = @user_id)
		`, pgx.NamedArgs{
			"id":          sectionUpdate.ID,
			"order_index": sectionUpdate.OrderIndex,
			"user_id":     userID,
		})
		if err != nil {
			return fmt.Errorf("failed to update section order for section_id=%s: %w", sectionUpdate.ID, err)
		}
	}

	return tx.Commit(ctx)
}

func (r *ResumeSectionRepository) DeleteSection(ctx context.Context, userID string, sectionID uuid.UUID) error {
	result, err := r.server.DB.Pool.Exec(ctx, `
		DELETE FROM resume_sections
		WHERE id = @id 
		AND resume_id IN (SELECT id FROM resumes WHERE user_id = @user_id)
	`, pgx.NamedArgs{
		"id":      sectionID,
		"user_id": userID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete section: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("section not found")
	}

	return nil
}
