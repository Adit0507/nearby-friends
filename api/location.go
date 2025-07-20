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
