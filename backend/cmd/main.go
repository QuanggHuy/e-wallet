package main

import (
	"log"

	"github.com/joho/godotenv"
	_ "github.com/nguyenhuy260301/e-wallet/docs"
	"github.com/nguyenhuy260301/e-wallet/internal/api"
	"github.com/nguyenhuy260301/e-wallet/internal/db"
	"github.com/nguyenhuy260301/e-wallet/internal/repository"
	"github.com/nguyenhuy260301/e-wallet/internal/worker"
)

// @title          GoWallet API
// @version        1.0
// @description    Digital Wallet Service — chuyển tiền, quản lý tài khoản
// @host           localhost:8080
// @BasePath       /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Không load được .env:", err)
	}

	pool, err := db.NewPool()
	if err != nil {
		log.Fatal("Không tạo được pool:", err)
	}
	defer pool.Close()

	accountRepo := repository.NewAccountRepository(pool)
	txRepo := repository.NewTransactionRepository(pool)
	userRepo := repository.NewUserRepository(pool)

	workerPool := worker.NewPool(txRepo, 3)

	server := api.NewServer(accountRepo, txRepo, userRepo, workerPool)
	log.Println("Server chạy tại http://localhost:8080")
	if err := server.Run(":8080"); err != nil {
		log.Fatal("Server lỗi:", err)
	}
}
