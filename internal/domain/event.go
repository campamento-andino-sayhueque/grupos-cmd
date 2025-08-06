package domain

import (
	"time"

	"github.com/google/uuid"
)

// EventHeader contains common fields for all events.
type EventHeader struct {
	EventID     string    `json:"eventId" firestore:"eventId"`
	EventType   string    `json:"eventType" firestore:"eventType"`
	Timestamp   time.Time `json:"timestamp" firestore:"timestamp"`
	AggregateID string    `json:"aggregateId" firestore:"aggregateId"`
}

// Event represents a domain event.
type Event struct {
	Header  EventHeader `json:"header" firestore:"header"`
	Payload interface{} `json:"payload" firestore:"payload"`
}

// NewEvent creates a new event with a header.
func NewEvent(aggregateID, eventType string, payload interface{}) Event {
	return Event{
		Header: EventHeader{
			EventID:     uuid.New().String(),
			EventType:   eventType,
			Timestamp:   time.Now().UTC(),
			AggregateID: aggregateID,
		},
		Payload: payload,
	}
}

// --- Event Payloads ---

// GrupoCreado is the payload for the GrupoCreado event.
type GrupoCreado struct {
	Nombre        string `json:"nombre"`
	FundacionFecha string `json:"fundacionFecha"`
}

// GrupoActualizado is the payload for the GrupoActualizado event.
type GrupoActualizado struct {
	Nombre string `json:"nombre"`
}

// DirigenteAsignadoAGrupo is the payload for the DirigenteAsignadoAGrupo event.
type DirigenteAsignadoAGrupo struct {
	DirigenteID string `json:"dirigenteId"`
}

// DirigenteRemovidoDeGrupo is the payload for the DirigenteRemovidoDeGrupo event.
type DirigenteRemovidoDeGrupo struct {
	DirigenteID string `json:"dirigenteId"`
}
