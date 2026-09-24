package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nguyenhuy260301/e-wallet/internal/domain"
)

type UserRepository interface {
	GetByUsername(username string) (*domain.User, error)
	Create(username, passwordHash, fullName string) (*domain.User, error)
}

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) GetByUsername(username string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, username, password_hash, account_id, created_at
              FROM users
              WHERE username = $1`

	u := &domain.User{}
	err := r.pool.QueryRow(ctx, query, username).
		Scan(&u.ID, &u.Username, &u.PasswordHash, &u.AccountID, &u.CreatedAt)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *userRepository) Create(username, passwordHash, fullName string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	accountNo := fmt.Sprintf("ACC%d", time.Now().UnixNano())

	var accountID int64
	err = tx.QueryRow(ctx,
		`INSERT INTO accounts (account_no, full_name, balance) VALUES ($1, $2, 0) RETURNING id`,
		accountNo, fullName,
	).Scan(&accountID)
	if err != nil {
		return nil, err
	}

	u := &domain.User{Username: username, PasswordHash: passwordHash, AccountID: accountID}
	err = tx.QueryRow(ctx,
		`INSERT INTO users (username, password_hash, account_id) VALUES ($1, $2, $3) RETURNING id, created_at`,
		u.Username, u.PasswordHash, u.AccountID,
	).Scan(&u.ID, &u.CreatedAt)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return u, nil
}
