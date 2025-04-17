rebuild:
	docker-compose down -v
	docker-compose build --no-cache
	docker-compose up -d
	docker-compose logs -f app

up:
	docker-compose up -d

down:
	docker-compose down -v

logs:
	docker-compose logs -f app

appRebuild:
	docker-compose build app
	docker-compose up -d --no-deps app
	docker-compose ps