EXECUTABLE = statistics

build:
	go build -o $(EXECUTABLE)

run: build
	./$(EXECUTABLE)

help:
	./$(EXECUTABLE) -help

test:
	go test -v

clean:
	go clean
	rm $(EXECUTABLE)
