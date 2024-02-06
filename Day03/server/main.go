// curl -X PUT "localhost:9200/places/_settings" -H 'Content-Type: application/json' -d'
// {
//	 "index" : {
//	   "max_result_window" : 20000
//	 }
// }'
// http://127.0.0.1:8888/?page=1
// http://127.0.0.1:8888/api/places?page=3
// http://127.0.0.1:8888/api/recommend?lat=55.674&lon=37.666
// http://127.0.0.1:8888/api/get_token
// curl -X GET "http://127.0.0.1:8888/api/recommend?lat=55.674&lon=37.666" -H "Authorization: Bearer <token>"

package main

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/golang-jwt/jwt/v5"

	types "Day03/Server/db"
)

type Store interface {
	GetPlaces(limit, offset int) ([]types.Place, int, error)
	GetPlacesRecommend(lat, lon float64) ([]types.Place, error)
}

func main() {
	secret := make([]byte, 32)  // генерируем срез длиной 32 байта
	_, err := rand.Read(secret) // заполняем срез случайными значениями
	if err != nil {
		fmt.Println("Failed to generate secret key:", err)
		return
	}

	store := types.NewElasticsearchStore()

	pwd, _ := os.Getwd()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(pwd+"/server/template"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleHome(w, r, store)
	})
	http.HandleFunc("/api/places", func(w http.ResponseWriter, r *http.Request) {
		handlePlaces(w, r, store)
	})
	http.HandleFunc("/api/recommend", func(w http.ResponseWriter, r *http.Request) {
		handlePlacesRecommend(w, r, store, secret)
	})
	http.HandleFunc("/api/get_token", func(w http.ResponseWriter, r *http.Request) {
		getToken(w, r, secret)
	})

	port := 8888

	log.Printf("Server started on port %d", port)
	log.Fatal(http.ListenAndServe("127.0.0.1:8888", nil))
}

// handleHome обрабатывает запрос для получения списка мест с использованием пагинации для HTML
func handleHome(w http.ResponseWriter, r *http.Request, s Store) {
	page, err := parsePageQueryParam(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	limit := 10
	offset := page * limit

	places, total, err := s.GetPlaces(limit, offset)
	if err != nil {
		log.Printf("Failed to get places from store: %s", err.Error())
		http.Error(w, "Error getting places", http.StatusInternalServerError)
		return
	}

	maxPage := (total+limit-1)/limit - 1
	if page > maxPage {
		errorMessage := fmt.Sprintf("Invalid 'page' value: '%d'", page)
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

	renderPageHTML(w, places, total, page)
}

// handlePlaces обрабатывает запрос для получения списка мест с использованием пагинации для JSON
func handlePlaces(w http.ResponseWriter, r *http.Request, s Store) {
	page, err := parsePageQueryParam(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	limit := 10
	offset := page * limit

	places, total, err := s.GetPlaces(limit, offset)
	if err != nil {
		log.Printf("Failed to get places from store: %s", err.Error())
		http.Error(w, "Error getting places", http.StatusInternalServerError)
		return
	}

	maxPage := (total+limit-1)/limit - 1
	if page > maxPage {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Invalid 'page' value: '%d'", page),
		})
		return
	}

	renderPage(w, places, total, page)
}

// parsePageQueryParam разбирает значение параметра "page" из URL и возвращает его в виде целого числа
func parsePageQueryParam(r *http.Request) (int, error) {
	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		return 0, nil
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 0 {
		return 0, fmt.Errorf("Invalid 'page' value: '%s'", pageStr)
	}

	return page, nil
}

// handlePlacesRecommend обрабатывает запрос для получения списка мест с использованием пагинации для JSON Recommend
func handlePlacesRecommend(w http.ResponseWriter, r *http.Request, s Store, secret []byte) {
	str := r.Header.Get("Authorization")
	token := strings.Split(str, " ")
	if len(token) != 2 {
		http.Error(w, "Empty jwt token", http.StatusUnauthorized)
		return
	}
	if err := parseToken(token[1], secret); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	latStr := r.URL.Query().Get("lat") // Получаем текущую широту из запроса
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, "Invalid 'lat' number", http.StatusBadRequest)
		return
	}

	lonStr := r.URL.Query().Get("lon") // Получаем текущую долготу из запроса
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		http.Error(w, "Invalid 'lon' number", http.StatusBadRequest)
		return
	}

	places, err := s.GetPlacesRecommend(lat, lon)
	if err != nil {
		log.Printf("Failed to get places from store: %s", err.Error())
		http.Error(w, "Error getting places", http.StatusInternalServerError)
		return
	}

	res := struct {
		Name   string        `json:"name"`
		Places []types.Place `json:"places"`
	}{
		Name:   "Recommendation",
		Places: places,
	}

	data, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		log.Println(err)
		return
	}

	writeHeader(w, data)
}

// writeHeader пишет ответ на запросы в JSON формате и ответ 200
func writeHeader(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// renderPageHTML возвращает данные в HTML
func renderPageHTML(w http.ResponseWriter, places []types.Place, total, currentPage int) {
	totalPages := (total + 9) / 10

	pageData := struct {
		FirstPage   int
		LastPage    int
		NextPage    int
		PrevPage    int
		CurrentPage int
		TotalPages  int
		Places      []types.Place
	}{
		FirstPage:   0,
		LastPage:    totalPages - 1,
		NextPage:    min(currentPage+1, totalPages-1),
		PrevPage:    max(currentPage-1, 0),
		CurrentPage: currentPage,
		TotalPages:  totalPages - 1,
		Places:      places,
	}

	tmpl, err := template.ParseFiles("server/template/index.html")
	if err != nil {
		log.Printf("Failed to parse template file: %s", err.Error())
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, pageData)
	if err != nil {
		log.Printf("Failed to execute template: %s", err.Error())
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}
}

// renderPage возвращает данные в JSON
func renderPage(w http.ResponseWriter, places []types.Place, total, currentPage int) {
	totalPages := (total + 9) / 10

	pageData := struct {
		Name     string
		Total    int
		Places   []types.Place
		PrevPage int
		NextPage int
		LastPage int
	}{
		Name:     "places",
		Total:    total,
		Places:   places,
		PrevPage: max(currentPage-1, 0),
		NextPage: min(currentPage+1, totalPages-1),
		LastPage: totalPages - 1,
	}

	data, err := json.MarshalIndent(pageData, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeHeader(w, data)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// getToken обрабатывает запрос для получения токена JSON Web Tokens
func getToken(w http.ResponseWriter, r *http.Request, secret []byte) {
	token, err := generateToken(secret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}

	data, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		log.Println(err)
		http.Error(w, "Error occurred while marshaling response", http.StatusInternalServerError)
		return
	}

	writeHeader(w, data)
}

// generateToken генерирует токен JSON Web Tokens
func generateToken(secret []byte) (string, error) {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(24 * time.Hour).Unix(), // время жизни
		"iat": time.Now().Unix(),                     // timestamp
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// parseToken проверяет входящий JWT (JSON Web Token) на допустимость и валидность
func parseToken(token string, secret []byte) error {
	parser := jwt.Parser{}
	_, err := parser.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid signing method")
		}

		return secret, nil
	})

	return err
}
