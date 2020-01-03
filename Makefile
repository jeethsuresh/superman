

setup:
	#TODO: install golang
	sudo curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.22.2
	
build: 
	go build -o bin/superman cmd/superman/main.go
	
image: 
	docker build .

clean:
	rm -rf bin/ 

lint: 
	golangci-lint run

fmt: 
	go fmt ./...

unit:
	./test/unit.sh

functional:
	./test/functional.sh

test: unit functional

db: 
	rm -rf ./data/superman.db 
	sqlite3 ./data/superman.db 'CREATE TABLE locations (id varchar(255) NOT NULL PRIMARY KEY, timestamp bigint UNIQUE, username varchar(255), lat float, lon float, speed int, radius int, ip string);'
	#sqlite3 ./data/superman.db "INSERT INTO locations (id, username, timestamp) VALUES ('33765386-C644-4AA4-A34B-EF04D2BD0E59', 'jeeth', 1578020060);"