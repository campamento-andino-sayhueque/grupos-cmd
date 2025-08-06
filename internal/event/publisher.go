package event

import (
	"context"
	"encoding/json"
	"grupos-cmd/internal/config"
	"grupos-cmd/internal/domain"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"
	"github.com/ThreeDotsLabs/watermill/message"
	"google.golang.org/api/option"
)

// Publisher defines the interface for publishing events.
type Publisher interface {
	Publish(ctx context.Context, event domain.Event) error
	Close() error
}

// watermillPublisher is the Watermill implementation of the Publisher.
type watermillPublisher struct {
	publisher message.Publisher
}

import (
	"context"
	"encoding/json"
	"grupos-cmd/internal/config"
	"grupos-cmd/internal/domain"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"
	"github.com/ThreeDotsLabs/watermill/message"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewWatermillPublisher creates a new Watermill publisher.
func NewWatermillPublisher(cfg *config.Config) (Publisher, error) {
	var clientOptions []option.ClientOption
	if os.Getenv("PUBSUB_EMULATOR_HOST") != "" {
		clientOptions = append(clientOptions, option.WithEndpoint(os.Getenv("PUBSUB_EMULATOR_HOST")))
		clientOptions = append(clientOptions, option.WithoutAuthentication())
		clientOptions = append(clientOptions, option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())))
	}

	pub, err := googlecloud.NewPublisher(
		googlecloud.PublisherConfig{
			ProjectID:     cfg.GoogleCloudProjectID,
			ClientOptions: clientOptions,
		},
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		return nil, err
	}

	return &watermillPublisher{publisher: pub}, nil
}

// Publish publishes an event to the Pub/Sub topic.
func (p *watermillPublisher) Publish(ctx context.Context, event domain.Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := message.NewMessage(event.Header.EventID, payload)
	msg.Metadata.Set("eventType", event.Header.EventType)

	return p.publisher.Publish("grupos", msg)
}

// Close closes the publisher.
func (p *watermillPublisher) Close() error {
	return p.publisher.Close()
}
