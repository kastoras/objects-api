package rpc

import (
	"encoding/json"
	"log"
	"time"
)

func (r *RPCServer) SetCache(cacheKey string, images []Image) {

	cacheValue, err := json.Marshal(images)
	if err != nil {
		log.Fatalf("Error marshalling images: %v", err)
		return
	}

	r.Cache.Set(cacheKey, string(cacheValue), 24*time.Hour)
}
