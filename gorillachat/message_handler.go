package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *Handler) CreateMessageHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	content := r.FormValue("content")
	//maxMessageTableRows := 2

	sender, _ := h.Store.GetUserByUsername(username)
	//messagesNumber, _ := h.Store.CountMessagesSent()

	roomID := r.FormValue("roomID")
	roomIDInt, _ := strconv.Atoi(roomID)

	if sender.Username != "" && content != "" && roomID != "" {
		// TODO : add a limit of messages per room or accept all messages
		//messageID, err := h.Store.AddMessage(MessageItem{Content: content, UserID: sender.ID, RoomID: roomIDInt, Username: sender.Username})
		//if err != nil {
		//	http.Error(w, err.Error(), http.StatusInternalServerError)
		//	return
		//}
		messageID := 0
		h.jsonResponse(w, http.StatusOK, map[string]interface{}{"message": "Message sent", "roomID": roomIDInt, "messageID": messageID, "userID": sender.ID})
	} else if sender.Username == "" {
		http.Error(w, "No user with this id found", http.StatusBadRequest)
		return
	} else {
		http.Error(w, "Requête non-satisfaite", http.StatusBadRequest)
		return
	}
}

func (h *Handler) GetMessageHandler(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "id")

	var id, err = strconv.Atoi(roomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	messages, err := h.Store.GetMessagesFromRoom(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.jsonResponse(w, http.StatusOK, map[string]interface{}{"message": "Messages found", "messages": messages})
}

func (h *Handler) DeleteMessageFromRoomHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		QueryId := chi.URLParam(request, "id")
		id, _ := strconv.Atoi(QueryId)

		err := h.Store.DeleteMessagesByRoomId(id)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			h.jsonResponse(writer, http.StatusInternalServerError, map[string]interface{}{"message": err.Error()})
			return
		}
		h.jsonResponse(writer, http.StatusOK, map[string]interface{}{"message": "L'historique des messages de cette salle a été supprimé"})
	}
}
