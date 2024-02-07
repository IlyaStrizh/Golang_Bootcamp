// curl -s --key ex02/minica/client/key.pem --cert ex02/minica/client/cert.pem --cacert ex02/minica/minica.pem -XPOST -H "Content-Type: application/json" -d '{"candyType": "NT", "candyCount": 1, "money": 34}' "https://localhost:3333/buy_candy"
// ./client -k AA -c 1 -m 50

package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Флаги командной строки
	key := flag.String("k", "", "Candy type abbreviation")
	count := flag.Int("c", 0, "Count of candy to buy")
	money := flag.Int("m", 0, "Amount of money given to machine")
	flag.Parse()

	// Загрузка клиентского certificate и key
	cert, err := tls.LoadX509KeyPair("ex02/minica/client/cert.pem", "ex02/minica/client/key.pem")
	if err != nil {
		log.Fatal("Failed to load client certificate and key:", err)
	}

	// Загрузка CA certificate
	caCert, err := os.ReadFile("ex02/minica/minica.pem")
	if err != nil {
		log.Fatal("Failed to read CA certificate:", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Конфигурирование TLS
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}

	// Создание HTTPS client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	// Запрос к серверу
	request := fmt.Sprintf(`{"money": %d, "candyType": "%s", "candyCount": %d}`, *money, *key, *count)
	resp, err := client.Post("https://localhost:3333/buy_candy", "application/json", strings.NewReader(request))
	if err != nil {
		log.Fatal("Failed to send request:", err)
	}
	defer resp.Body.Close()

	// Чтение ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Failed to read response:", err)
	}

	var data map[string]interface{}

	err = json.Unmarshal(body, &data) // распаковываем данные из формата JSON в структуру или мапу
	if err != nil {
		log.Fatal("Failed to Unmarshal:", err)
		return
	}

	// Печать ответа
	if len(data) > 1 {
		fmt.Printf("%s Your change is %v\n", data["thanks"], data["change"])
	} else {
		fmt.Println(data["error"])
	}
}
