package domain

import "time"

type User struct {
	ID           int64     `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	PasswordHash string    `db:"password_hash" json:"-"`
	AccountID    int64     `db:"account_id" json:"account_id"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}
