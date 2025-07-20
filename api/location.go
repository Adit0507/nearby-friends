package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Adit0507/nearby-friends/models"
	"github.com/Adit0507/nearby-friends/storage"
	"github.com/gorilla/mux"
)

func UpdateLocation(redisClient *storage.RedisClient, casssandraClient *storage.CassandraClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := mux.Vars(r)["id"]

		var loc models.Location
		if err := json.NewDecoder(r.Body).Decode(&loc); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		loc.Timestamp = time.Now().Unix()

		// updatin location in redis
		if err := redisClient.UpdateLocation(r.Context(), userID, loc.Latitude, loc.Longitude); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// save to cassandra
		if err := casssandraClient.SaveLocation(userID, loc.Latitude, loc.Longitude, loc.Timestamp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// publish location update
		if err := redisClient.PublishLocationUpdate(r.Context(), userID, loc.Latitude, loc.Longitude); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func GetNearbyFriends(redisClient *storage.RedisClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := mux.Vars(r)["id"]
		radius := 1000.0

		// get friends
		friends, err := redisClient.Client.SMembers(r.Context(), "friends: "+userID).Result()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// gettin nearby users
		nearby, err := redisClient.GetNearbyFriends(r.Context(), userID, radius)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// filter nearby friends
		var ans []models.NearbyFriend
		for _, friend := range friends {
			for _, nearbyID := range nearby {
				if friend == nearbyID {
					name, _ := redisClient.Client.HGet(r.Context(), "users", friend).Result()
					ans = append(ans, models.NearbyFriend{UserID: friend, Name: name, Distance: 0})
				}
			}
		}

		json.NewEncoder(w).Encode(ans)
	}
}
