package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/recreatedev/Resumify/internal/model/certification"
	"github.com/recreatedev/Resumify/internal/server"
)

type CertificationRepository struct {
	server *server.Server
}

func NewCertificationRepository(server *server.Server) *CertificationRepository {
	return &CertificationRepository{server: server}
}

func (r *CertificationRepository) CreateCertification(ctx context.Context, userID string, payload *certification.CreateCertificationRequest) (*certification.Certification, error) {
	stmt := `
		INSERT INTO
			certifications (
				resume_id,
				name,
				organization,
				issue_date,
				expiry_date,
				credential_id,
				credential_url,
				order_index
			)
		VALUES
			(
				@resume_id,
				@name,
				@organization,
				@issue_date,
				@expiry_date,
				@credential_id,
				@credential_url,
				@order_index
			)
		RETURNING
		*
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"resume_id":      payload.ResumeID,
		"name":           payload.Name,
		"organization":   payload.Organization,
		"issue_date":     payload.IssueDate,
		"expiry_date":    payload.ExpiryDate,
		"credential_id":  payload.CredentialID,
		"credential_url": payload.CredentialURL,
		"order_index":    payload.OrderIndex,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute create certification query for resume_id=%s: %w", payload.ResumeID.String(), err)
	}

	certificationItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[certification.Certification])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:certifications for resume_id=%s: %w", payload.ResumeID.String(), err)
	}

	return &certificationItem, nil
}

func (r *CertificationRepository) GetCertificationByID(ctx context.Context, userID string, certificationID uuid.UUID) (*certification.Certification, error) {
	stmt := `
		SELECT
			c.*
		FROM
			certifications c
		JOIN resumes r ON c.resume_id = r.id
		WHERE
			c.id=@id
			AND r.user_id=@user_id
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"id":      certificationID,
		"user_id": userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get certification by id query for certification_id=%s user_id=%s: %w", certificationID.String(), userID, err)
	}

	certificationItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[certification.Certification])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:certifications for certification_id=%s user_id=%s: %w", certificationID.String(), userID, err)
	}

	return &certificationItem, nil
}

func (r *CertificationRepository) GetCertificationsByResumeID(ctx context.Context, userID string, resumeID uuid.UUID) ([]certification.Certification, error) {
	stmt := `
		SELECT
			c.*
		FROM
			certifications c
		JOIN resumes r ON c.resume_id = r.id
		WHERE
			c.resume_id=@resume_id
			AND r.user_id=@user_id
		ORDER BY c.order_index ASC
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"resume_id": resumeID,
		"user_id":   userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get certifications by resume query for resume_id=%s user_id=%s: %w", resumeID.String(), userID, err)
	}

	certificationItems, err := pgx.CollectRows(rows, pgx.RowToStructByName[certification.Certification])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []certification.Certification{}, nil
		}
		return nil, fmt.Errorf("failed to collect rows from table:certifications for resume_id=%s user_id=%s: %w", resumeID.String(), userID, err)
	}

	return certificationItems, nil
}

func (r *CertificationRepository) UpdateCertification(ctx context.Context, userID string, certificationID uuid.UUID, payload *certification.UpdateCertificationRequest) (*certification.Certification, error) {
	stmt := `UPDATE certifications SET `
	args := pgx.NamedArgs{
		"id": certificationID,
	}
	setClauses := []string{}

	if payload.Name != nil {
		setClauses = append(setClauses, "name = @name")
		args["name"] = *payload.Name
	}
	if payload.Organization != nil {
		setClauses = append(setClauses, "organization = @organization")
		args["organization"] = *payload.Organization
	}
	if payload.IssueDate != nil {
		setClauses = append(setClauses, "issue_date = @issue_date")
		args["issue_date"] = *payload.IssueDate
	}
	if payload.ExpiryDate != nil {
		setClauses = append(setClauses, "expiry_date = @expiry_date")
		args["expiry_date"] = *payload.ExpiryDate
	}
	if payload.CredentialID != nil {
		setClauses = append(setClauses, "credential_id = @credential_id")
		args["credential_id"] = *payload.CredentialID
	}
	if payload.CredentialURL != nil {
		setClauses = append(setClauses, "credential_url = @credential_url")
		args["credential_url"] = *payload.CredentialURL
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
		return nil, fmt.Errorf("failed to execute update certification query for certification_id=%s user_id=%s: %w", certificationID.String(), userID, err)
	}

	certificationItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[certification.Certification])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:certifications for certification_id=%s user_id=%s: %w", certificationID.String(), userID, err)
	}

	return &certificationItem, nil
}

func (r *CertificationRepository) BulkUpdateCertificationOrder(ctx context.Context, userID string, payload *certification.BulkUpdateCertificationsRequest) error {
	tx, err := r.server.DB.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	for _, certificationUpdate := range payload.Certifications {
		_, err := tx.Exec(ctx, `
			UPDATE certifications 
			SET order_index = @order_index
			WHERE id = @id 
			AND resume_id IN (SELECT id FROM resumes WHERE user_id = @user_id)
		`, pgx.NamedArgs{
			"id":          certificationUpdate.ID,
			"order_index": certificationUpdate.OrderIndex,
			"user_id":     userID,
		})
		if err != nil {
			return fmt.Errorf("failed to update certification order for certification_id=%s: %w", certificationUpdate.ID, err)
		}
	}

	return tx.Commit(ctx)
}

func (r *CertificationRepository) DeleteCertification(ctx context.Context, userID string, certificationID uuid.UUID) error {
	result, err := r.server.DB.Pool.Exec(ctx, `
		DELETE FROM certifications
		WHERE id = @id 
		AND resume_id IN (SELECT id FROM resumes WHERE user_id = @user_id)
	`, pgx.NamedArgs{
		"id":      certificationID,
		"user_id": userID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete certification: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("certification not found")
	}

	return nil
}
