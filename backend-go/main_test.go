package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
)

// setupRouter configures a test router with security/CORS middleware and the target routes
func setupRouter() *gin.Engine {
	// Set gin to test mode to avoid verbose logging in tests
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(securityHeadersMiddleware())
	router.Use(corsMiddleware())

	api := router.Group("/api/v1")
	{
		api.GET("/health/metabolic", fetchMetabolicMetrics)
		api.POST("/health/metabolic", updateMetabolicMetrics)
		api.POST("/checkout", handleCheckout)
		api.POST("/calendar", handleAddCalendarSlot)
		api.GET("/calendar", fetchCalendarSlots)
	}

	return router
}

// TestFetchMetabolicMetrics verifies GET /api/v1/health/metabolic
func TestFetchMetabolicMetrics(t *testing.T) {
	t.Run("Fallback when db is nil", func(t *testing.T) {
		db = nil
		router := setupRouter()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/health/metabolic", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var resp MetabolicLog
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if resp.ConsumedCalories != 1420.50 || resp.ActiveBurn != 580.20 || resp.TargetAllowance != 2200.00 {
			t.Errorf("unexpected fallback response values: %+v", resp)
		}
	})

	t.Run("Success with DB metrics", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("error creating mock database: %v", err)
		}
		defer mockDB.Close()

		db = mockDB
		defer func() { db = nil }()

		rows := sqlmock.NewRows([]string{"consumed", "active", "target"}).
			AddRow(1850.25, 620.10, 2300.00)

		mock.ExpectQuery("^SELECT consumed, active, target FROM metabolic_logs ORDER BY id DESC LIMIT 1").
			WillReturnRows(rows)

		router := setupRouter()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/health/metabolic", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var resp MetabolicLog
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if resp.ConsumedCalories != 1850.25 || resp.ActiveBurn != 620.10 || resp.TargetAllowance != 2300.00 {
			t.Errorf("unexpected query response values: %+v", resp)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet mock expectations: %v", err)
		}
	})
}

// TestUpdateMetabolicMetrics verifies POST /api/v1/health/metabolic
func TestUpdateMetabolicMetrics(t *testing.T) {
	t.Run("Invalid payload", func(t *testing.T) {
		router := setupRouter()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/health/metabolic", bytes.NewBufferString("{invalid-json"))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Fallback when db is nil", func(t *testing.T) {
		db = nil
		router := setupRouter()
		w := httptest.NewRecorder()

		input := MetabolicLog{
			ConsumedCalories: 2000.0,
			ActiveBurn:       800.0,
			TargetAllowance:  2500.0,
		}
		body, _ := json.Marshal(input)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/health/metabolic", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var resp map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if resp["status"] != "success" {
			t.Errorf("expected status success, got %v", resp["status"])
		}
	})

	t.Run("Success insert to DB", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("error creating mock database: %v", err)
		}
		defer mockDB.Close()

		db = mockDB
		defer func() { db = nil }()

		mock.ExpectExec("^INSERT INTO metabolic_logs").
			WithArgs(1900.50, 750.30, 2400.00).
			WillReturnResult(sqlmock.NewResult(1, 1))

		router := setupRouter()
		w := httptest.NewRecorder()

		input := MetabolicLog{
			ConsumedCalories: 1900.50,
			ActiveBurn:       750.30,
			TargetAllowance:  2400.00,
		}
		body, _ := json.Marshal(input)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/health/metabolic", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet mock expectations: %v", err)
		}
	})
}

// TestHandleCheckout verifies POST /api/v1/checkout
func TestHandleCheckout(t *testing.T) {
	t.Run("Invalid payload", func(t *testing.T) {
		router := setupRouter()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/checkout", bytes.NewBufferString("{invalid-json"))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Fallback when db is nil", func(t *testing.T) {
		db = nil
		router := setupRouter()
		w := httptest.NewRecorder()

		input := CheckoutTransaction{
			TotalAmount:  120.50,
			ItemCount:    4,
			KMSSignature: "sig-abc-123",
		}
		body, _ := json.Marshal(input)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/checkout", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var resp map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if resp["status"] != "success" {
			t.Errorf("expected status success, got %v", resp["status"])
		}

		data, ok := resp["data"].(map[string]interface{})
		if !ok {
			t.Fatalf("data field missing or invalid in response: %+v", resp)
		}

		if data["kms_signature"] != "sig-abc-123" {
			t.Errorf("expected kms_signature sig-abc-123, got %v", data["kms_signature"])
		}

		if data["id"] == nil || data["id"] == 0 {
			t.Errorf("expected fallback ID to be generated, got: %v", data["id"])
		}
	})

	t.Run("Success insert to DB", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("error creating mock database: %v", err)
		}
		defer mockDB.Close()

		db = mockDB
		defer func() { db = nil }()

		fixedTime := time.Date(2026, 6, 13, 12, 0, 0, 0, time.UTC)
		rows := sqlmock.NewRows([]string{"id", "created_at"}).AddRow(15, fixedTime)

		mock.ExpectQuery("^INSERT INTO checkout_transactions").
			WithArgs(120.50, 4, "sig-abc-123").
			WillReturnRows(rows)

		router := setupRouter()
		w := httptest.NewRecorder()

		input := CheckoutTransaction{
			TotalAmount:  120.50,
			ItemCount:    4,
			KMSSignature: "sig-abc-123",
		}
		body, _ := json.Marshal(input)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/checkout", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var resp map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		data := resp["data"].(map[string]interface{})
		if data["id"].(float64) != 15 {
			t.Errorf("expected returned ID 15, got %v", data["id"])
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet mock expectations: %v", err)
		}
	})
}

// TestHandleAddCalendarSlot verifies POST /api/v1/calendar
func TestHandleAddCalendarSlot(t *testing.T) {
	t.Run("Invalid payload", func(t *testing.T) {
		router := setupRouter()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/calendar", bytes.NewBufferString("{invalid-json"))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Fallback when db is nil", func(t *testing.T) {
		db = nil
		router := setupRouter()
		w := httptest.NewRecorder()

		input := CalendarSlot{
			Title:    "Leg Day Workout",
			SlotTime: "04:30 PM",
		}
		body, _ := json.Marshal(input)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/calendar", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var resp map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		data := resp["data"].(map[string]interface{})
		if data["title"] != "Leg Day Workout" || data["slot_time"] != "04:30 PM" || data["synced"] != true {
			t.Errorf("unexpected calendar slot saved in fallback mode: %+v", data)
		}
	})

	t.Run("Success insert to DB", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("error creating mock database: %v", err)
		}
		defer mockDB.Close()

		db = mockDB
		defer func() { db = nil }()

		fixedTime := time.Date(2026, 6, 13, 12, 0, 0, 0, time.UTC)
		rows := sqlmock.NewRows([]string{"id", "created_at"}).AddRow(99, fixedTime)

		mock.ExpectQuery("^INSERT INTO calendar_slots").
			WithArgs("Leg Day Workout", "04:30 PM", true).
			WillReturnRows(rows)

		router := setupRouter()
		w := httptest.NewRecorder()

		input := CalendarSlot{
			Title:    "Leg Day Workout",
			SlotTime: "04:30 PM",
		}
		body, _ := json.Marshal(input)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/calendar", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var resp map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		data := resp["data"].(map[string]interface{})
		if data["id"].(float64) != 99 {
			t.Errorf("expected returned ID 99, got %v", data["id"])
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet mock expectations: %v", err)
		}
	})
}

// TestFetchCalendarSlots verifies GET /api/v1/calendar
func TestFetchCalendarSlots(t *testing.T) {
	t.Run("Fallback when db is nil", func(t *testing.T) {
		db = nil
		router := setupRouter()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/calendar", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var resp []CalendarSlot
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if len(resp) != 2 {
			t.Errorf("expected exactly 2 fallback calendar slots, got %d", len(resp))
		}

		if resp[0].Title != "Keto Breakfast Intake" || resp[1].Title != "Gym Metabolic Checkin" {
			t.Errorf("unexpected values in fallback calendar slots: %+v", resp)
		}
	})

	t.Run("Success retrieve from DB", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("error creating mock database: %v", err)
		}
		defer mockDB.Close()

		db = mockDB
		defer func() { db = nil }()

		fixedTime := time.Date(2026, 6, 13, 12, 0, 0, 0, time.UTC)
		rows := sqlmock.NewRows([]string{"id", "title", "slot_time", "synced", "created_at"}).
			AddRow(1, "Pilates Workout", "07:30 AM", true, fixedTime).
			AddRow(2, "Protein Shake Intake", "09:00 AM", false, fixedTime)

		mock.ExpectQuery("^SELECT id, title, slot_time, synced, created_at FROM calendar_slots").
			WillReturnRows(rows)

		router := setupRouter()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/calendar", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var resp []CalendarSlot
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if len(resp) != 2 {
			t.Errorf("expected exactly 2 calendar slots, got %d", len(resp))
		}

		if resp[0].Title != "Pilates Workout" || resp[1].Title != "Protein Shake Intake" {
			t.Errorf("unexpected calendar slots returned: %+v", resp)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet mock expectations: %v", err)
		}
	})
}

// TestMiddlewares verifies security headers and CORS configurations
func TestMiddlewares(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/health/metabolic", nil)
	router.ServeHTTP(w, req)

	// Verify CORS headers
	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("expected CORS Access-Control-Allow-Origin to be '*', got %s", w.Header().Get("Access-Control-Allow-Origin"))
	}

	// Verify Security headers
	if w.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Errorf("expected X-Content-Type-Options to be 'nosniff', got %s", w.Header().Get("X-Content-Type-Options"))
	}
	if w.Header().Get("X-Frame-Options") != "DENY" {
		t.Errorf("expected X-Frame-Options to be 'DENY', got %s", w.Header().Get("X-Frame-Options"))
	}
	if w.Header().Get("X-XSS-Protection") != "1; mode=block" {
		t.Errorf("expected X-XSS-Protection to be '1; mode=block', got %s", w.Header().Get("X-XSS-Protection"))
	}
}
