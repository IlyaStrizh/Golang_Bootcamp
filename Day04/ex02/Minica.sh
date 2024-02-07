#!/bin/bash

# Пути и имена файлов для сохранения ключей и сертификатов
CLIENT_KEY_FILE="client-key.pem"
CLIENT_CERT_FILE="client.pem"
LOCALHOST_KEY_FILE="localhost-key.pem"
LOCALHOST_CERT_FILE="localhost.pem"

# Запуск контейнера Docker с minica
docker run --rm \
  -v "${PWD}:/root/.minica" \
  -w "/root/.minica" \
  golang:latest \
  sh -c "git clone https://github.com/jsha/minica.git && cd minica && go build && ./minica --domains 'client'"

docker run --rm \
  -v "${PWD}:/root/.minica" \
  -w "/root/.minica" \
  golang:latest \
  sh -c "cd minica && ./minica --domains 'localhost'"

echo "Генерация ключей и сертификатов успешно выполнена."
echo "Созданные ключи и сертификаты находятся в текущей папке: ${PWD}"
echo "Ключ и сертификат для клиента:"
echo "  Ключ: ${CLIENT_KEY_FILE}"
echo "  Сертификат: ${CLIENT_CERT_FILE}"
echo "Ключ и сертификат для server:"
echo "  Ключ: ${LOCALHOST_KEY_FILE}"
echo "  Сертификат: ${LOCALHOST_CERT_FILE}"
