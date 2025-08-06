package repository

import (
	"context"
	"grupos-cmd/internal/config"
	"grupos-cmd/internal/domain"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

// EventRepository defines the interface for persisting events.
type EventRepository interface {
	Save(ctx context.Context, event domain.Event) error
}

// firestoreRepository is the Firestore implementation of the EventRepository.
type firestoreRepository struct {
	client     *firestore.Client
	collection string
}

import (
	"context"
	"grupos-cmd/internal/config"
	"grupos-cmd/internal/domain"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewFirestoreRepository creates a new Firestore event repository.
func NewFirestoreRepository(ctx context.Context, cfg *config.Config) (EventRepository, error) {
	var opts []option.ClientOption
	if os.Getenv("FIRESTORE_EMULATOR_HOST") != "" {
		opts = append(opts, option.WithEndpoint(os.Getenv("FIRESTORE_EMULATOR_HOST")))
		opts = append(opts, option.WithoutAuthentication())
		opts = append(opts, option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())))
	}

	client, err := firestore.NewClient(ctx, cfg.GoogleCloudProjectID, opts...)
	if err != nil {
		return nil, err
	}

	return &firestoreRepository{
		client:     client,
		collection: cfg.FirestoreCollection,
	}, nil
}

// Save persists an event to Firestore.
func (r *firestoreRepository) Save(ctx context.Context, event domain.Event) error {
	_, err := r.client.Collection(r.collection).Doc(event.Header.EventID).Set(ctx, event)
	return err
}
