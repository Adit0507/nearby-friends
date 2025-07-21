package websocket

import (
	"context"
	"log"
	"net/http"

	"github.com/Adit0507/nearby-friends/storage"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func HandleWebSocket(redisClient *storage.RedisClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := mux.Vars(r)["userID"]

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Websocket upgrade failed: %v", err)
			return
		}
		defer conn.Close()

		// gettin user's friends
		friends, err := redisClient.Client.SMembers(r.Context(), "friends: "+userID).Result()
		if err != nil {
			log.Printf("Failed to get friends: %v", err)
			return
		}

		// subscribin to friend's location updates
		ctx := context.Background()
		pubsub := redisClient.Subscribe(ctx, "location: "+userID)
		for _, friend := range friends {
			pubsub.Subscribe(ctx, "location: "+friend)
		}
		defer pubsub.Close()

		// handlin websocket messages
		go func() {
			for {
				_, _, err := conn.ReadMessage()
				if err != nil {
					log.Printf("websocket read error: %v", err)
					return
				}
			}
		}()

		// handle pub/sub messages
		for msg := range pubsub.Channel() {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
				log.Printf("Websocket write error: %v", err)
				return
			}
		}

	}
}
