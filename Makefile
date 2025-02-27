build:
	docker-compose build booking-system
run:
	docker-compose up
test:
	go test -v ./...
migrate_up:
	migrate -path ./schema -database 'postgres://postgres:postgres@0.0.0.0:5436/postgres?sslmode=disable' up
migrate_down:
	migrate -path ./schema -database 'postgres://postgres:postgres@0.0.0.0:5436/postgres?sslmode=disable' down
migrate_down_force:
	migrate -path ./schema -database 'postgres://postgres:postgres@0.0.0.0:5436/postgres?sslmode=disable' force 1