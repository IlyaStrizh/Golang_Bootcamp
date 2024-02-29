package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"Day06/ex02/internal/handle"

	"golang.org/x/time/rate"
)

const (
	dbHost     = "localhost"
	dbPort     = 5051
	dbUser     = "pitermar"
	dbPassword = "1243"
	dbName     = "article"
)

func main() {
	db := initDB()
	defer db.Close()

	var limiter = rate.NewLimiter(rate.Limit(100), 50)

	// Настроить обработчик для обслуживания статических файлов из папки "css"
	cssDir := http.Dir("../css/")
	cssHandler := http.StripPrefix("/css/", http.FileServer(cssDir))
	http.Handle("/css/", cssHandler)

	// Настроить обработчик для обслуживания статических файлов из папки "images"
	imagesDir := http.Dir("../images/")
	imagesHandler := http.StripPrefix("/images/", http.FileServer(imagesDir))
	http.Handle("/images/", imagesHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if limiter.Allow() {
			handle.HandleMainPage(w, r, db)
		} else {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		}
	})
	http.HandleFunc("/article/", func(w http.ResponseWriter, r *http.Request) {
		if limiter.Allow() {
			handle.HandleArticle(w, r, db)
		} else {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		}
	})
	http.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		if limiter.Allow() {
			handle.HandleAdminPanel(w, r, db)
		} else {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		}
	})

	log.Println("Server started on port 8888")
	log.Fatal(http.ListenAndServe("127.0.0.1:8888", nil))
}

func initDB() *sql.DB {
	dbURL := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS articles")
	if err != nil {
		log.Fatal("Failed to drop database: ", err)
	}

	tableQuery := `
						CREATE TABLE IF NOT EXISTS articles (
						    id SERIAL PRIMARY KEY,
						    title VARCHAR(100),
						    content TEXT,
						    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
						)`
	_, err = db.Exec(tableQuery)
	if err != nil {
		log.Fatal("Failed to create database: ", err)
	}
	handle.Insert("", "../md/article.md", db)

	return db
}
