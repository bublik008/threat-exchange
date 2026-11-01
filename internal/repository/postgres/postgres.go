// Package postgres реализует ThreatRepository с использованием библиотеки pgx и PostgreSQL.
package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/bublik008/threat-exchange/internal/api/gen" // ЗАМЕНИ ник!
	"github.com/bublik008/threat-exchange/internal/repository" // ЗАМЕНИ ник!
)

// Repository — реализация интерфейса repository.ThreatRepository.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository создает новый экземпляр репозитория.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Create добавляет новую угрозу в БД.
func (r *Repository) Create(ctx context.Context, req *gen.ThreatCreateRequest, now time.Time) (*gen.ThreatResponse, error) {
	query := `
		INSERT INTO threats (md5, name, severity, threat_type, description, tags, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING md5, name, severity, threat_type, description, tags, created_at, updated_at
	`

	// Делаем теги пустым массивом, если они не переданы, чтобы в БД не был NULL
	tags := []string{}
	if req.Tags != nil {
		tags = *req.Tags
	}

	var resp gen.ThreatResponse
	err := r.pool.QueryRow(ctx, query,
		req.Md5, req.Name, req.Severity,
		req.ThreatType, req.Description, tags, now, now,
	).Scan(
		&resp.Md5, &resp.Name, &resp.Severity,
		&resp.ThreatType, &resp.Description, &resp.Tags,
		&resp.CreatedAt, &resp.UpdatedAt,
	)

	if err != nil {
		// Код 23505 в Postgres — это нарушение уникальности (unique_violation).
		// Если он возник, значит MD5 уже существует.
		if pgErr, ok := err.(interface{ SQLState() string }); ok && pgErr.SQLState() == "23505" {
			return nil, repository.ErrConflict
		}
		return nil, fmt.Errorf("failed to insert threat: %w", err)
	}

	return &resp, nil
}

// GetByMD5 получает угрозу по MD5.
func (r *Repository) GetByMD5(ctx context.Context, md5 string) (*gen.ThreatResponse, error) {
	query := `
		SELECT md5, name, severity, threat_type, description, tags, created_at, updated_at
		FROM threats
		WHERE md5 = $1
	`

	var resp gen.ThreatResponse
	err := r.pool.QueryRow(ctx, query, md5).Scan(
		&resp.Md5, &resp.Name, &resp.Severity,
		&resp.ThreatType, &resp.Description, &resp.Tags,
		&resp.CreatedAt, &resp.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get threat: %w", err)
	}

	return &resp, nil
}

// Update обновляет поля угрозы по MD5.
func (r *Repository) Update(ctx context.Context, md5 string, req *gen.ThreatUpdateRequest) (*gen.ThreatResponse, error) {
	query := `
		UPDATE threats
		SET name = $1, severity = COALESCE($2, severity), threat_type = COALESCE($3, threat_type),
		    description = COALESCE($4, description), tags = COALESCE($5, tags), updated_at = now()
		WHERE md5 = $6
		RETURNING md5, name, severity, threat_type, description, tags, created_at, updated_at
	`

	var resp gen.ThreatResponse
	err := r.pool.QueryRow(ctx, query,
		req.Name, req.Severity, req.ThreatType, req.Description, req.Tags, md5,
	).Scan(
		&resp.Md5, &resp.Name, &resp.Severity,
		&resp.ThreatType, &resp.Description, &resp.Tags,
		&resp.CreatedAt, &resp.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("failed to update threat: %w", err)
	}

	return &resp, nil
}