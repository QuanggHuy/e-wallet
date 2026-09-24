package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nguyenhuy260301/e-wallet/internal/domain"
	"github.com/nguyenhuy260301/e-wallet/internal/repository"
)

func TestGetAccount_Found(t *testing.T) {
	mockRepo := &repository.MockAccountRepository{
		GetByIDFn: func(id int64) (*domain.Account, error) {
			return &domain.Account{
				ID:        id,
				AccountNo: "ACC001",
				FullName:  "Nguyen Van A",
				Balance:   1000000,
			}, nil
		},
	}

	h := NewHandler(mockRepo, nil, nil, nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/accounts/me", func(c *gin.Context) {
		c.Set("account_id", int64(1))
		h.GetAccount(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/accounts/me", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var acc domain.Account
	if err := json.Unmarshal(w.Body.Bytes(), &acc); err != nil {
		t.Fatal("không parse được response body")
	}

	if acc.AccountNo != "ACC001" {
		t.Errorf("expected ACC001, got %s", acc.AccountNo)
	}
}

func TestGetAccount_NotFound(t *testing.T) {
	mockRepo := &repository.MockAccountRepository{
		GetByIDFn: func(id int64) (*domain.Account, error) {
			return nil, errors.New("not found")
		},
	}

	h := NewHandler(mockRepo, nil, nil, nil)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/accounts/me", func(c *gin.Context) {
		c.Set("account_id", int64(99))
		h.GetAccount(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/accounts/me", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}
