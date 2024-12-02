package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type RPCPayload struct {
	Bucket string
	Folder string
}

type Image struct {
	ImageURL  string
	ImageName string
}

var cacheKey string

func (r *RPCServer) GetImages(payload RPCPayload, response *[]Image) error {

	cacheKey = createCacheKey(payload.Bucket, payload.Folder)

	cachedData, err := r.Cache.Get(cacheKey)
	if err == nil && cachedData != "" {
		var cachedImages []Image
		err = json.Unmarshal([]byte(cachedData), &cachedImages)
		if err == nil {
			*response = cachedImages
			return nil
		} else {
			log.Printf("Error unmarshalling cached data: %v", err)
		}
	} else if err != nil {
		log.Printf("Error getting data from cache: %v", err)
	}

	images, err := r.listBucketFiles(payload.Bucket, payload.Folder)
	if err != nil {
		return err
	}

	*response = images

	return nil
}

func (r *RPCServer) listBucketFiles(bucket string, folder string) ([]Image, error) {

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	it := r.StorageClient.Bucket(bucket).Objects(ctx, &storage.Query{
		Prefix:                   folder,
		IncludeFoldersAsPrefixes: false,
	})

	images := []Image{}

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break // No more objects
		}
		if err != nil {
			log.Fatalf("Error listing objects: %v", err)
		}

		if attrs.Name == fmt.Sprintf("%s/", folder) {
			continue
		}

		publicURL := getPublicURL(bucket, attrs.Name)

		images = append(images, Image{ImageURL: publicURL, ImageName: attrs.Name})
	}

	r.SetCache(cacheKey, images)

	return images, nil
}

func getPublicURL(bucket string, objectName string) string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucket, objectName)
}

func createCacheKey(bucket string, folder string) string {
	return fmt.Sprintf("%s:%s", bucket, folder)
}
