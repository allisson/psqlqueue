.PHONY: lint
lint:
	golangci-lint run -v --fix

.PHONY: test
test:
	go test -covermode=count -coverprofile=count.out -v ./...

.PHONY: build
build:
	go build -ldflags="-s -w" -o ./psqlqueue ./cmd/psqlqueue

.PHONY: build-image
build-image:
	docker build --rm -t psqlqueue .

.PHONY: run-db
run-db:
	docker run --name postgres-psqlqueue \
		--restart unless-stopped \
		-e POSTGRES_USER=psqlqueue \
		-e POSTGRES_PASSWORD=psqlqueue \
		-e POSTGRES_DB=psqlqueue \
		-p 5432:5432 \
		-d postgres:15-alpine

.PHONY: rm-db
rm-db:
	docker kill $$(docker ps -aqf name=postgres-psqlqueue)
	docker container rm $$(docker ps -aqf name=postgres-psqlqueue)

.PHONY: run-test-db
run-test-db:
	docker run --name postgres-psqlqueue-test \
		--restart unless-stopped \
		-e POSTGRES_USER=psqlqueue \
		-e POSTGRES_PASSWORD=psqlqueue \
		-e POSTGRES_DB=psqlqueue-test \
		-p 5432:5432 \
		-d postgres:15-alpine

.PHONY: rm-test-db
rm-test-db:
	docker kill $$(docker ps -aqf name=postgres-psqlqueue-test)
	docker container rm $$(docker ps -aqf name=postgres-psqlqueue-test)

.PHONY: run-migration
run-migration:
	go run cmd/psqlqueue/main.go migrate
