run-docker:
	cd ./deploy && docker-compose up -d

down-docker:
	cd ./deploy && docker-compose down

run-app:
	go run ./cmd/main.go
