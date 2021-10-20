
front:
	cd frontend && npm run build

docker-build:
	docker compose build

build: front docker-build

run:
	docker compose up