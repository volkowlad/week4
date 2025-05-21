memory:
	STORAGE=memory docker compose --env-file .env -f docker-compose.yml up

postgres:
	STORAGE=postgres docker compose --env-file .env -f docker-compose.yml up

down:
	docker compose -f docker-compose.yml down
