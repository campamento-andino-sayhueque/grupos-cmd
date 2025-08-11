package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"grupos-cmd/internal/domain"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockEventRepository struct {
	mock.Mock
}

func (m *MockEventRepository) Save(ctx context.Context, event domain.Event) error {
	args := m.Called(event)
	return args.Error(0)
}

type MockEventPublisher struct {
	mock.Mock
}

func (m *MockEventPublisher) Publish(ctx context.Context, event domain.Event) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockEventPublisher) Close() error {
	args := m.Called()
	return args.Error(0)
}

// --- Tests ---

func TestHandler_CreateGrupo_Success(t *testing.T) {
	// --- Setup ---
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, r := gin.CreateTestContext(rec)

	reqBody := CreateGrupoRequest{
		Nombre:         "Test Group",
		FundacionFecha: "2024-01-01",
	}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/v1/grupos", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	mockRepo := new(MockEventRepository)
	mockPublisher := new(MockEventPublisher)
	h := NewHandler(mockRepo, mockPublisher)
	r.POST("/v1/grupos", h.CreateGrupo)

	// Mock expectations
	mockRepo.On("Save", mock.AnythingOfType("domain.Event")).Return(nil)
	mockPublisher.On("Publish", mock.AnythingOfType("domain.Event")).Return(nil)

	// --- Execute ---
	h.CreateGrupo(c)

	// --- Assert ---
	assert.Equal(t, http.StatusAccepted, rec.Code)

	var responseEvent domain.Event
	err := json.Unmarshal(rec.Body.Bytes(), &responseEvent)
	assert.NoError(t, err)
	assert.Equal(t, "GrupoCreado", responseEvent.Header.EventType)

	// Since the payload is now a struct, we need to unmarshal it into the correct type
	var payload domain.GrupoCreado
	payloadBytes, _ := json.Marshal(responseEvent.Payload)
	err = json.Unmarshal(payloadBytes, &payload)
	assert.NoError(t, err)
	assert.Equal(t, "Test Group", payload.Nombre)

	mockRepo.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}

func TestHandler_CreateGrupo_BadRequest(t *testing.T) {
	// --- Setup ---
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, r := gin.CreateTestContext(rec)

	reqBody := CreateGrupoRequest{
		Nombre:        "", // Invalid
		FundacionFecha: "2024-01-01",
	}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/v1/grupos", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	mockRepo := new(MockEventRepository)
	mockPublisher := new(MockEventPublisher)
	h := NewHandler(mockRepo, mockPublisher)
	r.POST("/v1/grupos", h.CreateGrupo)

	// --- Execute ---
	h.CreateGrupo(c)

	// --- Assert ---
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockRepo.AssertNotCalled(t, "Save", mock.Anything)
	mockPublisher.AssertNotCalled(t, "Publish", mock.Anything)
}
