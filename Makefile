.PHONY: all build-frontend build-backend run docker-build docker-run

all: build-frontend build-backend

build-frontend:
	cd web && npm install && npm run build

build-backend:
	rm -rf cmd/server/dist
	mkdir -p cmd/server/dist
	cp -r web/dist/* cmd/server/dist/
	go build -o server-moni ./cmd/server

run: build-backend
	./server-moni

docker-build:
	docker build -t server-moni .

docker-run:
	docker run -p 8080:8080 server-moni
