# Contributing

Thanks for considering a contribution to ai-rag-platform!

## Branch strategy

- `main` — stable, deployable at all times.
- `develop` — integration branch; feature branches target this.
- `feature/<service>-<short-desc>` — e.g. `feature/rag-service-retrieval-pipeline`.

Open PRs against `develop` (or `main` for urgent hotfixes).

## Commit style

Use [Conventional Commits](https://www.conventionalcommits.org/): `feat(auth-service): add login endpoint`, `fix(api-gateway): correct proxy path stripping`, `chore: ...`, `docs: ...`, `ci: ...`.

## Before opening a PR

- Go services: `gofmt -l .` clean, `go vet ./...` clean, `go test ./...` passing.
- Python services: `ruff check .` clean, `pytest` passing.
- Docker images for any touched service build successfully (`docker compose build <service>`).

These three checks map directly to the `go-ci`, `python-ci`, and `docker-build` GitHub Actions workflows, which must pass before merge.

## Code style

- Go: standard `gofmt`, idiomatic error handling (no panics for expected errors), keep handlers thin — business logic in `internal/` packages, not in route handlers.
- Python: `ruff` for linting, type hints on public functions, Pydantic models for request/response schemas.

## Code of Conduct

Participation in this project is governed by the [Code of Conduct](CODE_OF_CONDUCT.md).
