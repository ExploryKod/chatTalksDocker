package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
)

type TemplateData struct {
	Titre   string
	Content any
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Extract registration data
	username := r.FormValue("username")
	password := r.FormValue("password")

	userID, err := h.Store.AddUser(UserItem{Username: username, Password: password})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	h.jsonResponse(w, http.StatusOK, map[string]interface{}{"message": "Registration successful", "userID": userID})
}

var tokenAuth *jwtauth.JWTAuth

const Secret = "mysecretamaury"

func init() {
	tokenAuth = jwtauth.New("HS256", []byte(Secret), nil)
}

func MakeToken(name string) string {
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"username": name})
	return tokenString
}

//func loginJWTHandler(w http.ResponseWriter, r *http.Request) {
//	username, password, ok := r.BasicAuth()
//	if !ok {
//		http.Error(w, "Unauthorized", http.StatusUnauthorized)
//		return
//	}
//
//	// Perform authentication (e.g., check credentials against a database)
//	if isValidCredentials(username, password) {
//		// Generate a JWT
//		_, tokenString, _ := tokenAuth.Encode(jwt.MapClaims{"username": username, "exp": time.Now().Add(time.Hour).Unix()})
//
//		// Respond with the JWT
//		response := map[string]string{"token": tokenString}
//		json.NewEncoder(w).Encode(response)
//	} else {
//		http.Error(w, "Unauthorized", http.StatusUnauthorized)
//	}
//}
//
//func isValidCredentials(username, password string) bool {
//	// Implement your authentication logic here (e.g., check against a database)
//	// Return true if the credentials are valid, otherwise false.
//	return true
//}

func (h *Handler) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Extract username and password from the request body or form data
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Validate user credentials against the database
		user, err := h.Store.GetUserByUsername(username)
		if err != nil {
			// Handle database error
			h.jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"message": "Internal Server Error",
			})
			return
		}

		if user.Username == "" || user.Password == "" {
			http.Error(w, "Il reste des champs vide", http.StatusBadRequest)
			return
		}

		// Check if the user exists and the password matches
		if user.Username == username && user.Password == password {
			token := MakeToken(username)

			http.SetCookie(w, &http.Cookie{
				HttpOnly: true,
				Expires:  time.Now().Add(7 * 24 * time.Hour),
				SameSite: http.SameSiteLaxMode,
				// Uncomment below for HTTPS:
				// Secure: true,
				Name:  "jwt", // Must be named "jwt" or else the token cannot be searched for by jwtauth.Verifier.
				Value: token,
			})
			// Successful login
			client = newClient(
				upgrader.Upgrade(w, r, nil),
				NewWebsocketServer(),
				username,
			)
				})
			response := map[string]string{"message": "Vous êtes bien connecté", "redirect": "/", "token": token}
			h.jsonResponse(w, http.StatusOK, response)
		} else if user.Password != password {
			// Failed login
			h.jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
				"message": "Mot de passe incorrect",
			})
		} else if user.Username != username {
			// Failed login
			h.jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
				"message": "Nom d'utilisateur incorrect",
			})
		} else {
			// Failed login
			h.jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
				"message": "Nom d'utilisateur et mot de passe incorrects",
			})
		}
	}
}

//func (h *Handler) GetUsers() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		users, err := h.Store.GetUsers()
//		if err != nil {
//			// Handle database error
//			h.jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
//				"message": "Internal Server Error",
//			})
//			return
//		}
//
//		// Respond with the users in JSON format
//		h.jsonResponse(w, http.StatusOK, users)
//	}
//}

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
			roomId, err := h.Store.AddRoom(RoomItem{Name: roomName, Description: "room de " + user.Username})
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
