// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-sql-driver/mysql"
	"log"
	"net/http"
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
	conf := mysql.Config{
		User:                 "u6ncknqjamhqpa3d",
		Passwd:               "O1Bo5YwBLl31ua5agKoq",
		Net:                  "tcp",
		Addr:                 "bnouoawh6epgx2ipx4hl-mysql.services.clever-cloud.com:3306",
		DBName:               "bnouoawh6epgx2ipx4hl",
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
	hub := newHub()
	go hub.run()

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

	// Define your routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	r.Get("/home", serveHome)

	r.Get("/mychat", serveChatPage)

	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	r.Post("/auth/register", handler.RegisterHandler)

	r.Post("/auth/logged", handler.LoginHandler())
	r.Group(func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(jwtauth.Verifier(tokenAuth))

		r.Use(jwtauth.Authenticator)
		r.Get("/chat", func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			w.Write([]byte(fmt.Sprintf("protected area. hi %v", claims["user_id"])))
		})

		r.Get("/chat", func(w http.ResponseWriter, r *http.Request) {

		})

		server := &http.Server{
			Addr:              ":8000", // Replace with your desired address
			ReadHeaderTimeout: 3 * time.Second,
			Handler:           r, // Use the chi router as the handler
		}

		err = server.ListenAndServe()
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	})
}
