postgres-up:
	migrate -path ./migrations -database "postgres://katawa:root@localhost:5432/data_generator_mock?sslmode=disable" up

postgres-down:
	migrate -path ./migrations -database "postgres://katawa:root@localhost:5432/data_generator_mock?sslmode=disable" down

docker-up:
	docker-compose up -d
	sleep 3