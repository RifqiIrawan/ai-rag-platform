.PHONY: up down build ps logs restart ollama-pull test-go test-python lint-go lint-python clean

GO_SERVICES := api-gateway auth-service document-service rag-service notification-service
PY_SERVICES := ocr-service embedding-service
OLLAMA_MODEL ?= llama3.2

up:
	docker compose up -d --build

down:
	docker compose down

build:
	docker compose build

ps:
	docker compose ps

logs:
	docker compose logs -f $(s)

restart:
	docker compose restart $(s)

ollama-pull:
	docker compose exec ollama ollama pull $(OLLAMA_MODEL)

test-go:
	@for svc in $(GO_SERVICES); do \
		echo "== $$svc =="; \
		(cd services/$$svc && go vet ./... && go test ./...) || exit 1; \
	done

test-python:
	@for svc in $(PY_SERVICES); do \
		echo "== $$svc =="; \
		(cd services/$$svc && pytest --maxfail=1 --disable-warnings) || exit 1; \
	done

lint-go:
	@for svc in $(GO_SERVICES); do \
		echo "== $$svc =="; \
		test -z "$$(cd services/$$svc && gofmt -l .)" || (echo "gofmt issues in $$svc"; exit 1); \
	done

lint-python:
	@for svc in $(PY_SERVICES); do \
		echo "== $$svc =="; \
		(cd services/$$svc && ruff check .) || exit 1; \
	done

clean:
	docker compose down -v --remove-orphans
