SERVER = apiServer
DATABASE = elasticDatabase

ex00:
	cd Database && go build -o $(DATABASE)
	mv Database/$(DATABASE) ./

ex01_04:
	cd Server && go build -o $(SERVER)
	mv Server/$(SERVER) ./

build: ex00 ex01_04

clean:
	go clean
	rm -f $(SERVER) $(DATABASE)

.PHONY: ex00 ex01_04