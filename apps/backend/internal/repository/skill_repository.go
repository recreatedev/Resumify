package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/recreatedev/Resumify/internal/model/skill"
	"github.com/recreatedev/Resumify/internal/server"
)

type SkillRepository struct {
	server *server.Server
}

func NewSkillRepository(server *server.Server) *SkillRepository {
	return &SkillRepository{server: server}
}

func (r *SkillRepository) CreateSkill(ctx context.Context, userID string, payload *skill.CreateSkillRequest) (*skill.Skill, error) {
	stmt := `
		INSERT INTO
			skills (
				resume_id,
				name,
				level,
				category,
				order_index
			)
		VALUES
			(
				@resume_id,
				@name,
				@level,
				@category,
				@order_index
			)
		RETURNING
		*
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"resume_id":   payload.ResumeID,
		"name":        payload.Name,
		"level":       payload.Level,
		"category":    payload.Category,
		"order_index": payload.OrderIndex,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute create skill query for resume_id=%s: %w", payload.ResumeID.String(), err)
	}

	skillItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[skill.Skill])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:skills for resume_id=%s: %w", payload.ResumeID.String(), err)
	}

	return &skillItem, nil
}

func (r *SkillRepository) GetSkillByID(ctx context.Context, userID string, skillID uuid.UUID) (*skill.Skill, error) {
	stmt := `
		SELECT
			s.*
		FROM
			skills s
		JOIN resumes r ON s.resume_id = r.id
		WHERE
			s.id=@id
			AND r.user_id=@user_id
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"id":      skillID,
		"user_id": userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get skill by id query for skill_id=%s user_id=%s: %w", skillID.String(), userID, err)
	}

	skillItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[skill.Skill])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:skills for skill_id=%s user_id=%s: %w", skillID.String(), userID, err)
	}

	return &skillItem, nil
}

func (r *SkillRepository) GetSkillsByResumeID(ctx context.Context, userID string, resumeID uuid.UUID) ([]skill.Skill, error) {
	stmt := `
		SELECT
			s.*
		FROM
			skills s
		JOIN resumes r ON s.resume_id = r.id
		WHERE
			s.resume_id=@resume_id
			AND r.user_id=@user_id
		ORDER BY s.order_index ASC
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"resume_id": resumeID,
		"user_id":   userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get skills by resume query for resume_id=%s user_id=%s: %w", resumeID.String(), userID, err)
	}

	skillItems, err := pgx.CollectRows(rows, pgx.RowToStructByName[skill.Skill])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []skill.Skill{}, nil
		}
		return nil, fmt.Errorf("failed to collect rows from table:skills for resume_id=%s user_id=%s: %w", resumeID.String(), userID, err)
	}

	return skillItems, nil
}

func (r *SkillRepository) GetSkillsByCategory(ctx context.Context, userID string, resumeID uuid.UUID) ([]skill.SkillsByCategoryResponse, error) {
	stmt := `
		SELECT
			s.*
		FROM
			skills s
		JOIN resumes r ON s.resume_id = r.id
		WHERE
			s.resume_id=@resume_id
			AND r.user_id=@user_id
		ORDER BY s.category ASC, s.order_index ASC
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"resume_id": resumeID,
		"user_id":   userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get skills by category query for resume_id=%s user_id=%s: %w", resumeID.String(), userID, err)
	}

	skillItems, err := pgx.CollectRows(rows, pgx.RowToStructByName[skill.Skill])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []skill.SkillsByCategoryResponse{}, nil
		}
		return nil, fmt.Errorf("failed to collect rows from table:skills for resume_id=%s user_id=%s: %w", resumeID.String(), userID, err)
	}

	// Group skills by category
	categoryMap := make(map[string][]skill.SkillResponse)
	for _, skillItem := range skillItems {
		category := "Other"
		if skillItem.Category != nil {
			category = *skillItem.Category
		}

		// Convert Skill to SkillResponse
		skillResponse := skill.SkillResponse{
			ID:         skillItem.ID.String(),
			ResumeID:   skillItem.ResumeID,
			Name:       skillItem.Name,
			Level:      skillItem.Level,
			Category:   skillItem.Category,
			OrderIndex: skillItem.OrderIndex,
		}

		categoryMap[category] = append(categoryMap[category], skillResponse)
	}

	// Convert to response format
	var result []skill.SkillsByCategoryResponse
	for category, skills := range categoryMap {
		result = append(result, skill.SkillsByCategoryResponse{
			Category: category,
			Skills:   skills,
		})
	}

	return result, nil
}

func (r *SkillRepository) UpdateSkill(ctx context.Context, userID string, skillID uuid.UUID, payload *skill.UpdateSkillRequest) (*skill.Skill, error) {
	stmt := `UPDATE skills SET `
	args := pgx.NamedArgs{
		"id": skillID,
	}
	setClauses := []string{}

	if payload.Name != nil {
		setClauses = append(setClauses, "name = @name")
		args["name"] = *payload.Name
	}
	if payload.Level != nil {
		setClauses = append(setClauses, "level = @level")
		args["level"] = *payload.Level
	}
	if payload.Category != nil {
		setClauses = append(setClauses, "category = @category")
		args["category"] = *payload.Category
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
		return nil, fmt.Errorf("failed to execute update skill query for skill_id=%s user_id=%s: %w", skillID.String(), userID, err)
	}

	skillItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[skill.Skill])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:skills for skill_id=%s user_id=%s: %w", skillID.String(), userID, err)
	}

	return &skillItem, nil
}

func (r *SkillRepository) BulkUpdateSkillOrder(ctx context.Context, userID string, payload *skill.BulkUpdateSkillsRequest) error {
	tx, err := r.server.DB.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	for _, skillUpdate := range payload.Skills {
		_, err := tx.Exec(ctx, `
			UPDATE skills 
			SET order_index = @order_index
			WHERE id = @id 
			AND resume_id IN (SELECT id FROM resumes WHERE user_id = @user_id)
		`, pgx.NamedArgs{
			"id":          skillUpdate.ID,
			"order_index": skillUpdate.OrderIndex,
			"user_id":     userID,
		})
		if err != nil {
			return fmt.Errorf("failed to update skill order for skill_id=%s: %w", skillUpdate.ID, err)
		}
	}

	return tx.Commit(ctx)
}

func (r *SkillRepository) DeleteSkill(ctx context.Context, userID string, skillID uuid.UUID) error {
	result, err := r.server.DB.Pool.Exec(ctx, `
		DELETE FROM skills
		WHERE id = @id 
		AND resume_id IN (SELECT id FROM resumes WHERE user_id = @user_id)
	`, pgx.NamedArgs{
		"id":      skillID,
		"user_id": userID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete skill: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("skill not found")
	}

	return nil
}
