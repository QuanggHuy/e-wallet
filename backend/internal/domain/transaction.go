package domain

import "time"

type Transaction struct {
	ID        int64     `db:"id" json:"id"`
	TxCode    string    `db:"tx_code" json:"tx_code"`
	FromAccID int64     `db:"from_acc_id" json:"from_acc_id"`
	ToAccID   int64     `db:"to_acc_id" json:"to_acc_id"`
	Amount    float64   `db:"amount" json:"amount"`
	Status    string    `db:"status" json:"status"`
	Note      string    `db:"note" json:"note,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
