package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nguyenhuy260301/e-wallet/internal/domain"
)

type AccountRepository interface {
	GetByID(id int64) (*domain.Account, error)
	GetByNo(accountNo string) (*domain.Account, error)
	Deposit(id int64, amount float64) error
}

type accountRepository struct {
	pool *pgxpool.Pool
}

func NewAccountRepository(pool *pgxpool.Pool) AccountRepository {
	return &accountRepository{pool: pool}
}

func (r *accountRepository) GetByID(id int64) (*domain.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, account_no, full_name, balance, created_at, updated_at
              FROM accounts
              WHERE id = $1`

	acc := &domain.Account{}
	err := r.pool.QueryRow(ctx, query, id).
		Scan(&acc.ID, &acc.AccountNo, &acc.FullName, &acc.Balance, &acc.CreatedAt, &acc.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (r *accountRepository) Deposit(id int64, amount float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE accounts SET balance = balance + $1, updated_at = NOW() WHERE id = $2`

	_, err := r.pool.Exec(ctx, query, amount, id)
	return err
}

func (r *accountRepository) GetByNo(accountNo string) (*domain.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, account_no, full_name, balance, created_at, updated_at
              FROM accounts
              WHERE account_no = $1`

	acc := &domain.Account{}
	err := r.pool.QueryRow(ctx, query, accountNo).
		Scan(&acc.ID, &acc.AccountNo, &acc.FullName, &acc.Balance, &acc.CreatedAt, &acc.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return acc, nil
}
