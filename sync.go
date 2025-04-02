package matchguru

import (
	"context"
	"errors"
	"encoding/json"
	"log"
	"os"
	"net/http"
	"time"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/firestore"
	"cloud.google.com/go/logging"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
)

const (
	soccerLivescoreFirestoreID = "48lNVsFhqOWufeFzZBba"
	livescoreURL              = "https://api.sportmonks.com/v3/football/livescores?include=events.type;participants;scores;venue;lineups.player;statistics;referees.referee;state;periods&filters=fixturerefereeTypes:6"
)

func init() {
	functions.CloudEvent("Sync", sync)
}

func sync(ctx context.Context, e event.Event) error {
	logger := initLogger(ctx)

	logger.Println("sync function called")

	projectID, err := metadata.ProjectIDWithContext(ctx)
	if err != nil {
		logger.Fatalf("failed to get project ID: %v", err)
	}

	firestoreClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		logger.Printf("failed to initiate firestore client: %v", err)
		return err
	}
	defer firestoreClient.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, livescoreURL, http.NoBody)
	if err != nil {
		logger.Printf("failed to create request: %v", err)
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", os.Getenv("SPORTMONKS_API_KEY"))

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Printf("failed to send request: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Printf("unexpected status code: %d", resp.StatusCode)
		return errors.New("unexpected status code")
	}

	soccerLivescoreData := map[string]any{}
	if err := json.NewDecoder(resp.Body).Decode(&soccerLivescoreData); err != nil {
		logger.Printf("error while decoding request: %v", err)
		return err
	}

	if _, ok := soccerLivescoreData["data"]; !ok {
		logger.Printf("no livescore data found in response")
		soccerLivescoreData = map[string]any{
			"data": []any{},
		}
	} else {
		// only keep the data field
		soccerLivescoreData = map[string]any{
			"data": soccerLivescoreData["data"],
		}
	}

	_, err = firestoreClient.Collection("livescores").Doc(soccerLivescoreFirestoreID).Set(ctx, soccerLivescoreData)
	if err != nil {
		logger.Printf("failed to write to firestore: %v", err)
		return err
	}

	logger.Printf("request successful %+v", resp)

	return nil
}

func initLogger(ctx context.Context) *log.Logger {
	projectID, err := metadata.ProjectIDWithContext(ctx)
	if err != nil {
		log.Fatalf("failed to get project ID: %v", err)
	}
	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("failed to create logging client: %v", err)
	}
	return client.Logger("livescore").StandardLogger(logging.Info)
}
