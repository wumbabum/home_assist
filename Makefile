# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: run quality control checks
.PHONY: audit
audit: test
	go mod tidy -diff
	go mod verify
	test -z "$(shell gofmt -l .)" 
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

## test: run all tests
.PHONY: test
test: test/db/setup
	gotestsum --format testname -- -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

## upgradeable: list direct dependencies that have upgrades available
.PHONY: upgradeable
upgradeable:
	@go run github.com/oligot/go-mod-upgrade@latest

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## tidy: tidy modfiles and format .go files
.PHONY: tidy
tidy:
	go mod tidy -v
	go fmt ./...

## build: build the cmd/web application
.PHONY: build
build:
	go build -o=/tmp/bin/web ./cmd/web
	
## run: run the cmd/web application
.PHONY: run
run: build
	/tmp/bin/web

## run/live: run the application with reloading on file changes
.PHONY: run/live
run/live:
	go run github.com/cosmtrek/air@v1.43.0 \
		--build.cmd "make build" --build.bin "/tmp/bin/web" --build.delay "100" \
		--build.exclude_dir "" \
		--build.include_ext "go, tpl, tmpl, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico" \
		--misc.clean_on_exit "true"


# ==================================================================================== #
# SQL MIGRATIONS
# ==================================================================================== #

## migrations/new name=$1: create a new database migration
.PHONY: migrations/new
migrations/new:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest create -seq -ext=.sql -dir=./assets/migrations ${name}

## migrations/up: apply all up database migrations
.PHONY: migrations/up
migrations/up:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="${DB_DSN}" up

## migrations/down: apply all down database migrations
.PHONY: migrations/down
migrations/down:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="${DB_DSN}" down

## migrations/goto version=$1: migrate to a specific version number
.PHONY: migrations/goto
migrations/goto:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="${DB_DSN}" goto ${version}

## migrations/force version=$1: force database migration
.PHONY: migrations/force
migrations/force:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="${DB_DSN}" force ${version}

## migrations/version: print the current in-use migration version
.PHONY: migrations/version
migrations/version:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="${DB_DSN}" version

# ==================================================================================== #
# DATABASE TESTING
# ==================================================================================== #

## test/db/setup: create and migrate test database
.PHONY: test/db/setup
test/db/setup:
	@echo "Creating test database..."
	@psql postgres://home_assist_user:dev_password@127.0.0.1:5432/postgres -c "DROP DATABASE IF EXISTS home_assist_test;" 2>/dev/null || true
	@psql postgres://home_assist_user:dev_password@127.0.0.1:5432/postgres -c "CREATE DATABASE home_assist_test;"
	@echo "Running migrations..."
	@TEST_DB_DSN="postgres://home_assist_user:dev_password@127.0.0.1:5432/home_assist_test?sslmode=disable" go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path=./assets/migrations -database="$$TEST_DB_DSN" up

## test/db/teardown: drop test database
.PHONY: test/db/teardown
test/db/teardown:
	@echo "Dropping test database..."
	@psql postgres://home_assist_user:dev_password@127.0.0.1:5432/postgres -c "DROP DATABASE IF EXISTS home_assist_test;" 2>/dev/null || true

## test/db: run database tests with test database
.PHONY: test/db
test/db: test/db/setup
	@TEST_DB_DSN="postgres://home_assist_user:dev_password@127.0.0.1:5432/home_assist_test?sslmode=disable" go test -v ./internal/database
	@$(MAKE) test/db/teardown
