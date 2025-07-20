package main

import (
	"log"

	"github.com/Adit0507/nearby-friends/api"
	"github.com/Adit0507/nearby-friends/storage"
	"github.com/gorilla/mux"
)

func main() {
	redisClient, err := storage.NewRedisClient("localhost:6379", "")
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	cassandraClient, err := storage.NewCassandraClient([]string{"localhost"}, "nearby_friends")
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/user", api.CreateUser(redisClient, cassandraClient)).Methods("POST")
    r.HandleFunc("/user/{id}/friends", api.AddFriend(redisClient)).Methods("POST")
    r.HandleFunc("/user/{id}/location", api.UpdateLocation(redisClient, cassandraClient)).Methods("POST")
    r.HandleFunc("/user/{id}/nearby", api.GetNearbyFriends(redisClient)).Methods("GET")

}