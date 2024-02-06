#!/bin/bash

if ! [[ $(docker image ls -q elasticsearch:7.9.2) ]]; then
    docker pull elasticsearch:7.9.2
fi

# docker network create elastic
# docker run --name es01 --net elastic -p 9200:9200 -it -m 1GB docker.elastic.co/elasticsearch/elasticsearch:8.12.0

docker run -d --name elasticsearch \
    -p 9200:9200 -p 9300:9300 \
    -e "discovery.type=single-node" \
    -v elasticsearch:/usr/share/elasticsearch/data \
    elasticsearch:7.9.2