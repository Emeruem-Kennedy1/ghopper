# Variables
DOCKER_COMPOSE = docker-compose
DOCKER_COMPOSE_DEV = $(DOCKER_COMPOSE) --env-file .env.development -f docker-compose.dev.yml
DOCKER_COMPOSE_PROD = $(DOCKER_COMPOSE) --env-file .env.production -f docker-compose.prod.yml

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


#prod commands
.PHONY: prod
prod:
	$(DOCKER_COMPOSE_PROD) up --build -d

.PHONY: prod-down
prod-down:
	$(DOCKER_COMPOSE_PROD) down

.PHONY: prod-logs
prod-logs:
	$(DOCKER_COMPOSE_PROD) logs -f


# Utility commands
.PHONY: clean
clean:
	$(DOCKER_COMPOSE_DEV) down -v --remove-orphans
	$(DOCKER_COMPOSE_PROD) down -v --remove-orphans

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
	@echo "  help         Show this help message"