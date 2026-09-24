package worker

import (
	"log"

	"github.com/nguyenhuy260301/e-wallet/internal/domain"
	"github.com/nguyenhuy260301/e-wallet/internal/repository"
)

type TransferJob struct {
	FromID      int64
	ToID        int64
	ToAccountNo string // mang theo để hiển thị kết quả — vì kết quả trả về không đảm bảo đúng thứ tự đã submit
	Amount      float64
	Reply       chan TransferResult // kênh riêng để nhận đúng kết quả của job này — tránh lẫn kết quả giữa các request khác nhau chạy cùng lúc
}

type TransferResult struct {
	Job TransferJob
	Trx *domain.Transaction
	Err error
}

type Pool struct {
	jobs chan TransferJob
	repo repository.TransactionRepository
}

func NewPool(repo repository.TransactionRepository, numWorkers int) *Pool {
	p := &Pool{
		jobs: make(chan TransferJob, 100),
		repo: repo,
	}
	p.Start(numWorkers)
	return p
}

func (p *Pool) Start(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		go p.worker()
	}
}

func (p *Pool) worker() {
	for job := range p.jobs {
		trx, err := p.repo.Transfer(job.FromID, job.ToID, job.Amount)
		if err != nil {
			log.Printf("Transfer failed %d->%d: %v", job.FromID, job.ToID, err)
		}
		job.Reply <- TransferResult{Job: job, Trx: trx, Err: err}
	}
}

func (p *Pool) Submit(job TransferJob) {
	p.jobs <- job
}
