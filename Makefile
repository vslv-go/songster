PROJECT_NAME=songster
PROJECT_PATH=src/songster
GOPATH?=$(shell go env GOPATH)

run-all: infra-up build-docs
	make -j 2 info-api-up run

run: build-docs
	cd $(PROJECT_PATH) && go run $(PROJECT_NAME)/cmd

build: build-docs
	cd $(PROJECT_PATH) && go build -o ../../bin/$(PROJECT_NAME) $(PROJECT_NAME)/cmd

test:
	cd $(PROJECT_PATH) && go test ./...

tidy:
	cd $(PROJECT_PATH) && go mod tidy

vendor: tidy
	cd $(PROJECT_PATH) && go mod vendor

build-docs:
	go install github.com/swaggo/swag/cmd/swag@latest
	cd $(PROJECT_PATH) && $(GOPATH)/bin/swag init -d ./api/rest/handlers/ -g handler.go -o ./api/rest/docs

info-api-up:
	cd fake_music_info && go run .

infra-up:
	cd infra && docker-compose -p $(PROJECT_NAME)-infra -f docker-compose-infra.yml up -d

infra-down:
	cd infra && docker-compose -p $(PROJECT_NAME)-infra -f docker-compose-infra.yml down
