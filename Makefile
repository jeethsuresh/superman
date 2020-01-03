

setup:
	#TODO: install golang
	sudo curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.22.2
	
build: 
	go build -o bin/superman cmd/superman/main.go
	
image: 
	docker build -t superman:v1 .

clean:
	rm -rf bin/ 

lint: 
	golangci-lint run

fmt: 
	go fmt ./...

unit:
	./test/unit.sh

test: unit

db: 
	rm -rf ./data/superman.db 
	sqlite3 ./data/superman.db 'CREATE TABLE locations (id varchar(255) NOT NULL PRIMARY KEY, timestamp bigint, username varchar(255), lat float, lon float, radius int, ip string);'
