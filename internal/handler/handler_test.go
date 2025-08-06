package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"grupos-cmd/internal/domain"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
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
	e := echo.New()
	reqBody := CreateGrupoRequest{
		Nombre:         "Test Group",
		FundacionFecha: "2024-01-01",
	}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/v1/grupos", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRepo := new(MockEventRepository)
	mockPublisher := new(MockEventPublisher)

	// Mock expectations
	mockRepo.On("Save", mock.AnythingOfType("domain.Event")).Return(nil)
	mockPublisher.On("Publish", mock.AnythingOfType("domain.Event")).Return(nil)

	h := NewHandler(mockRepo, mockPublisher)

	// --- Execute ---
	err := h.CreateGrupo(c)

	// --- Assert ---
	assert.NoError(t, err)
	assert.Equal(t, http.StatusAccepted, rec.Code)

	var responseEvent domain.Event
	err = json.Unmarshal(rec.Body.Bytes(), &responseEvent)
	assert.NoError(t, err)
	assert.Equal(t, "GrupoCreado", responseEvent.Header.EventType)

	payload, ok := responseEvent.Payload.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "Test Group", payload["nombre"])

	mockRepo.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}

func TestHandler_CreateGrupo_BadRequest(t *testing.T) {
	// --- Setup ---
	e := echo.New()
	reqBody := CreateGrupoRequest{
		Nombre:        "", // Invalid
		FundacionFecha: "2024-01-01",
	}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/v1/grupos", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRepo := new(MockEventRepository)
	mockPublisher := new(MockEventPublisher)
	h := NewHandler(mockRepo, mockPublisher)

	// --- Execute ---
	err := h.CreateGrupo(c)

	// --- Assert ---
	assert.Error(t, err)
	httpError, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpError.Code)
	mockRepo.AssertNotCalled(t, "Save", mock.Anything)
	mockPublisher.AssertNotCalled(t, "Publish", mock.Anything)
}
