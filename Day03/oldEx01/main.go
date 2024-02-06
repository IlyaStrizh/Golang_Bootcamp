// curl -X PUT "localhost:9200/places/_settings" -H 'Content-Type: application/json' -d'
//
//	{
//	    "index" : {
//	        "max_result_window" : 20000
//	    }
//	}
//
// '

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type Place struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

type Store interface {
	GetPlaces(limit int, offset int) ([]Place, int, error)
}

type ElasticsearchStore struct {
	es *elasticsearch.Client
}

func NewElasticsearchStore() *ElasticsearchStore {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Ошибка создания Elasticsearch client: %s", err.Error())
	}
	return &ElasticsearchStore{es}
}

func (s *ElasticsearchStore) GetPlaces(limit int, offset int) ([]Place, int, error) {
	query := map[string]interface{}{
		"from": offset,
		"size": limit,
		//	"query": map[string]interface{}{
		//		"match_all": map[string]interface{}{},
		//	},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, 0, err
	}

	req := esapi.SearchRequest{
		Index:          []string{"places"},
		Body:           &buf,
		TrackTotalHits: true,
	}

	res, err := req.Do(context.Background(), s.es)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, 0, fmt.Errorf("failed to search: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, 0, err
	}

	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	total := int(result["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))

	places := make([]Place, 0)
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		place := Place{
			Name:    source["name"].(string),
			Address: source["address"].(string),
			Phone:   source["phone"].(string),
		}
		places = append(places, place)
	}

	return places, total, nil
}

func renderPage(w http.ResponseWriter, places []Place, total int, currentPage int) {
	totalPages := (total-1)/10 + 1

	// Создание массива для номеров страниц
	pageNums := make([]int, totalPages)
	for i := 0; i < totalPages; i++ {
		pageNums[i] = i + 1
	}

	pageData := struct {
		TotalPages  []int
		CurrentPage int
		Places      []Place
	}{
		TotalPages:  pageNums,
		CurrentPage: currentPage,
		Places:      places,
	}

	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, pageData)
	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}
}

func parsePageQueryParam(r *http.Request) (int, error) {
	pageStr := r.FormValue("page")
	if pageStr == "" {
		return 1, nil
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return 0, fmt.Errorf("Invalid 'page' value: '%s'", pageStr)
	}

	return page, nil
}

func handlePlaces(w http.ResponseWriter, r *http.Request, store Store) {
	page, err := parsePageQueryParam(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	limit := 10
	offset := (page - 1) * limit

	places, total, err := store.GetPlaces(limit, offset)
	if err != nil {
		http.Error(w, "Error retrieving places", http.StatusInternalServerError)
		return
	}

	renderPage(w, places, total, page)
}

func main() {
	store := NewElasticsearchStore()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlePlaces(w, r, store)
	})
	log.Println("Server started on port 8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}
