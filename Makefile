.PHONY: dev prod build-frontend clean deps setup-frontend

# Development commands
dev:
	@echo "Starting development environment..."
	@echo "Backend: http://localhost:8080"
	@echo "Frontend: http://localhost:3000"
	@echo "Run these commands in separate terminals:"
	@echo "1. make backend"
	@echo "2. make frontend"

deps:
	go mod download
	go mod tidy

backend: deps
	go run cmd/server/main.go

setup-frontend:
	mkdir -p frontend/public frontend/src
	@if [ ! -f "frontend/package.json" ]; then \
		echo "Initializing frontend project..."; \
		cd frontend && npm init -y; \
		npm install react react-dom react-scripts axios; \
	fi

frontend: setup-frontend
	cd frontend && npm start

# Production commands
prod: deps build-frontend
	go run cmd/server/main.go

build-frontend: setup-frontend
	cd frontend && npm run build

clean:
	rm -rf frontend/build
	rm -f employees.db
	go clean -modcache 