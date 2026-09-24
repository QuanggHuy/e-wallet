package repository

import "github.com/nguyenhuy260301/e-wallet/internal/domain"

type MockAccountRepository struct {
	GetByIDFn func(id int64) (*domain.Account, error)
	GetByNoFn func(accountNo string) (*domain.Account, error)
	DepositFn func(id int64, amount float64) error
}

func (m *MockAccountRepository) GetByID(id int64) (*domain.Account, error) {
	return m.GetByIDFn(id)
}

func (m *MockAccountRepository) GetByNo(accountNo string) (*domain.Account, error) {
	return m.GetByNoFn(accountNo)
}

func (m *MockAccountRepository) Deposit(id int64, amount float64) error {
	return m.DepositFn(id, amount)
}
