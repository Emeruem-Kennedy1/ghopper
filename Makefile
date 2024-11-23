# Variables
DOCKER_COMPOSE = docker-compose
DOCKER_COMPOSE_DEV = $(DOCKER_COMPOSE) --env-file .env.development -f docker-compose.dev.yml
DOCKER_USERNAME = kennedyemeruem
VERSION=latest

#dev commands
.PHONY: dev
dev:
	$(DOCKER_COMPOSE_DEV) up --build -d

.PHONY: dev-down
dev-down:
	$(DOCKER_COMPOSE_DEV) down

.PHONY: dev-logs
dev-logs:
	$(DOCKER_COMPOSE_DEV) logs -f


# prod commands
.PHONY: prod
apply-secrets:
	kubectl apply -f k3s/config/secrets.yaml
apply-app:
	kubectl apply -f k3s/apps/
apply-all: apply-secrets apply-app

# Build and push images to dockerhub
.PHONY: build push deploy

build:
	docker build -t $(DOCKER_USERNAME)/ghopper-frontend:$(VERSION) -f frontend/Dockerfile.prod .
	docker build -t $(DOCKER_USERNAME)/ghopper-backend:$(VERSION) -f backend/Dockerfile.prod .

push:
	docker push $(DOCKER_USERNAME)/ghopper-frontend:$(VERSION)
	docker push $(DOCKER_USERNAME)/ghopper-backend:$(VERSION)

# Build and push the images
deploy: build push


# Utility commands
.PHONY: clean
clean:
	$(DOCKER_COMPOSE_DEV) down -v --remove-orphans

.PHONY: prune
prune:
	docker system prune -af

.PHONY: backend-shell
backend-shell:
	$(DOCKER_COMPOSE_DEV) exec backend sh

.PHONY: frontend-shell
frontend-shell:
	$(DOCKER_COMPOSE_DEV) exec frontend sh




# help
.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available commands:"
	@echo "  dev          Start the development environment"
	@echo "  dev-down     Stop the development environment"
	@echo "  dev-logs     Show logs of the development environment"
	@echo "  prod         Start the production environment"
	@echo "  prod-down    Stop the production environment"
	@echo "  prod-logs    Show logs of the production environment"
	@echo "  clean        Stop and remove all containers, networks, images, and volumes"
	@echo "  prune        Remove all unused containers, networks, images, and volumes"
	@echo "  backend-shell    Open a shell in the backend container"
	@echo "  frontend-shell    Open a shell in the frontend container"
	@echo "  build        Build the images"
	@echo "  push         Push the images to Docker Hub"
	@echo "  deploy       Build and push the images"
	@echo "  apply-secrets    Apply the secrets"
	@echo "  apply-app    Apply the app"
	@echo "  apply-all    Apply all"
	@echo "  help         Show this help message"