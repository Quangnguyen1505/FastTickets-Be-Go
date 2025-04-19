PATH_RUN=cmd/server/main.go
GOOSE_DRIVER=postgres
GOOSE_DBSTRING=user=postgres password=1234567890 dbname=shopdevgov2 host=127.0.0.1 port=25432 sslmode=disable
GOOSE_MIGRATION_DIR=sql/schema

dev:
	go run ${PATH_RUN}
run:
	docker-compose up -d && go run ${PATH_RUN}
docker_build:
	docker-compose up -d --build
	docker-compose ps
docker_up:
	docker-compose up -d
docker_down:
	docker-compose down
docker_kill:
	docker-compose kill
up_by_one:
	set "GOOSE_DRIVER=${GOOSE_DRIVER}" && set "GOOSE_DBSTRING=${GOOSE_DBSTRING}" && goose -dir ${GOOSE_MIGRATION_DIR} up-by-one
#create new migration
create_migration:
	@goose -dir=${GOOSE_MIGRATION_DIR} create ${name} sql
upGoose:
	set "GOOSE_DRIVER=${GOOSE_DRIVER}" && set "GOOSE_DBSTRING=${GOOSE_DBSTRING}" && goose -dir ${GOOSE_MIGRATION_DIR} up

downGoose:
	set "GOOSE_DRIVER=${GOOSE_DRIVER}" && set "GOOSE_DBSTRING=${GOOSE_DBSTRING}" && goose -dir ${GOOSE_MIGRATION_DIR} down

resetGoose:
	set "GOOSE_DRIVER=${GOOSE_DRIVER}" && set "GOOSE_DBSTRING=${GOOSE_DBSTRING}" && goose -dir ${GOOSE_MIGRATION_DIR} reset

sqlgen:
	sqlc generate
swag:
	swag init -g ./cmd/server/main.go ./cmd/swag/docs

.PHONY: run docker_build docker_up docker_down docker_kill
.PHONY: upGoose downGoose resetGoose