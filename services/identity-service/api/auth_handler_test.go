package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/k1ngalph0x/atlas/services/identity-service/api"
	"github.com/k1ngalph0x/atlas/services/identity-service/config"
	"github.com/k1ngalph0x/atlas/services/identity-service/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)


func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

func setupRouter(db *gorm.DB) (*api.AuthHandler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		TOKEN: config.TokenConfig{JwtKey: "test-secret-key"},
	}
	h := api.NewAuthHandler(db, cfg)
	r := gin.New()
	r.POST("/auth/signup", h.SignUp)
	r.POST("/auth/signin", h.SignIn)
	return h, r
}

func post(t *testing.T, r *gin.Engine, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func decodeBody(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var out map[string]any
	if err := json.NewDecoder(w.Body).Decode(&out); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	return out
}


func TestSignUp_Success(t *testing.T) {
	_, r := setupRouter(setupTestDB(t))

	w := post(t, r, "/auth/signup", map[string]any{
		"email":    "john@example.com",
		"password": "password123",
	})

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d — body: %s", w.Code, w.Body.String())
	}

	body := decodeBody(t, w)
	if body["token"] == nil || body["token"] == "" {
		t.Error("expected token in response")
	}
	user, ok := body["user"].(map[string]any)
	if !ok || user["email"] != "john@example.com" {
		t.Error("expected user.email in response")
	}
}

func TestSignUp_DuplicateEmail(t *testing.T) {
	_, r := setupRouter(setupTestDB(t))

	payload := map[string]any{"email": "dup@example.com", "password": "password123"}
	post(t, r, "/auth/signup", payload)

	w := post(t, r, "/auth/signup", payload) 
	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", w.Code)
	}
}

func TestSignUp_InvalidEmail(t *testing.T) {
	_, r := setupRouter(setupTestDB(t))

	w := post(t, r, "/auth/signup", map[string]any{
		"email":    "not-an-email",
		"password": "password123",
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestSignUp_ShortPassword(t *testing.T) {
	_, r := setupRouter(setupTestDB(t))

	w := post(t, r, "/auth/signup", map[string]any{
		"email":    "short@example.com",
		"password": "abc",
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestSignUp_EmailNormalized(t *testing.T) {
	db := setupTestDB(t)
	_, r := setupRouter(db)

	w := post(t, r, "/auth/signup", map[string]any{
		"email":    "UPPER@Example.COM",
		"password": "password123",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d — body: %s", w.Code, w.Body.String())
	}

	var u models.User
	db.Where("email = ?", "upper@example.com").First(&u)
	if u.Email != "upper@example.com" {
		t.Errorf("expected normalized email, got %q", u.Email)
	}
}


func TestSignIn_Success(t *testing.T) {
	_, r := setupRouter(setupTestDB(t))

	post(t, r, "/auth/signup", map[string]any{
		"email":    "jane@example.com",
		"password": "password123",
	})

	w := post(t, r, "/auth/signin", map[string]any{
		"email":    "jane@example.com",
		"password": "password123",
	})

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}

	body := decodeBody(t, w)
	if body["token"] == nil || body["token"] == "" {
		t.Error("expected token in signin response")
	}
}

func TestSignIn_WrongPassword(t *testing.T) {
	_, r := setupRouter(setupTestDB(t))

	post(t, r, "/auth/signup", map[string]any{
		"email":    "jane@example.com",
		"password": "correctpassword",
	})

	w := post(t, r, "/auth/signin", map[string]any{
		"email":    "jane@example.com",
		"password": "wrongpassword",
	})

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestSignIn_NonExistentUser(t *testing.T) {
	_, r := setupRouter(setupTestDB(t))

	w := post(t, r, "/auth/signin", map[string]any{
		"email":    "ghost@example.com",
		"password": "password123",
	})

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestSignIn_MissingFields(t *testing.T) {
	_, r := setupRouter(setupTestDB(t))

	w := post(t, r, "/auth/signin", map[string]any{"email": "jane@example.com"})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}


func TestGenerateJWT_ValidToken(t *testing.T) {
	cfg := &config.Config{TOKEN: config.TokenConfig{JwtKey: "test-secret"}}
	h := api.NewAuthHandler(nil, cfg)

	token, err := h.GenerateJWT("user-123", "test@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token == "" {
		t.Error("expected non-empty token")
	}
}