package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/redismock/v8"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// setupMockApp skapar en App med mockade beroenden
func setupMockApp(t *testing.T) (*App, sqlmock.Sqlmock, redismock.ClientMock) {
	// Mock PostgreSQL med MonitorPingsOption aktiverat
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Kunde inte skapa sql mock: %v", err)
	}

	// Mock Redis
	redisClient, redisMock := redismock.NewClientMock()

	app := &App{
		DB:    db,
		Redis: redisClient,
		Ctx:   context.Background(),
	}

	return app, mock, redisMock
}

// TestGetEnv testar miljövariabel-funktionen
func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "Använd default när env saknas",
			key:          "MISSING_KEY",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
		{
			name:         "Använd env värde när satt",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "custom",
			expected:     "custom",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				t.Setenv(tt.key, tt.envValue)
			}
			result := getEnv(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestHealthHandler testar health endpoint
func TestHealthHandler(t *testing.T) {
	app, dbMock, redisMock := setupMockApp(t)
	defer app.DB.Close()

	tests := []struct {
		name           string
		dbPingErr      error
		redisPingErr   error
		expectedStatus string
		expectedDB     string
	}{
		{
			name:           "Allt fungerar",
			dbPingErr:      nil,
			redisPingErr:   nil,
			expectedStatus: "healthy",
			expectedDB:     "healthy",
		},
		{
			name:           "Databas nere",
			dbPingErr:      sql.ErrConnDone,
			redisPingErr:   nil,
			expectedStatus: "degraded",
			expectedDB:     "unhealthy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks för databas ping
			if tt.dbPingErr != nil {
				dbMock.ExpectPing().WillReturnError(tt.dbPingErr)
			} else {
				dbMock.ExpectPing()
			}

			// Setup mocks för Redis ping
			if tt.redisPingErr == nil {
				redisMock.ExpectPing().SetVal("PONG")
			} else {
				redisMock.ExpectPing().SetErr(tt.redisPingErr)
			}

			// Skapa request
			req := httptest.NewRequest("GET", "/health", nil)
			w := httptest.NewRecorder()

			// Kör handler
			app.healthHandler(w, req)

			// Verifiera response
			assert.Equal(t, http.StatusOK, w.Code)

			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, response["status"])
			assert.Equal(t, tt.expectedDB, response["database"])
		})
	}
}

// TestGetEntriesHandler_CacheHit testar när cache träffar
func TestGetEntriesHandler_CacheHit(t *testing.T) {
	app, _, redisMock := setupMockApp(t)
	defer app.DB.Close()

	// Mockdata
	cachedEntries := []Entry{
		{ID: 1, Name: "Test User", Message: "Test message", CreatedAt: time.Now()},
	}
	cachedJSON, _ := json.Marshal(cachedEntries)

	// Setup Redis mock
	redisMock.ExpectGet("entries:all").SetVal(string(cachedJSON))

	// Skapa request
	req := httptest.NewRequest("GET", "/api/entries", nil)
	w := httptest.NewRecorder()

	// Kör handler
	app.getEntriesHandler(w, req)

	// Verifiera
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "HIT", w.Header().Get("X-Cache"))

	var entries []Entry
	err := json.NewDecoder(w.Body).Decode(&entries)
	assert.NoError(t, err)
	assert.Len(t, entries, 1)
	assert.Equal(t, "Test User", entries[0].Name)
}

// TestGetEntriesHandler_CacheMiss testar när cache missar
func TestGetEntriesHandler_CacheMiss(t *testing.T) {
	app, dbMock, redisMock := setupMockApp(t)
	defer app.DB.Close()

	now := time.Now()

	// Setup Redis mock (cache miss)
	redisMock.ExpectGet("entries:all").RedisNil()

	// Setup DB mock
	rows := sqlmock.NewRows([]string{"id", "name", "message", "created_at"}).
		AddRow(1, "Test User", "Test message", now).
		AddRow(2, "Another User", "Another message", now)

	dbMock.ExpectQuery("SELECT id, name, message, created_at FROM entries").
		WillReturnRows(rows)

	// Setup Redis Set mock
	redisMock.ExpectSet("entries:all", sqlmock.AnyArg(), 30*time.Second).SetVal("OK")

	// Skapa request
	req := httptest.NewRequest("GET", "/api/entries", nil)
	w := httptest.NewRecorder()

	// Kör handler
	app.getEntriesHandler(w, req)

	// Verifiera
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "MISS", w.Header().Get("X-Cache"))

	var entries []Entry
	err := json.NewDecoder(w.Body).Decode(&entries)
	assert.NoError(t, err)
	assert.Len(t, entries, 2)
	assert.Equal(t, "Test User", entries[0].Name)

	// Verifiera att alla förväntningar uppfylldes
	assert.NoError(t, dbMock.ExpectationsWereMet())
}

// TestCreateEntryHandler_Success testar lyckad skapande
func TestCreateEntryHandler_Success(t *testing.T) {
	app, dbMock, redisMock := setupMockApp(t)
	defer app.DB.Close()

	now := time.Now()
	entry := Entry{
		Name:    "Test User",
		Message: "Test message",
	}

	// Setup DB mock
	rows := sqlmock.NewRows([]string{"id", "created_at"}).
		AddRow(1, now)

	dbMock.ExpectQuery("INSERT INTO entries").
		WithArgs(entry.Name, entry.Message).
		WillReturnRows(rows)

	// Setup Redis mocks
	redisMock.ExpectDel("entries:all").SetVal(1)
	redisMock.ExpectIncr("stats:total_entries").SetVal(1)

	// Skapa request
	body, _ := json.Marshal(entry)
	req := httptest.NewRequest("POST", "/api/entries", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Kör handler
	app.createEntryHandler(w, req)

	// Verifiera
	assert.Equal(t, http.StatusCreated, w.Code)

	var response Entry
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, 1, response.ID)
	assert.Equal(t, "Test User", response.Name)
	assert.Equal(t, "Test message", response.Message)

	assert.NoError(t, dbMock.ExpectationsWereMet())
}

// TestCreateEntryHandler_InvalidData testar ogiltig data
func TestCreateEntryHandler_InvalidData(t *testing.T) {
	app, _, _ := setupMockApp(t)
	defer app.DB.Close()

	tests := []struct {
		name           string
		entry          Entry
		expectedStatus int
	}{
		{
			name:           "Tomt namn",
			entry:          Entry{Name: "", Message: "Test"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Tomt meddelande",
			entry:          Entry{Name: "Test", Message: ""},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Både namn och meddelande tomma",
			entry:          Entry{Name: "", Message: ""},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.entry)
			req := httptest.NewRequest("POST", "/api/entries", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			app.createEntryHandler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

// TestCreateEntryHandler_InvalidJSON testar felaktig JSON
func TestCreateEntryHandler_InvalidJSON(t *testing.T) {
	app, _, _ := setupMockApp(t)
	defer app.DB.Close()

	req := httptest.NewRequest("POST", "/api/entries", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	app.createEntryHandler(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestStatsHandler testar statistik endpoint
func TestStatsHandler(t *testing.T) {
	app, dbMock, redisMock := setupMockApp(t)
	defer app.DB.Close()

	// Setup DB mock
	rows := sqlmock.NewRows([]string{"count"}).AddRow(42)
	dbMock.ExpectQuery("SELECT COUNT").WillReturnRows(rows)

	// Setup Redis mocks
	redisMock.ExpectGet("stats:total_entries").SetVal("50")
	redisMock.ExpectInfo("stats").SetVal("# Stats\r\ntotal_keys:10")

	// Skapa request
	req := httptest.NewRequest("GET", "/api/stats", nil)
	w := httptest.NewRecorder()

	// Kör handler
	app.statsHandler(w, req)

	// Verifiera
	assert.Equal(t, http.StatusOK, w.Code)

	var stats map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&stats)
	assert.NoError(t, err)
	assert.Equal(t, float64(42), stats["total_entries_db"])
	assert.Equal(t, "50", stats["total_entries_created"])
	assert.Equal(t, true, stats["cache_available"])

	assert.NoError(t, dbMock.ExpectationsWereMet())
}

// TestCORSMiddleware testar CORS middleware
func TestCORSMiddleware(t *testing.T) {
	// Skapa en enkel handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Applicera CORS middleware
	wrappedHandler := corsMiddleware(handler)

	t.Run("OPTIONS request", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/test", nil)
		w := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "GET, POST, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
		assert.Equal(t, "Content-Type", w.Header().Get("Access-Control-Allow-Headers"))
	})

	t.Run("GET request", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	})
}

// TestRouting testar att routes är korrekt konfigurerade
func TestRouting(t *testing.T) {
	app, _, _ := setupMockApp(t)
	defer app.DB.Close()

	r := mux.NewRouter()
	r.Use(corsMiddleware)
	r.HandleFunc("/health", app.healthHandler).Methods("GET")
	r.HandleFunc("/api/entries", app.getEntriesHandler).Methods("GET")
	r.HandleFunc("/api/entries", app.createEntryHandler).Methods("POST")
	r.HandleFunc("/api/stats", app.statsHandler).Methods("GET")

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "Health endpoint GET",
			method:         "GET",
			path:           "/health",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Stats endpoint GET",
			method:         "GET",
			path:           "/api/stats",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Ogiltigt endpoint",
			method:         "GET",
			path:           "/invalid",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
