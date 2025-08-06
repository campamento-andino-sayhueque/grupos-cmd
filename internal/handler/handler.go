package handler

import (
	"grupos-cmd/internal/domain"
	"grupos-cmd/internal/event"
	"grupos-cmd/internal/repository"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// --- Request Structs ---

// CreateGrupoRequest is the request payload for creating a new group.
type CreateGrupoRequest struct {
	Nombre        string `json:"nombre" validate:"required"`
	FundacionFecha string `json:"fundacionFecha" validate:"required"`
}

// UpdateGrupoRequest is the request payload for updating a group.
type UpdateGrupoRequest struct {
	Nombre string `json:"nombre" validate:"required"`
}

// AssignDirigenteRequest is the request payload for assigning a leader to a group.
type AssignDirigenteRequest struct {
	DirigenteID string `json:"dirigenteId" validate:"required"`
}

// --- Handler ---

// Handler holds the dependencies for the HTTP handlers.
type Handler struct {
	repo      repository.EventRepository
	publisher event.Publisher
}

// NewHandler creates a new Handler.
func NewHandler(repo repository.EventRepository, publisher event.Publisher) *Handler {
	return &Handler{
		repo:      repo,
		publisher: publisher,
	}
}

// RegisterRoutes registers the HTTP routes for the service.
func (h *Handler) RegisterRoutes(e *echo.Echo) {
	v1 := e.Group("/v1")
	grupos := v1.Group("/grupos")

	grupos.POST("", h.CreateGrupo)
	grupos.PUT("/:id", h.UpdateGrupo)
	grupos.POST("/:id/dirigentes", h.AssignDirigente)
	grupos.DELETE("/:id/dirigentes/:dirigenteId", h.RemoveDirigente)
}

// CreateGrupo handles the creation of a new group.
// @Summary Create a new group
// @Description Creates a new group and publishes a GrupoCreado event.
// @Tags grupos
// @Accept json
// @Produce json
// @Param grupo body CreateGrupoRequest true "Group information"
// @Success 202 {object} domain.Event
// @Router /v1/grupos [post]
func (h *Handler) CreateGrupo(c echo.Context) error {
	var req CreateGrupoRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Basic validation
	if req.Nombre == "" || req.FundacionFecha == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "nombre and fundacionFecha are required")
	}
	if _, err := time.Parse("2006-01-02", req.FundacionFecha); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid fundacionFecha format, expected YYYY-MM-DD")
	}

	aggregateID := uuid.New().String()
	event := domain.NewEvent(
		aggregateID,
		"GrupoCreado",
		domain.GrupoCreado{
			Nombre:        req.Nombre,
			FundacionFecha: req.FundacionFecha,
		},
	)

	if err := h.repo.Save(c.Request().Context(), event); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save event")
	}

	if err := h.publisher.Publish(c.Request().Context(), event); err != nil {
		// Log the error but don't fail the request, as the event is already saved.
		c.Logger().Error("failed to publish event: ", err)
	}

	return c.JSON(http.StatusAccepted, event)
}

// UpdateGrupo handles the update of a group.
// @Summary Update a group's name
// @Description Updates a group's name and publishes a GrupoActualizado event.
// @Tags grupos
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Param grupo body UpdateGrupoRequest true "Group new name"
// @Success 202 {object} domain.Event
// @Router /v1/grupos/{id} [put]
func (h *Handler) UpdateGrupo(c echo.Context) error {
	var req UpdateGrupoRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if req.Nombre == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "nombre is required")
	}

	aggregateID := c.Param("id")
	event := domain.NewEvent(
		aggregateID,
		"GrupoActualizado",
		domain.GrupoActualizado{Nombre: req.Nombre},
	)

	if err := h.repo.Save(c.Request().Context(), event); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save event")
	}

	if err := h.publisher.Publish(c.Request().Context(), event); err != nil {
		c.Logger().Error("failed to publish event: ", err)
	}

	return c.JSON(http.StatusAccepted, event)
}

// AssignDirigente handles assigning a leader to a group.
// @Summary Assign a leader to a group
// @Description Assigns a leader to a group and publishes a DirigenteAsignadoAGrupo event.
// @Tags grupos
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Param dirigente body AssignDirigenteRequest true "Leader ID"
// @Success 202 {object} domain.Event
// @Router /v1/grupos/{id}/dirigentes [post]
func (h *Handler) AssignDirigente(c echo.Context) error {
	var req AssignDirigenteRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if req.DirigenteID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "dirigenteId is required")
	}

	aggregateID := c.Param("id")
	event := domain.NewEvent(
		aggregateID,
		"DirigenteAsignadoAGrupo",
		domain.DirigenteAsignadoAGrupo{DirigenteID: req.DirigenteID},
	)

	if err := h.repo.Save(c.Request().Context(), event); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save event")
	}

	if err := h.publisher.Publish(c.Request().Context(), event); err != nil {
		c.Logger().Error("failed to publish event: ", err)
	}

	return c.JSON(http.StatusAccepted, event)
}

// RemoveDirigente handles removing a leader from a group.
// @Summary Remove a leader from a group
// @Description Removes a leader from a group and publishes a DirigenteRemovidoDeGrupo event.
// @Tags grupos
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Param dirigenteId path string true "Leader ID"
// @Success 202 {object} domain.Event
// @Router /v1/grupos/{id}/dirigentes/{dirigenteId} [delete]
func (h *Handler) RemoveDirigente(c echo.Context) error {
	aggregateID := c.Param("id")
	dirigenteID := c.Param("dirigenteId")

	event := domain.NewEvent(
		aggregateID,
		"DirigenteRemovidoDeGrupo",
		domain.DirigenteRemovidoDeGrupo{DirigenteID: dirigenteID},
	)

	if err := h.repo.Save(c.Request().Context(), event); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save event")
	}

	if err := h.publisher.Publish(c.Request().Context(), event); err != nil {
		c.Logger().Error("failed to publish event: ", err)
	}

	return c.JSON(http.StatusAccepted, event)
}
