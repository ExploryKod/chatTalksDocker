package main

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) GetUserDiscussionsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		QueryId := chi.URLParam(r, "id")
		id, _ := strconv.Atoi(QueryId)
		discussions, err := h.Store.GetUserDiscussions(id)
		if err != nil {
			log.Fatal(err)
			h.jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"message": "Internal Server Error",
			})
			return
		}

		h.jsonResponse(w, http.StatusOK, discussions)
	}
}

func (h *Handler) DeleteUserRoomHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		QueryUserId := chi.URLParam(r, "userid")
		userid, _ := strconv.Atoi(QueryUserId)
		QueryRoomId := chi.URLParam(r, "roomid")
		roomid, _ := strconv.Atoi(QueryRoomId)

		err := h.Store.DeleteUserRoom(userid, roomid)
		if err != nil {
			log.Fatal(err)
			h.jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"message": "Internal Server Error",
			})
			return
		}
		h.jsonResponse(w, http.StatusOK, map[string]interface{}{
			"message": "User room deleted",
		})
	}
}
