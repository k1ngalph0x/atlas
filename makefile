.PHONY: all stop identity ingestion issue alert intelligence build test tidy

SERVICES = identity ingestion issue alert intelligence

all:
	@echo "Starting all Atlas services..."
	@start "identity-service"    cmd /k "cd services/identity-service    && go run main.go"
	@start "ingestion-service"   cmd /k "cd services/ingestion-service   && go run main.go"
	@start "issue-service"       cmd /k "cd services/issue-service       && go run main.go"
	@start "alert-service"       cmd /k "cd services/alert-service       && go run main.go"
	@start "intelligence-service" cmd /k "cd services/intelligence-service && go run main.go"
	@echo "All services started in separate windows"

identity:
	cd services/identity-service && go run main.go

ingestion:
	cd services/ingestion-service && go run main.go

issue:
	cd services/issue-service && go run main.go

alert:
	cd services/alert-service && go run main.go

intelligence:
	cd services/intelligence-service && go run main.go


build:
	@echo "Building all services..."
	cd services/identity-service     && go build -o bin/identity.exe     main.go
	cd services/ingestion-service    && go build -o bin/ingestion.exe    main.go
	cd services/issue-service        && go build -o bin/issue.exe        main.go
	cd services/alert-service        && go build -o bin/alert.exe        main.go
	cd services/intelligence-service && go build -o bin/intelligence.exe main.go
	@echo "Build complete"


stop:
	@echo "Stopping all services..."
	@for %%p in (8080 8081 8082 8083 8084) do \
		@for /f "tokens=5" %%a in ('netstat -ano ^| findstr :%%p ^| findstr LISTENING') do \
			taskkill /PID %%a /F 2>nul || true
	@echo "Done"


test:
	cd services/identity-service  && go test ./api/... -v
	cd services/ingestion-service && go test ./api/... -v


tidy:
	cd services/identity-service     && go mod tidy
	cd services/ingestion-service    && go mod tidy
	cd services/issue-service        && go mod tidy
	cd services/alert-service        && go mod tidy
	cd services/intelligence-service && go mod tidy