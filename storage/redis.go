package storage

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(addr, password string) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}

	return &RedisClient{Client: client}, nil
}

func (r *RedisClient) UpdateLocation(ctx context.Context, userId string, lat, lon float64) error {
	// storin location in Redis with ttl
	return r.Client.GeoAdd(ctx, "user_locations", &redis.GeoLocation{
		Name:      userId,
		Latitude:  lat,
		Longitude: lon,
	}).Err()
}

func (r *RedisClient) GetNearbyFriends(ctx context.Context, userId string, radius float64) ([]string, error) {
	res, err := r.Client.GeoRadiusByMember(ctx, "user_locations", userId, &redis.GeoRadiusQuery{
		Radius:    radius,
		Unit:      "m",
		WithCoord: true,
		Sort:      "ASC",
	}).Result()
	if err != nil {
		return nil, err
	}

	var nearby []string
	for _, loc := range res {
		if loc.Name != userId {
			nearby = append(nearby, loc.Name)
		}
	}
	return nearby, nil
}

func (r *RedisClient) PublishLocationUpdate(ctx context.Context, userID string, lat, lon float64) error {
	return r.Client.Publish(ctx, "location: "+userID, fmt.Sprintf("%f, %f", lat, lon)).Err()
}

func (r *RedisClient) Subscribe(ctx context.Context, channel string) *redis.PubSub {
	return r.Client.Subscribe(ctx, channel)
}
