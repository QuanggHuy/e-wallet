package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nguyenhuy260301/e-wallet/internal/domain"
)

type TransactionRepository interface {
	Transfer(fromID, toID int64, amount float64) (*domain.Transaction, error)
	GetByAccountID(accountID int64) ([]*domain.Transaction, error)
	GetByCode(code string) (*domain.Transaction, error)
}

type transactionRepository struct {
	pool *pgxpool.Pool
}

func (r *transactionRepository) GetByAccountID(accountID int64) ([]*domain.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, tx_code, from_account_id, to_account_id, amount, status, COALESCE(note, ''), created_at
              FROM transactions
              WHERE from_account_id = $1 OR to_account_id = $1
              ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []*domain.Transaction{}
	for rows.Next() {
		trx := &domain.Transaction{}
		err := rows.Scan(&trx.ID, &trx.TxCode, &trx.FromAccID, &trx.ToAccID, &trx.Amount, &trx.Status, &trx.Note, &trx.CreatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, trx)
	}

	return result, nil
}

func (r *transactionRepository) GetByCode(code string) (*domain.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, tx_code, from_account_id, to_account_id, amount, status, COALESCE(note, ''), created_at
              FROM transactions
              WHERE tx_code ILIKE $1`

	trx := &domain.Transaction{}
	err := r.pool.QueryRow(ctx, query, code).
		Scan(&trx.ID, &trx.TxCode, &trx.FromAccID, &trx.ToAccID, &trx.Amount, &trx.Status, &trx.Note, &trx.CreatedAt)
	if err != nil {
		return nil, err
	}

	return trx, nil
}

func NewTransactionRepository(pool *pgxpool.Pool) TransactionRepository {
	return &transactionRepository{pool: pool}
}

func (r *transactionRepository) Transfer(fromID, toID int64, amount float64) (*domain.Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("số tiền chuyển phải lớn hơn 0")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var fromBalance float64
	err = tx.QueryRow(ctx, "SELECT balance FROM accounts WHERE id = $1 FOR UPDATE", fromID).Scan(&fromBalance)
	if err != nil {
		return nil, err
	}

	if fromBalance < amount {
		// dùng r.pool.Exec (không phải tx.Exec) — vì tx sẽ bị rollback do return sớm, không commit được
		failCode := fmt.Sprintf("TX%d", time.Now().UnixNano())
		r.pool.Exec(context.Background(),
			`INSERT INTO transactions (tx_code, from_account_id, to_account_id, amount, status, note)
              VALUES ($1, $2, $3, $4, 'failed', 'số dư không đủ')`,
			failCode, fromID, toID, amount,
		)
		return nil, errors.New("số dư không đủ")
	}

	_, err = tx.Exec(ctx,
		`UPDATE accounts SET balance = balance - $1, updated_at = NOW() WHERE id = $2`,
		amount, fromID,
	)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(ctx,
		`UPDATE accounts SET balance = balance + $1, updated_at = NOW() WHERE id = $2`,
		amount, toID,
	)
	if err != nil {
		return nil, err
	}

	txCode := fmt.Sprintf("TX%d", time.Now().UnixNano())

	trx := &domain.Transaction{}
	err = tx.QueryRow(ctx,
		`INSERT INTO transactions (tx_code, from_account_id, to_account_id, amount, status)
              VALUES ($1, $2, $3, $4, 'success')
              RETURNING id, tx_code, from_account_id, to_account_id, amount, status, created_at`,
		txCode, fromID, toID, amount,
	).Scan(&trx.ID, &trx.TxCode, &trx.FromAccID, &trx.ToAccID, &trx.Amount, &trx.Status, &trx.CreatedAt)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return trx, nil
}
