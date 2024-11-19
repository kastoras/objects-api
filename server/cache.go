package server

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/kastoras/go-utilities"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func CreateCacheClient() error {

	redisHost, err := utilities.GetEnvParam("REDIS_HOST", "")
	if err != nil {
		return fmt.Errorf("error getting REDIS_HOST: %v", err)
	}

	redisPass, err := utilities.GetEnvParam("REDIS_PASS", "")
	if err != nil {
		return fmt.Errorf("error getting REDIS_PASS: %v", err)
	}

	redisDb, err := utilities.GetEnvParam("REDIS_DB", "")
	if err != nil {
		return fmt.Errorf("error getting REDIS_DB: %v", err)
	}

	db, err := strconv.Atoi(redisDb)
	if err != nil {
		return fmt.Errorf("error parsing REDIS_DB: %v", err)
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPass,
		DB:       db,
	})

	return nil
}

func CloseCacheClient() {
	rdb.Close()
}

func (cache *Cache) Set(key string, value string, expiration time.Duration) {
	ctx := context.Background()

	rdb.Set(ctx, key, value, expiration)
}

func (cache *Cache) Get(key string) (string, error) {
	ctx := context.Background()
	val := rdb.Get(ctx, key)
	result, err := val.Result()
	if err != nil {
		return result, err
	}
	return result, nil
}

func (cache *Cache) Remove(key string) {
	ctx := context.Background()
	rdb.Del(ctx, key)
}

func (cache *Cache) Ping() string {

	ctx := context.Background()

	status, err := rdb.Ping(ctx).Result()
	cacheStatus := ""
	if err != nil {
		cacheStatus = fmt.Sprintf("Cannot connect to cache server, error : %v", err.Error())
	} else {
		cacheStatus = fmt.Sprintf("Ping - %s", status)
	}

	return cacheStatus
}
