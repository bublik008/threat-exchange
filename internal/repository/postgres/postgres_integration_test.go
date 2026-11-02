//go:build integration

package postgres

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/bublik008/threat-exchange/internal/api/gen"
	repoerrors "github.com/bublik008/threat-exchange/internal/repository"
)

var testPool *pgxpool.Pool
var testDBURL string

func TestMain(m *testing.M) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()

	testDBURL, err = pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("failed to get connection string: %s", err)
	}

	testPool, err = pgxpool.New(ctx, testDBURL)
	if err != nil {
		log.Fatalf("failed to connect to db: %s", err)
	}
	defer testPool.Close()

	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(b), "..", "..", "..")
	migrationPath := filepath.Join(projectRoot, "migrations", "0001_init_sql.up.sql")
	migrationSQL, err := os.ReadFile(migrationPath)
	if err != nil {
		log.Fatalf("failed to read migration file: %s", err)
	}

	if _, err := testPool.Exec(ctx, string(migrationSQL)); err != nil {
		log.Fatalf("failed to apply migration: %s", err)
	}

	code := m.Run()
	os.Exit(code)
}

func TestRepository_CreateAndGet(t *testing.T) {
	repo := NewRepository(testPool)
	ctx := context.Background()

	threatType := "trojan"
	tags := []string{"test", "integration"}
	
	createReq := &gen.ThreatCreateRequest{
		Md5:        "d41d8cd98f00b204e9800998ecf8427e",
		Name:       "Integration Test Trojan",
		Severity:   "high",
		ThreatType: &threatType,
		Tags:       &tags,
	}

	created, err := repo.Create(ctx, createReq, time.Now())
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	
	// Безопасно проверяем указатель
	if created.Md5 == nil || *created.Md5 != createReq.Md5 {
		t.Fatalf("Created MD5 mismatch")
	}

	fetched, err := repo.GetByMD5(ctx, createReq.Md5)
	if err != nil {
		t.Fatalf("GetByMD5 failed: %v", err)
	}
	
	if fetched.Name == nil || *fetched.Name != (*createReq).Name {
		t.Fatalf("Fetched Name mismatch: got %v", fetched.Name)
	}
	if fetched.Severity == nil || string(*fetched.Severity) != string((*createReq).Severity) {
		t.Fatalf("Fetched Severity mismatch: got %v", fetched.Severity)
	}
	if fetched.Tags == nil || len(*fetched.Tags) != len(tags) {
		t.Fatalf("Fetched Tags mismatch: got %v", fetched.Tags)
	}
}

func TestRepository_CreateDuplicate(t *testing.T) {
	repo := NewRepository(testPool)
	ctx := context.Background()

	createReq := &gen.ThreatCreateRequest{
		Md5:      "e99a18c428cb38d5f260853678922e03",
		Name:     "Duplicate Test",
		Severity: "low",
	}

	_, err := repo.Create(ctx, createReq, time.Now())
	require.NoError(t, err)


	_, err = repo.Create(ctx, createReq, time.Now())
	require.Error(t, err)
	assert.ErrorIs(t, err, repoerrors.ErrConflict)
}

func TestRepository_GetNotFound(t *testing.T) {
	repo := NewRepository(testPool)
	ctx := context.Background()

	_, err := repo.GetByMD5(ctx, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	require.Error(t, err)
	assert.ErrorIs(t, err, repoerrors.ErrNotFound)
}

func stringPtr(s string) *string { return &s }