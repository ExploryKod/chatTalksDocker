package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
	"strconv"
)

//func (h *Handler) CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
//	// Extract registration data
//	name := r.FormValue("name")
//
//	roomID, err := h.Store.AddRoom(RoomItem{Name: name})
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	// Respond with a success message
//	h.jsonResponse(w, http.StatusOK, map[string]interface{}{"message": "Room created", "roomID": roomID})
//}
//
//func (h *Handler) GetRoomHandler(w http.ResponseWriter, r *http.Request) {
//	// Extract registration data
//	name := r.FormValue("name")
//
//	room, err := h.Store.GetRoomByName(name)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	// Respond with a success message
//	h.jsonResponse(w, http.StatusOK, map[string]interface{}{"message": "Room found", "room": room})
//
//}
//
//func (h *Handler) GetRoomByIdHandler(w http.ResponseWriter, r *http.Request) {
//	// Extract registration data
//	id := r.FormValue("id")
//
//	room, err := h.Store.GetRoomById(id)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	// Respond with a success message
//	h.jsonResponse(w, http.StatusOK, map[string]interface{}{"message": "Room found", "room": room})
//}

func (h *Handler) JoinRoomHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomID := chi.URLParam(r, "id")
		var id, err = strconv.Atoi(roomID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		room, err := h.Store.GetRoomById(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, claims, _ := jwtauth.FromContext(r.Context())
		if username, ok := claims["username"].(string); ok {
			user, err := h.Store.GetUserByUsername(username)
			if err != nil {
				// Handle database error
				h.jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
					"message": "Internal Server Error",
				})
				return
			}
			fromRoom, err := h.GetOneUserFromRoom(room.ID, user.ID)
			if err != nil {
				h.jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
					"message": "Internal Server Error DB",
				})
				return
			}
			if fromRoom.Username != "" {
				h.jsonResponse(w, http.StatusOK, map[string]interface{}{"message": "Hi " + username + "Welcome back in your room"})
				return
			}
			err = h.Store.AddUserToRoom(room.ID, user.ID)
			if err != nil {
				return
			}
			h.jsonResponse(w, http.StatusOK, map[string]interface{}{"message": "you joined the room " + room.Name})
		} else {
			h.jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{"error": "Unauthorized"})
		}
	}
}

func (h *Handler) GetRooms() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rooms, err := h.Store.GetRooms()
		if err != nil {
			// Handle database error
			h.jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"message": "Internal Server Error",
			})
			return
		}

		h.jsonResponse(w, http.StatusOK, rooms)
	}
}

func (h *Handler) CreateRoomHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())
		if username, ok := claims["username"].(string); ok {
			user, err := h.Store.GetUserByUsername(username)
			if err != nil {
				// Handle database error
				h.jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
					"message": "Internal Server Error",
				})
				return
			}
			roomName := r.FormValue("roomName")
			roomId, err := h.Store.AddRoom(RoomItem{Name: roomName, Description: "coquelicots"})
			if err != nil {
				// Handle database error
				h.jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
					"message": "Internal Server Error",
				})
				return
			}

			h.jsonResponse(w, http.StatusOK, map[string]interface{}{"message": "Welcome " + username, "roomID": roomId})
		} else {
			h.jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{"error": "Unauthorized"})
		}
	}
}
