package main

import (
	"net/http"
	"strconv"
)

func (h *Handler) CreateMessageHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	content := r.FormValue("content")

	sender, _ := h.Store.GetUserByUsername(username)

	roomID := r.FormValue("roomID")
	roomIDInt, _ := strconv.Atoi(roomID)

	if sender.Username != "" {
		messageID, err := h.Store.AddMessage(MessageItem{Content: content, UserID: sender.ID, RoomID: roomIDInt, Username: sender.Username})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		h.jsonResponse(w, http.StatusOK, map[string]interface{}{"message": "Message sent", "roomID": roomIDInt, "messageID": messageID, "userID": sender.ID})
	} else {
		println("Noooo")
		http.Error(w, "No user with this id found", http.StatusBadRequest)
		return
	}
}
