package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nguyenhuy260301/e-wallet/internal/auth"
	"github.com/nguyenhuy260301/e-wallet/internal/domain"
	"github.com/nguyenhuy260301/e-wallet/internal/repository"
	"github.com/nguyenhuy260301/e-wallet/internal/worker"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	accountRepo repository.AccountRepository
	txRepo      repository.TransactionRepository
	userRepo    repository.UserRepository
	pool        *worker.Pool
}

func NewHandler(
	accountRepo repository.AccountRepository,
	txRepo repository.TransactionRepository,
	userRepo repository.UserRepository,
	pool *worker.Pool,
) *Handler {
	return &Handler{
		accountRepo: accountRepo,
		txRepo:      txRepo,
		userRepo:    userRepo,
		pool:        pool,
	}
}

// @Summary     Lấy thông tin tài khoản đang đăng nhập
// @Tags        accounts
// @Produce     json
// @Success     200 {object} domain.Account
// @Failure     404 {object} map[string]string
// @Router      /accounts/me [get]
// @Security    BearerAuth
func (h *Handler) GetAccount(c *gin.Context) {
	id := c.MustGet("account_id").(int64)

	acc, err := h.accountRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "không tìm thấy tài khoản"})
		return
	}

	c.JSON(http.StatusOK, acc)
}

// @Summary     Nạp tiền vào tài khoản của chính mình (DEMO — mô phỏng nộp tiền mặt tại quầy, không kết nối cổng thanh toán thật)
// @Tags        accounts
// @Accept      json
// @Produce     json
// @Param       deposit body domain.DepositRequest true "Số tiền nạp"
// @Success     200 {object} map[string]string
// @Failure     400 {object} map[string]string
// @Router      /deposit [post]
// @Security    BearerAuth
func (h *Handler) Deposit(c *gin.Context) {
	var req domain.DepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "số tiền nạp phải lớn hơn 0"})
		return
	}

	accountID := c.MustGet("account_id").(int64)

	if err := h.accountRepo.Deposit(accountID, req.Amount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "nạp tiền thành công (demo)"})
}

// @Summary     Chuyển tiền giữa 2 tài khoản
// @Tags        transactions
// @Accept      json
// @Produce     json
// @Param       transfer body domain.TransferRequest true "Thông tin chuyển tiền"
// @Success     200 {object} domain.Transaction
// @Failure     400 {object} map[string]string
// @Router      /transfer [post]
// @Security    BearerAuth
func (h *Handler) Transfer(c *gin.Context) {
	var req domain.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "số tiền chuyển phải lớn hơn 0"})
		return
	}

	fromID := c.MustGet("account_id").(int64)

	toAccount, err := h.accountRepo.GetByNo(req.ToAccountNo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "không tìm thấy số tài khoản nhận"})
		return
	}

	trx, err := h.txRepo.Transfer(fromID, toAccount.ID, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trx)
}

// @Summary     Chuyển tiền hàng loạt (nhiều lệnh cùng lúc), xử lý song song có giới hạn qua Worker Pool
// @Tags        transactions
// @Accept      json
// @Produce     json
// @Param       transfers body domain.BulkTransferRequest true "Danh sách lệnh chuyển tiền"
// @Success     200 {object} domain.BulkTransferResponse
// @Failure     400 {object} map[string]string
// @Router      /transfer/bulk [post]
// @Security    BearerAuth
func (h *Handler) BulkTransfer(c *gin.Context) {
	var req domain.BulkTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fromID := c.MustGet("account_id").(int64)
	replyCh := make(chan worker.TransferResult, len(req.Transfers))

	for _, item := range req.Transfers {
		toAccount, err := h.accountRepo.GetByNo(item.ToAccountNo)
		if err != nil {
			replyCh <- worker.TransferResult{
				Job: worker.TransferJob{FromID: fromID, ToAccountNo: item.ToAccountNo, Amount: item.Amount},
				Err: errors.New("không tìm thấy số tài khoản nhận"),
			}
			continue
		}

		h.pool.Submit(worker.TransferJob{
			FromID:      fromID,
			ToID:        toAccount.ID,
			ToAccountNo: item.ToAccountNo,
			Amount:      item.Amount,
			Reply:       replyCh,
		})
	}

	resp := domain.BulkTransferResponse{
		Total:   len(req.Transfers),
		Results: []domain.BulkTransferItemResult{},
	}

	for i := 0; i < len(req.Transfers); i++ {
		result := <-replyCh

		itemResult := domain.BulkTransferItemResult{
			ToAccountNo: result.Job.ToAccountNo,
			Amount:      result.Job.Amount,
		}
		if result.Err != nil {
			itemResult.Status = "failed"
			itemResult.Error = result.Err.Error()
			resp.Failed++
		} else {
			itemResult.Status = "success"
			itemResult.TxCode = result.Trx.TxCode
			resp.Success++
		}
		resp.Results = append(resp.Results, itemResult)
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary     Lấy lịch sử giao dịch của tài khoản đang đăng nhập
// @Tags        transactions
// @Produce     json
// @Success     200 {array} domain.Transaction
// @Failure     500 {object} map[string]string
// @Router      /transactions [get]
// @Security    BearerAuth
func (h *Handler) GetTransactions(c *gin.Context) {
	accountID := c.MustGet("account_id").(int64)

	txs, err := h.txRepo.GetByAccountID(accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, txs)
}

// @Summary     Tìm 1 giao dịch theo mã giao dịch (chỉ xem được giao dịch của chính mình)
// @Tags        transactions
// @Produce     json
// @Param       code path string true "Mã giao dịch"
// @Success     200 {object} domain.Transaction
// @Failure     404 {object} map[string]string
// @Router      /transactions/{code} [get]
// @Security    BearerAuth
func (h *Handler) SearchTransaction(c *gin.Context) {
	code := c.Param("code")

	accountID := c.MustGet("account_id").(int64)

	trx, err := h.txRepo.GetByCode(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "không tìm thấy giao dịch"})
		return
	}

	if trx.FromAccID != accountID && trx.ToAccID != accountID {
		c.JSON(http.StatusNotFound, gin.H{"error": "không tìm thấy giao dịch"})
		return
	}

	c.JSON(http.StatusOK, trx)
}

// @Summary     Đăng ký tài khoản mới (tự tạo Account + User)
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       info body domain.RegisterRequest true "Username, password, họ tên"
// @Success     201 {object} map[string]interface{}
// @Failure     400 {object} map[string]string
// @Failure     409 {object} map[string]string
// @Router      /register [post]
func (h *Handler) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username không được để trống"})
		return
	}

	if len(req.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password phải có ít nhất 6 ký tự"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "không tạo được mật khẩu"})
		return
	}

	user, err := h.userRepo.Create(req.Username, string(hash), req.FullName)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "username đã tồn tại hoặc lỗi tạo tài khoản"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "đăng ký thành công",
		"username":   user.Username,
		"account_id": user.AccountID,
	})
}

// @Summary     Đăng nhập lấy JWT token
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       credentials body domain.LoginRequest true "Username và password"
// @Success     200 {object} map[string]string
// @Failure     401 {object} map[string]string
// @Router      /login [post]
func (h *Handler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userRepo.GetByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "sai username hoặc password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "sai username hoặc password"})
		return
	}

	token, err := auth.GenerateToken(user.Username, user.AccountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "không tạo được token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
