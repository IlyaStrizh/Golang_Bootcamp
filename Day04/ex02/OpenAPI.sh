#!/bin/bash

docker run --rm \
  -v ${PWD}:/local openapitools/openapi-generator-cli generate \
  -i /local/vending.yaml \
  -g go-gin-server \
  -o /local/server/go