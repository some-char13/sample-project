package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	// needed for db connect
	_ "github.com/jackc/pgx/v4/stdlib"
	"sample_project/internal/model/check"
	"sample_project/internal/model/service"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(dsn string) (*PostgresRepository, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) Close() error {
	return r.db.Close()
}

func (r *PostgresRepository) Ping() error {
	return r.db.Ping()
}

func (r *PostgresRepository) AddService(ctx context.Context, s *service.Service) (*service.Service, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, fmt.Errorf("cannot begin transaction: %w", err)
	}

	defer tx.Rollback()

	query := `INSERT INTO services (name, url, interval_seconds) VALUES ($1, $2, $3) RETURNING id, created_at`
	if err := tx.QueryRowContext(ctx, query, s.Name, s.URL, s.Interval).Scan(&s.ID, &s.Created); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("cannot commit transaction: %w", err)
	}

	return s, nil
}

func (r *PostgresRepository) GetServices(ctx context.Context) ([]*service.Service, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, fmt.Errorf("cannot begin transaction: %w", err)
	}

	defer tx.Rollback()

	query := `SELECT id, name, url, interval_seconds, created_at FROM services ORDER BY id`
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var services []*service.Service
	for rows.Next() {
		var s service.Service
		err := rows.Scan(&s.ID, &s.Name, &s.URL, &s.Interval, &s.Created)
		if err != nil {
			return nil, err
		}
		services = append(services, &s)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("cannot commit transaction: %w", err)
	}

	return services, nil
}

func (r *PostgresRepository) GetServiceByID(ctx context.Context, id int) (*service.Service, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, fmt.Errorf("cannot begin transaction: %w", err)
	}

	defer tx.Rollback()

	query := `SELECT id, name, url, interval_seconds, created_at FROM services WHERE id = $1`
	var s service.Service
	err = tx.QueryRowContext(ctx, query, id).Scan(&s.ID, &s.Name, &s.URL, &s.Interval, &s.Created)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrServiceNotFound
	}
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("cannot commit transaction: %w", err)
	}

	return &s, nil
}

func (r *PostgresRepository) GetServiceByName(ctx context.Context, name string) (*service.Service, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, fmt.Errorf("cannot begin transaction: %w", err)
	}

	defer tx.Rollback()

	query := `SELECT id, name, url, interval_seconds, created_at FROM services WHERE name = $1`
	var s service.Service
	err = tx.QueryRowContext(ctx, query, name).Scan(&s.ID, &s.Name, &s.URL, &s.Interval, &s.Created)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrServiceNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("query service by name: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("cannot commit transaction: %w", err)
	}

	return &s, nil
}

func (r *PostgresRepository) DeleteService(ctx context.Context, id int) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return fmt.Errorf("cannot begin transaction: %w", err)
	}

	defer tx.Rollback()

	var exists int
	chk := `SELECT 1 FROM services WHERE id = $1`
	err = tx.QueryRowContext(ctx, chk, id).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("service with ID %d not found", id)
		}
		return fmt.Errorf("cannot check service existence: %w", err)
	}

	query := `DELETE FROM services WHERE id = $1`
	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete service with id: %d", id)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("cannot commit transaction: %w", err)
	}

	return err
}

func (r *PostgresRepository) AddCheckResult(ctx context.Context, result *check.Result) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return fmt.Errorf("cannot begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO check_results (service_id, status_code, response_time_ms)
		VALUES ($1, $2, $3)
		RETURNING id, checked_at
	`

	err = tx.QueryRowContext(ctx, query, result.ServiceID, result.ResponseCode, result.RespDuration).
		Scan(&result.ID, &result.TimeChecked)
	if err != nil {
		return fmt.Errorf("failed to insert check result: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("cannot commit transaction: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetCheckResults(ctx context.Context, serviceID int, limit int) ([]*check.Result, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, fmt.Errorf("cannot begin transaction: %w", err)
	}

	defer tx.Rollback()

	query := `
		SELECT id, service_id, status_code, response_time_ms, checked_at 
        FROM check_results 
		WHERE service_id = $1 
		ORDER BY checked_at DESC LIMIT $2
	`

	rows, err := tx.QueryContext(ctx, query, serviceID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var results []*check.Result
	for rows.Next() {
		var result check.Result
		err := rows.Scan(&result.ID, &result.ServiceID, &result.ResponseCode,
			&result.RespDuration, &result.TimeChecked)
		if err != nil {
			return nil, err
		}
		results = append(results, &result)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("cannot commit transaction: %w", err)
	}

	return results, nil
}

func (r *PostgresRepository) GetLastServiceCheck(
	ctx context.Context,
	serviceID int,
) (*check.ResultRequest, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, fmt.Errorf("cannot begin transaction: %w", err)
	}

	defer tx.Rollback()

	query := `
		SELECT service_id, status_code, response_time_ms, checked_at
        FROM check_results 
		WHERE service_id = $1 
		ORDER BY checked_at DESC LIMIT 1
	`

	var result check.ResultRequest

	err = tx.QueryRowContext(ctx, query, serviceID).Scan(&result.ServiceID, &result.ResponseCode,
		&result.RespDuration, &result.TimeChecked)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrCheckNotFound
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("cannot commit transaction: %w", err)
	}

	return &result, nil
}

func (r *PostgresRepository) GetCheckResultsByStatus(
	ctx context.Context,
	serviceID int,
	respCodes []int,
	limit int,
) ([]*check.Result, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, fmt.Errorf("cannot begin transaction: %w", err)
	}

	defer tx.Rollback()

	query := `
		SELECT id, service_id, status_code, response_time_ms, checked_at 
        FROM check_results 
		WHERE service_id = $1 AND status_code = ANY($2) 
		ORDER BY checked_at DESC LIMIT $3
	`

	rows, err := tx.QueryContext(ctx, query, serviceID, respCodes, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var results []*check.Result
	for rows.Next() {
		var result check.Result
		err := rows.Scan(&result.ID, &result.ServiceID, &result.ResponseCode,
			&result.RespDuration, &result.TimeChecked)
		if err != nil {
			return nil, err
		}
		results = append(results, &result)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("cannot commit transaction: %w", err)
	}

	return results, nil
}
