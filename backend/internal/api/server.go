package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nguyenhuy260301/e-wallet/internal/repository"
	"github.com/nguyenhuy260301/e-wallet/internal/worker"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

type Server struct {
	router  *gin.Engine
	handler *Handler
}

func NewServer(
	accountRepo repository.AccountRepository,
	txRepo repository.TransactionRepository,
	userRepo repository.UserRepository,
	pool *worker.Pool,
) *Server {
	h := NewHandler(accountRepo, txRepo, userRepo, pool)
	r := gin.Default()
	r.Use(corsMiddleware())

	// Public routes — không cần token
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Protected routes — phải có JWT token
	protected := r.Group("/")
	protected.Use(JWTMiddleware())
	{
		protected.GET("/accounts/me", h.GetAccount)
		protected.POST("/deposit", h.Deposit)
		protected.POST("/transfer", h.Transfer)
		protected.POST("/transfer/bulk", h.BulkTransfer)
		protected.GET("/transactions", h.GetTransactions)
		protected.GET("/transactions/:code", h.SearchTransaction)
	}

	return &Server{router: r, handler: h}
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
