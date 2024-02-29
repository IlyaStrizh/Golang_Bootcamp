package handle

import (
	"bufio"
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/russross/blackfriday"
)

const (
	adminLogin    = "login"
	adminPassword = "password"
)

type Article struct {
	ID      int
	Title   string
	Content string
}

type page struct {
	FirstPage   int
	LastPage    int
	NextPage    int
	PrevPage    int
	CurrentPage int
	TotalPages  int
	Articles    []Article
}

func HandleMainPage(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	tmpl, err := template.ParseFiles("../html/mainPage.html")
	if err != nil {
		log.Printf("Failed to parse template file: %s", err.Error())
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}

	currentPage, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		currentPage = 1
	}

	articles, totalArticles, err := getArticles(db, currentPage)
	if err != nil {
		log.Println("Failed to fetch articles from the database: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	totalPages := (totalArticles + 2) / 3
	homePage := page{
		FirstPage:   1,
		LastPage:    totalPages,
		NextPage:    min(currentPage+1, totalPages),
		PrevPage:    max(currentPage-1, 1),
		CurrentPage: currentPage,
		TotalPages:  totalPages,
		Articles:    articles,
	}
	err = tmpl.Execute(w, homePage)
	if err != nil {
		log.Printf("Failed to execute template: %s", err.Error())
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}
}

func getArticles(db *sql.DB, page int) ([]Article, int, error) {
	const LIMIT = 3
	offset := (page - 1) * LIMIT

	rows, err := db.Query("SELECT id, title, SUBSTRING(content, 1, 210) as content FROM articles LIMIT $1 OFFSET $2", LIMIT, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var article Article
		err := rows.Scan(&article.ID, &article.Title, &article.Content)
		if err != nil {
			return nil, 0, err
		}
		articles = append(articles, article)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	var count int
	err = db.QueryRow("SELECT count(*) FROM articles").Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return articles, count, nil
}

func HandleArticle(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	tmpl, err := template.ParseFiles("../html/article.html")
	if err != nil {
		log.Printf("Failed to parse template file: %s", err.Error())
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	articleID, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		http.NotFound(w, r)
		return
	}

	article, err := getArticleByID(db, articleID)
	if err != nil {
		log.Println("Failed to fetch article from the database: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	article.Content = string(blackfriday.MarkdownCommon([]byte(article.Content)))

	err = tmpl.Execute(w, article)
	if err != nil {
		log.Printf("Failed to execute template: %s", err.Error())
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}
}

func getArticleByID(db *sql.DB, id int) (*Article, error) {
	var article Article

	err := db.QueryRow("SELECT id, title, content FROM articles WHERE id = $1", id).Scan(&article.ID, &article.Title, &article.Content)
	if err != nil {
		return nil, err
	}

	return &article, nil
}

func HandleAdminPanel(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	login := r.FormValue("login")
	password := r.FormValue("password")

	if login != "" && password != "" {
		if login != adminLogin || password != adminPassword {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		postArticle(w, r, db)
	}

	serve(w, r, "../html/admin.html")
}

func postArticle(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		path := r.FormValue("path")

		Insert(title, path, db)
		// http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	serve(w, r, "../html/newPost.html")
}

func serve(w http.ResponseWriter, r *http.Request, filePath string) {
	adminFile, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer adminFile.Close()
	// Устанавливаем MIME-тип для файла HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Отправляем файл в ответе HTTP
	http.ServeFile(w, r, filePath)
}

func Insert(title string, path string, db *sql.DB) {
	var file *os.File
	if ext := filepath.Ext(path); ext == ".md" {
		_, err := os.Stat(path)
		if err != nil {
			log.Println("File is not exsist:", err)
			return
		} else {
			file, err = os.Open(path)
			if err != nil {
				log.Println("Error opening the file:", err)
				return
			}
			defer file.Close()
		}
		parseMD(title, file, db)
	} else {
		log.Println("File is not .md")
	}
}

func parseMD(title string, file *os.File, db *sql.DB) {
	var content string

	scanner := bufio.NewScanner(file)
	if title == "" && scanner.Scan() {
		title = strings.TrimPrefix(scanner.Text(), "### ")
	}

	if scanner.Scan() {
		text := scanner.Text() + "\n"
		if !strings.HasPrefix(text, "### ") {
			content += text
		}
	}

	for scanner.Scan() {
		content += scanner.Text() + "\n"
	}

	insertArticle(title, content, db)
}

func insertArticle(title, content string, db *sql.DB) error {
	_, err := db.Exec("INSERT INTO articles (title, content) VALUES ($1, $2)", title, content)
	if err != nil {
		return err
	}

	return nil
}
