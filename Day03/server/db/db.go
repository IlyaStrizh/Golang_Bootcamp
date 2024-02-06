package db

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type Place struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	Location struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"location"`
}

type ElasticsearchStore struct {
	es *elasticsearch.Client
}

// NewElasticsearchStore создает новый экземпляр ElasticsearchStore
func NewElasticsearchStore() *ElasticsearchStore {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create Elasticsearch client: %s", err.Error())
	}
	return &ElasticsearchStore{es}
}

// GetPlaces возвращает список мест из хранилища с использованием пагинации
func (s *ElasticsearchStore) GetPlaces(limit, offset int) ([]Place, int, error) {
	query := map[string]interface{}{
		"from": offset,
		"size": limit,
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}

	result, err := s.getRequest(query)
	if err != nil {
		return nil, 0, err
	}

	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	total := int(result["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))

	return makePlaces(hits), total, nil
}

// GetPlacesRecommend возвращает рекомендованные места на основе координат
func (s *ElasticsearchStore) GetPlacesRecommend(lat, lon float64) ([]Place, error) {
	// Добавляем сортировку по расстоянию
	sort := map[string]interface{}{
		"size": 3,
		"sort": map[string]interface{}{
			"_geo_distance": map[string]interface{}{
				"location": map[string]interface{}{
					"lat": lat,
					"lon": lon,
				},
				"order":           "asc",
				"unit":            "km",
				"mode":            "min",
				"distance_type":   "arc",
				"ignore_unmapped": true,
			},
		},
	}

	result, err := s.getRequest(sort)
	if err != nil {
		return nil, err
	}

	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})

	return makePlaces(hits), nil
}

// getRequest выполняет запрос Elasticsearch и возвращает результат
func (s *ElasticsearchStore) getRequest(query map[string]interface{}) (map[string]interface{}, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	// Создание и отправка запроса Elasticsearch
	req := esapi.SearchRequest{
		Index:          []string{"places"},
		Body:           &buf,
		TrackTotalHits: true,
	}

	res, err := req.Do(context.Background(), s.es)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("Elasticsearch search error: %s", res.String())
	}

	// Обработка ответа Elasticsearch и возврат результата
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// makePlaces создает срез мест на основе результатов Elasticsearch
func makePlaces(hits []interface{}) []Place {
	places := make([]Place, 0)
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"]
		placeBytes, err := json.Marshal(source)
		if err != nil {
			continue
		}

		var place Place
		if err := json.Unmarshal(placeBytes, &place); err != nil {
			continue
		}

		places = append(places, place)
	}

	return places
}
