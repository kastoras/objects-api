package server

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	"github.com/kastoras/go-utilities"
	"google.golang.org/api/option"
)

func InitStorageClient() (*storage.Client, error) {
	serviceAccountKey, err := utilities.GetEnvParam("GOOGLE_STORAGE_KEY", "")
	if err != nil {
		return nil, fmt.Errorf("failed to get GOOGLE_STORAGE_KEY: %w", err)
	}

	options := option.WithCredentialsFile(serviceAccountKey)

	ctx := context.Background()
	client, err := storage.NewClient(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	//defer client.Close()

	return client, nil
}
