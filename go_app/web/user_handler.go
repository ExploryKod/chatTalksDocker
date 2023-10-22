package web

import (
	"chatHTTP"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/jwtauth/v5"
)

type TemplateData struct {
	Titre   string
	Content any
}

//func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
//	// Extract registration data
//	username := r.FormValue("username")
//	password := r.FormValue("password")
//
//	userID, err := h.Store.AddUser(chatHTTP.UserItem{Username: username, Password: password})
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	// Respond with a success message
//	h.jsonResponse(w, http.StatusOK, map[string]interface{}{"message": "Registration successful", "userID": userID})
//}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		h.jsonResponse(w, http.StatusBadRequest, map[string]interface{}{"error": "Username and password are required"})
		return
	}

	if _, err := h.Store.GetUserByUsername(username); err == nil {
		h.jsonResponse(w, http.StatusConflict, map[string]interface{}{"error": "Username already taken"})
		return
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	userID, err := h.Store.AddUser(chatHTTP.UserItem{Username: username, Password: hashedPassword})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.jsonResponse(w, http.StatusOK, map[string]interface{}{"message": "Registration successful", "userID": userID})
}

func hashPassword(password string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

var tokenAuth *jwtauth.JWTAuth

func GenerateRandomSecret(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func init() {
	secret, err := GenerateRandomSecret(64)
	if err != nil {
		log.Fatal("Error generating random secret:", err)
	}
	tokenAuth = jwtauth.New("HS256", []byte(secret), nil)
}

func MakeToken(name string) string {
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"username": name})
	return tokenString
}

func loginJWTHandler(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Perform authentication (e.g., check credentials against a database)
	if isValidCredentials(username, password) {
		// Generate a JWT
		_, tokenString, _ := tokenAuth.Encode(jwt.MapClaims{"username": username, "exp": time.Now().Add(time.Hour).Unix()})

		// Respond with the JWT
		response := map[string]string{"token": tokenString}
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

func isValidCredentials(username, password string) bool {
	// Implement your authentication logic here (e.g., check against a database)
	// Return true if the credentials are valid, otherwise false.
	return true
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, token, err := jwtauth.FromContext(r.Context())
		if err != nil || token == nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Vous pouvez extraire des informations du token si nécessaire
		// username := token.Claims.(jwt.MapClaims)["username"].(string)
		// ...

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract username and password from the request body or form data
		username := r.FormValue("username")
		password := r.FormValue("password")

		user, err := h.Store.GetUserByUsername(username)
		if err != nil {
			h.jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"message": "Nom d'utilisateur incorrect",
			})
			return
		}

		if user.Username == "" {
			h.jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
				"message": "Vous devez entrer un nom d'utilisateur qui existe",
			})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err == nil {

			token := MakeToken(username)

			http.SetCookie(w, &http.Cookie{
				HttpOnly: true,
				Expires:  time.Now().Add(7 * 24 * time.Hour),
				//Reduce cross-site request forgery (CSRF) -only GET requests are allowed :
				SameSite: http.SameSiteLaxMode,
				//Secure: true, for https
				Name:  "jwt",
				Value: token,
			})

			// Successful login
			response := map[string]string{"message": "Vous êtes bien connecté", "redirect": "/", "token": token}
			h.jsonResponse(w, http.StatusOK, response)
		} else {
			// Passwords do not match
			h.jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
				"message": "Mot de passe incorrect",
			})
		}
	}
}

//func (h *Handler) LoginHandler() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//
//		// Extract username and password from the request body or form data
//		username := r.FormValue("username")
//		password := r.FormValue("password")
//
//		user, err := h.Store.GetUserByUsername(username)
//		if err != nil {
//			// Handle database error
//			h.jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
//				"message": "Internal Server Error",
//			})
//			return
//		}
//
//		if user.Username == "" {
//			// User not found
//			h.jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
//				"message": "Nom d'utilisateur incorrect",
//			})
//			return
//		}
//
//		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
//		if err != nil {
//			// Handle database error
//			h.jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
//				"message": "Mot de passe incorrect",
//			})
//			return
//		}
//
//		if user.Username == "" || user.Password == "" {
//			http.Error(w, "Il reste des champs vides", http.StatusBadRequest)
//			return
//		}
//
//		// Check if the user exists and the password matches
//		if user.Username == username && user.Password == password {
//			token := MakeToken(username)
//
//			http.SetCookie(w, &http.Cookie{
//				HttpOnly: true,
//				Expires:  time.Now().Add(7 * 24 * time.Hour),
//				SameSite: http.SameSiteLaxMode,
//				// Uncomment below for HTTPS:
//				// Secure: true,
//				Name:  "jwt", // Must be named "jwt" or else the token cannot be searched for by jwtauth.Verifier.
//				Value: token,
//			})
//			// Successful login
//
//			response := map[string]string{"message": "Vous êtes bien connecté", "redirect": "/", "token": token}
//			h.jsonResponse(w, http.StatusOK, response)
//		} else if user.Password != password {
//			// Failed login
//			h.jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
//				"message": "Mot de passe incorrect",
//			})
//		} else if user.Username != username {
//			// Failed login
//			h.jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
//				"message": "Nom d'utilisateur incorrect",
//			})
//		} else {
//			// Failed login
//			h.jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
//				"message": "Nom d'utilisateur et mot de passe incorrects",
//			})
//		}
//	}
//}

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
//		h.jsonResponse(w, http.StatusOK, users)
//	}
//}

func (h *Handler) GetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Accédez aux informations du token JWT
		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			// Gérez l'erreur, par exemple, si le token n'est pas présent ou invalide
			h.jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
				"message": "Token invalide",
			})
			return
		}

		// Vous pouvez maintenant accéder aux claims, y compris le nom d'utilisateur
		if username, ok := claims["username"].(string); ok {
			// Utilisez le nom d'utilisateur dans votre logique métier
			user, err := h.Store.GetUserByUsername(username)
			if err != nil {
				// Gérez l'erreur de base de données
				h.jsonResponse(w, http.StatusInternalServerError, map[string]interface{}{
					"message": "Erreur de base de données",
				})
				return
			}

			// Utilisez les informations sur l'utilisateur dans la réponse JSON
			h.jsonResponse(w, http.StatusOK, map[string]interface{}{
				"message":  "Liste des utilisateurs",
				"username": username,
				"user":     user,
			})
		} else {
			h.jsonResponse(w, http.StatusUnauthorized, map[string]interface{}{
				"message": "Claims invalides",
			})
		}
	}
}

func (h *Handler) DeleteUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		QueryId := chi.URLParam(request, "id")
		id, _ := strconv.Atoi(QueryId)

		err := h.Store.DeleteUser(id)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(writer, request, "/", http.StatusSeeOther)
	}
}
