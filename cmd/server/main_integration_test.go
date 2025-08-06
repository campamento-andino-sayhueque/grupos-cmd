package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"grupos-cmd/internal/domain"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	baseAddress = "http://localhost:8080"
	projectID   = "test-project"
	topicID     = "grupos"
)

var (
	firestoreClient *firestore.Client
	pubsubClient    *pubsub.Client
)

func TestMain(m *testing.M) {
	// --- Setup ---
	ctx := context.Background()
	var err error

	// Firestore Emulator Client
	firestoreClient, err = firestore.NewClient(ctx, projectID,
		option.WithEndpoint("localhost:8081"),
		option.WithoutAuthentication(),
		option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	if err != nil {
		fmt.Printf("Failed to create firestore client: %v\n", err)
		os.Exit(1)
	}

	// Pub/Sub Emulator Client
	pubsubClient, err = pubsub.NewClient(ctx, projectID,
		option.WithEndpoint("localhost:8085"),
		option.WithoutAuthentication(),
		option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	if err != nil {
		fmt.Printf("Failed to create pubsub client: %v\n", err)
		os.Exit(1)
	}

	// Create Pub/Sub topic
	topic, err := pubsubClient.CreateTopic(ctx, topicID)
	if err != nil {
		// Ignore if topic already exists
		if pubsub.IsTopicExistsError(err) {
			topic = pubsubClient.Topic(topicID)
		} else {
			fmt.Printf("Failed to create pubsub topic: %v\n", err)
			os.Exit(1)
		}
	}

	// --- Run Tests ---
	code := m.Run()

	// --- Teardown ---
	_ = topic.Delete(ctx)
	_ = firestoreClient.Close()
	_ = pubsubClient.Close()

	os.Exit(code)
}

func TestCreateGrupo_Integration(t *testing.T) {
	// --- Setup ---
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a subscription to the topic to receive messages
	sub, err := pubsubClient.CreateSubscription(ctx, "test-sub", pubsub.SubscriptionConfig{Topic: pubsubClient.Topic(topicID)})
	require.NoError(t, err)
	defer sub.Delete(ctx)

	// --- Execute ---
	reqBody := map[string]string{
		"nombre":         "Integration Test Group",
		"fundacionFecha": "2024-01-01",
	}
	jsonBody, _ := json.Marshal(reqBody)
	resp, err := http.Post(baseAddress+"/v1/grupos", "application/json", bytes.NewReader(jsonBody))
	require.NoError(t, err)
	defer resp.Body.Close()

	// --- Assert HTTP Response ---
	assert.Equal(t, http.StatusAccepted, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var createdEvent domain.Event
	err = json.Unmarshal(body, &createdEvent)
	require.NoError(t, err)
	assert.Equal(t, "GrupoCreado", createdEvent.Header.EventType)

	// --- Assert Firestore ---
	doc, err := firestoreClient.Collection("eventos").Doc(createdEvent.Header.EventID).Get(ctx)
	require.NoError(t, err)
	assert.True(t, doc.Exists())

	// --- Assert Pub/Sub ---
	var receivedMsg *pubsub.Message
	cctx, cancelSub := context.WithCancel(ctx)
	err = sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		receivedMsg = msg
		msg.Ack()
		cancelSub()
	})
	require.NoError(t, err)
	require.NotNil(t, receivedMsg)

	var publishedEvent domain.Event
	err = json.Unmarshal(receivedMsg.Data, &publishedEvent)
	require.NoError(t, err)
	assert.Equal(t, createdEvent.Header.EventID, publishedEvent.Header.EventID)
	assert.Equal(t, "GrupoCreado", publishedEvent.Header.EventType)
}
