package db

import (
	"context"
	"fmt"
	"time"

	"github.com/VidarHUN/app_server/internal/config"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type Room struct {
	Id     string              `json:"id"`
	Users  []User              `json:"users"`
	Server config.QuicrqServer `json:"server"`
}

type User struct {
	Id   string          `json:"id"`
	Conn *websocket.Conn `json:"connection"`
}

// NewRedisClient creates and returns a new Redis client instance.
func NewRedisClient(addr, password string, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Failed to connect to Redis:", err)
		return client, err
	}

	fmt.Println("Redis connected:", pong)

	return client, err
}
