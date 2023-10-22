package main

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
