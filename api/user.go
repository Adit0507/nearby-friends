package api

import (
	"encoding/json"
	"net/http"

	"github.com/Adit0507/nearby-friends/models"
	"github.com/Adit0507/nearby-friends/storage"
)

func CreateUser(redisClient *storage.RedisClient, cassandraClient *storage.CassandraClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return 
		}

		err := redisClient.Client.HSet(r.Context(), "users", user.ID, user.Name).Err()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return 
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}
