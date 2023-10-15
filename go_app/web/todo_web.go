package web

import (
	"demoHTTP"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

type TemplateData struct {
	Titre   string
	Content any
}

func (h *Handler) WebShowTodos() http.HandlerFunc {
	// Placer cette déclaration avant de retourner le handler
	// permet de ne créer qu'une seule fois cette struct
	// plutôt que de la créer à chaque requête
	data := TemplateData{Titre: "Tous les todos"}

	return func(writer http.ResponseWriter, request *http.Request) {
		todos, err := h.Store.GetTodos()
		data.Content = todos

		// ParseFS fonctionne exactement comme ParseFiles mais va chercher
		// dans un fileSystem donné plutôt que dans celui de l'hôte
		tmpl, err := template.ParseFS(demoHTTP.EmbedTemplates, "templates/layout.gohtml", "templates/list.gohtml")
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}

		// Je passe mes données ici
		err = tmpl.ExecuteTemplate(writer, "layout", data)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (h *Handler) WebCreateTodoForm() http.HandlerFunc {
	data := TemplateData{Titre: "Add a todo"}

	return func(writer http.ResponseWriter, request *http.Request) {
		tmpl, err := template.ParseFS(demoHTTP.EmbedTemplates, "templates/layout.gohtml", "templates/form.gohtml")
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}

		err = tmpl.ExecuteTemplate(writer, "layout", data)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (h *Handler) WebAddTodo() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := request.ParseForm()
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}

		_, err = h.Store.AddTodo(demoHTTP.TodoItem{Title: request.FormValue("new-todo")})
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
		http.Redirect(writer, request, "/", http.StatusSeeOther)
	}
}

func (h *Handler) WebToogleTodo() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		QueryId := chi.URLParam(request, "id")
		id, _ := strconv.Atoi(QueryId)

		err := h.Store.ToggleTodo(id)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(writer, request, "/", http.StatusSeeOther)
	}
}

func (h *Handler) WebDeleteTodo() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		QueryId := chi.URLParam(request, "id")
		id, _ := strconv.Atoi(QueryId)

		err := h.Store.DeleteTodo(id)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(writer, request, "/", http.StatusSeeOther)
	}
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Extract registration data
	username := r.FormValue("username")
	password := r.FormValue("password")

	userID, err := h.Store.AddUser(demoHTTP.UserItem{Username: username, Password: password})
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

			response := map[string]string{"message": "Vous êtes bien connecté", "redirect": "/profile"}
			jsonResponse(w, http.StatusOK, response)
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
