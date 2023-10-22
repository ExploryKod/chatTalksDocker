package web

import (
	database "chatHTTP/mysql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/gorilla/websocket"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// var tokenAuth *jwtauth.JWTAuth

// const Secret = "42a00d84-9914-4a77-b6bd-d2a9d09c6795"

type Handler struct {
	*chi.Mux
	*database.Store
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

func NewHandler(store *database.Store) *Handler {
	handler := &Handler{
		chi.NewRouter(),
		store,
	}

	//client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}

	handler.Use(middleware.Logger)

	handler.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true, // initialement en false mais nécessaire à true pour les httpOnly cookie dans les credentials include des requête du front
		MaxAge:           300,  // Maximum value not ignored by any of major browsers
	}))

	handler.Post("/auth/register", handler.RegisterHandler)
	handler.Post("/auth/logged", handler.LoginHandler())
	handler.Get("/user-list", handler.GetUsers())
	handler.Get("/delete-user/{id}", handler.DeleteUser())
	// Il faut encore déplacer les fonction qui sont dans pakage main actuellement dans des handler

	//handler.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
	//	ServeWs(client.hub, w, r)
	//})

	handler.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
	})

	return handler
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to upgrade connection to WebSocket", http.StatusInternalServerError)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

func (h *Handler) jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Log encoding error
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
