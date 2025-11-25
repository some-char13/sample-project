package repository_test

import (
	"context"
	"log"
	"testing"
	"time"

	"sample_project/internal/model/check"
	"sample_project/internal/model/service"
	"sample_project/internal/repository"

	"sample_project/internal/migrations"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func testDB(t *testing.T) (context.Context, *repository.PostgresRepository, func()) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:14",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		//postgres.WithInitScripts("../../migrations/001_create_tables.sql"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	require.NoError(t, err)

	dsn, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	require.NoError(t, migrations.Migrate(dsn))

	repo, err := repository.NewPostgresRepository(dsn)
	require.NoError(t, err)

	cleanup := func() {
		repo.Close()
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Printf("failed to terminate container: %v", err)
		}
	}

	return ctx, repo, cleanup
}
func TestPostgresRepository_FullLifecycle(t *testing.T) {
	ctx, repo, cleanup := testDB(t)
	defer cleanup()

	svcReq := &service.ServiceRequest{
		Name:     "корпорация зла",
		Url:      "https://google.com",
		Interval: 30,
	}
	require.NoError(t, svcReq.Validate())

	created, err := repo.AddService(ctx, service.NewService(svcReq.Name, svcReq.Url, svcReq.Interval))
	require.NoError(t, err)
	require.NotZero(t, created.Id)
	assert.Equal(t, "корпорация зла", created.Name)

	found, err := repo.GetServiceByID(ctx, created.Id)
	require.NoError(t, err)
	assert.Equal(t, created.Id, found.Id)
	assert.Equal(t, "https://google.com", found.Url)

	result1 := check.NewResult(created.Id, 200, 145)
	result2 := check.NewResult(created.Id, 503, 3200)

	require.NoError(t, repo.AddCheckResult(ctx, result1))
	require.NoError(t, repo.AddCheckResult(ctx, result2))

	results, err := repo.GetCheckResults(ctx, created.Id, 20)
	require.NoError(t, err)
	require.Len(t, results, 2)
	assert.Equal(t, 503, results[0].ResponseCode)
	assert.Equal(t, 200, results[1].ResponseCode)

	latest, err := repo.GetLastServiceCheck(ctx, created.Id)
	require.NoError(t, err)
	require.NotNil(t, latest)
	assert.Equal(t, 503, latest.ResponseCode)

	failed, err := repo.GetCheckResultsByStatus(ctx, created.Id, []int{503}, 10)
	require.NoError(t, err)
	require.Len(t, failed, 1)
	assert.Equal(t, 503, failed[0].ResponseCode)

	success, err := repo.GetCheckResultsByStatus(ctx, created.Id, []int{200}, 10)
	require.NoError(t, err)
	require.Len(t, success, 1)
	assert.Equal(t, 200, success[0].ResponseCode)

	require.NoError(t, repo.DeleteService(ctx, created.Id))

	_, err = repo.GetServiceByID(ctx, created.Id)
	assert.Error(t, err)

	resultsAfterDelete, err := repo.GetCheckResults(ctx, created.Id, 10)
	require.NoError(t, err)
	assert.Empty(t, resultsAfterDelete)
}

func TestPostgresRepository_GetAllServicesStatus_Integration(t *testing.T) {
	ctx, repo, cleanup := testDB(t)
	defer cleanup()

	s1 := service.NewService("Yandex", "https://ya.ru", 60)
	s2 := service.NewService("GitHub", "https://github.com", 30)

	var err error
	s1, err = repo.AddService(ctx, s1)
	require.NoError(t, err)

	s2, err = repo.AddService(ctx, s2)
	require.NoError(t, err)

	require.NoError(t, repo.AddCheckResult(ctx, check.NewResult(s1.Id, 200, 120)))
	require.NoError(t, repo.AddCheckResult(ctx, check.NewResult(s2.Id, 500, 5000)))

	services, err := repo.GetServices(ctx)
	require.NoError(t, err)
	require.Len(t, services, 2)

	latest1, err := repo.GetLastServiceCheck(ctx, s1.Id)
	require.NoError(t, err)
	assert.Equal(t, 200, latest1.ResponseCode)
}
