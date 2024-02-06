// curl -X PUT "localhost:9200/places/_settings" -H 'Content-Type: application/json' -d'
// {
//	 "index" : {
//	   "max_result_window" : 20000
//	 }
// }'
// http://localhost:9200/places/?pretty
// http://localhost:9200/places/_doc/1?pretty

// curl -X DELETE "localhost:9200/places"

package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

const (
	indexName       = "places"
	mappingFileName = "database/schema.json"
	dataCSVFileName = "../materials/data.csv"
	elasticHost     = "http://localhost:9200"
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

func main() {
	es := createIndex()
	setMappings(es)
	uploadData(es)
}

// createIndex Создает индекс в elasticsearch
func createIndex() *elasticsearch.Client {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Failed to create Elasticsearch client: %s", err)
	}
	deleteIndex(es)

	reqCreate := esapi.IndicesCreateRequest{
		Index: indexName,
	}

	createRes, err := reqCreate.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Failed to create index: %s", err)
	}
	defer createRes.Body.Close()

	if createRes.IsError() {
		log.Fatalf("Failed to create index: %s", createRes.String())
	}

	log.Println("Index successfully created")
	return es
}

// deleteIndex удаляет индекс если он уже существует
func deleteIndex(es *elasticsearch.Client) {
	req := esapi.IndicesExistsRequest{
		Index: []string{indexName},
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Failed to check if index exists: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		deleteReq := esapi.IndicesDeleteRequest{
			Index: []string{indexName},
		}

		deleteRes, err := deleteReq.Do(context.Background(), es)
		if err != nil {
			log.Fatalf("Failed to delete index: %s", err)
		}
		defer deleteRes.Body.Close()

		if deleteRes.IsError() {
			log.Fatalf("Failed to delete index: %s", deleteRes.String())
		}
		log.Println("Index successfully deleted")
	}
}

// setMappings устанавливает Mapping в созданном индексе elasticsearch
func setMappings(es *elasticsearch.Client) {
	mappingBytes, err := os.ReadFile(mappingFileName)
	if err != nil {
		log.Fatalf("Failed to read mapping file: %s", err)
	}

	req := esapi.IndicesPutMappingRequest{
		Index: []string{indexName},
		Body:  bytes.NewReader(mappingBytes),
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Failed to set mappings: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Failed to set mappings: %s", res.String())
	}

	log.Println("Mappings successfully set")
}

// uploadData загружает данные из файла в базу elasticsearch
func uploadData(es *elasticsearch.Client) {
	file, err := os.Open(dataCSVFileName)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1

	_, err = reader.Read()
	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}

	places := makePlaces(reader)

	bulkRequestBody := strings.Builder{}
	for i, place := range places {
		data, err := json.Marshal(place)
		if err != nil {
			log.Fatalf("Failed to marshal data: %s", err)
		}
		action := `{"index" : { "_index" : "` + indexName + `", "_id" : "` + strconv.Itoa(i+1) + `" } }` + "\n" // используем i+1 в качестве идентификатора документа
		bulkRequestBody.WriteString(action)
		bulkRequestBody.WriteString(string(data))
		bulkRequestBody.WriteString("\n")
	}
	req := esapi.BulkRequest{
		Body: strings.NewReader(bulkRequestBody.String()),
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Failed to upload data: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Failed to upload data: %s", res.String())
	}
	log.Println("Data successfully uploaded")
}

// makePlaces создает и заполняет слайс структур данных из файла для загрузки в базу
func makePlaces(reader *csv.Reader) []Place {
	var places []Place
	for {
		line, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("Failed to read CSV line: %s", err)
		}
		id, err := strconv.Atoi(line[0])
		if err != nil {
			log.Fatalf("Failed to parse ID: %s", err)
		}

		lat, err := strconv.ParseFloat(line[5], 64)
		if err != nil {
			log.Fatalf("Failed to parse latitude: %s", err)
		}
		lon, err := strconv.ParseFloat(line[4], 64)
		if err != nil {
			log.Fatalf("Failed to parse longitude: %s", err)
		}

		place := Place{
			ID:      id + 1,
			Name:    line[1],
			Address: line[2],
			Phone:   line[3],
			Location: struct {
				Lat float64 `json:"lat"`
				Lon float64 `json:"lon"`
			}{
				Lat: lat,
				Lon: lon,
			},
		}
		places = append(places, place)
	}

	return places
}
