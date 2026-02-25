package cache

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(addr string, db int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("[redis] failed to connect: addr=%s db=%d error=%v", addr, db, err)
	}

	log.Printf("[redis] connected: addr=%s db=%d", addr, db)
	return rdb
}
