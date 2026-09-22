package repository

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"

	"orderApi/internal/domain"
)

type PostgresOrderRepository struct {
	db *sqlx.DB
}

type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	GetByID(ctx context.Context, id string) (*domain.Order, error)
	List(ctx context.Context) ([]domain.Order, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	Delete(ctx context.Context, id string) error
}

func NewPostgresOrderRepository(db *sqlx.DB) *PostgresOrderRepository {
	return &PostgresOrderRepository{
		db: db,
	}
}

func (r *PostgresOrderRepository) Create(ctx context.Context, order *domain.Order) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO orders (id, user_id, status)
		VALUES ($1, $2, $3)`,
		order.ID,
		order.UserID,
		order.Status,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresOrderRepository) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	var order domain.Order

	err := r.db.GetContext(
		ctx,
		&order,
		`SELECT id, user_id, status
		FROM orders
		WHERE id = $1`,
		id,
	)
	if err != nil {
		slog.Error("failed to get order by id", "error", err)
		return nil, err
	}
	return &order, nil
}

func (r *PostgresOrderRepository) List(ctx context.Context) ([]domain.Order, error) {
	orders := []domain.Order{}

	err := r.db.SelectContext(
		ctx,
		&orders,
		"SELECT id, user_id, status FROM orders",
	)

	if err != nil {
		return nil, fmt.Errorf("list order: %w", err)
	}
	return orders, nil
}
