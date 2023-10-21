package main

import (
	database "chatHTTP/mysql"
	"chatHTTP/web"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-sql-driver/mysql"
)

func main() {
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

	store := database.CreateStore(db)
	mux := web.NewHandler(store)

	err = http.ListenAndServe(":8000", mux)
	if err != nil {
		_ = fmt.Errorf("impossible de lancer le serveur : %w", err)
		return
	}
}
