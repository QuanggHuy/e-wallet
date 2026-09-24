package domain

import "time"

type Account struct {
	ID        int64     `db:"id" json:"id"`
	AccountNo string    `db:"account_no" json:"account_no"`
	FullName  string    `db:"full_name" json:"full_name"`
	Balance   float64   `db:"balance" json:"balance"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TransferRequest struct {
	ToAccountNo string  `json:"to_account_no"`
	Amount      float64 `json:"amount"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

type DepositRequest struct {
	Amount float64 `json:"amount"`
}

type BulkTransferRequest struct {
	Transfers []TransferRequest `json:"transfers"`
}

type BulkTransferItemResult struct {
	ToAccountNo string  `json:"to_account_no"`
	Amount      float64 `json:"amount"`
	Status      string  `json:"status"`
	TxCode      string  `json:"tx_code,omitempty"`
	Error       string  `json:"error,omitempty"`
}

type BulkTransferResponse struct {
	Total   int                       `json:"total"`
	Success int                       `json:"success"`
	Failed  int                       `json:"failed"`
	Results []BulkTransferItemResult  `json:"results"`
}
