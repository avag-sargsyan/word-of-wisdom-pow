install:
	go mod download

test:
	go clean --testcache
	go test ./...

start:
	docker-compose up --abort-on-container-exit --force-recreate --build server client
