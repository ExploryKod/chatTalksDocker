package main

import (
	"database/sql"
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
	"time"
)

type Handler struct {
	*chi.Mux
	*Store
}

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/home" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func serveChatPage(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Perform a redirect to /chat
	http.Redirect(w, r, "/chat", http.StatusSeeOther)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // Default to port 8000 if PORT environment variable is not set
	}

	//conf := mysql.Config{
	//	User:                 "u6ncknqjamhqpa3d",
	//	Passwd:               "O1Bo5YwBLl31ua5agKoq",
	//	Net:                  "tcp",
	//	Addr:                 "bnouoawh6epgx2ipx4hl-mysql.services.clever-cloud.com:3306",
	//	DBName:               "bnouoawh6epgx2ipx4hl",
	//	AllowNativePasswords: true,
	//}

	conf := mysql.Config{
		User:                 "root",
		Passwd:               os.Getenv("MARIADB_ROOT_PASSWORD"),
		Net:                  "tcp",
		Addr:                 "database:3306",
		DBName:               os.Getenv("MARIADB_DATABASE"),
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", conf.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	store := CreateStore(db)
	//mux := NewHandler(store)

	handler := &Handler{
		chi.NewRouter(),
		store,
	}

	flag.Parse()
	wsServer := NewWebsocketServer()
	go wsServer.Run()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true, // initialement en false
		MaxAge:           300,  // Maximum value not ignored by any of major browsers
	}))

	r.Post("/auth/register", handler.RegisterHandler)
	r.Post("/auth/logged", handler.LoginHandler())

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))

		r.Use(jwtauth.Authenticator)
		// use JoinHub method to join a hub
		r.Get("/chat/{id}", handler.JoinRoomHandler())
		r.Get("/chat/rooms", handler.GetRooms())
		r.Post("/chat/create", handler.CreateRoomHandler())
		r.Get("/user-list", handler.GetUsers())
		r.Delete("/delete-user/{id}", handler.DeleteUserHandler())
		r.Get("/update-user", handler.UpdateHandler)
	})
	// Define your routes

	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(wsServer, w, r)
	})

	server := &http.Server{
		Addr:              port, // Replace with your desired address
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           r, // Use the chi router as the handler
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
