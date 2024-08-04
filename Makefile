run-docker:
	cd ./deploy && docker-compose up -d

run-app:
	go run ./cmd/main.go
