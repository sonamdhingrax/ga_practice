# ga_practice

A tiny Go HTTP service used as a hands-on lab for mastering **GitHub Actions** —
from solid CI to multi-arch images on GHCR, supply-chain security, and keyless
deployment to **AWS ECS Fargate** via OIDC.

## The app

| Endpoint   | Purpose                                            |
|------------|----------------------------------------------------|
| `GET /`        | Returns JSON build metadata (version, commit, date) |
| `GET /healthz` | Liveness probe, returns `200 ok`                    |

Build metadata is injected at compile time with `-ldflags -X`, so the running
service tells you exactly which commit/tag produced it.

## Run locally

```bash
go test ./...                 # run the test suite
go run .                      # starts on :8080 (override with PORT)
curl localhost:8080/ ; echo   # {"version":"dev",...}
curl localhost:8080/healthz   # ok
```

## Build the container

```bash
docker build -t ga .
docker run --rm -p 8080:8080 ga
```

## Learning roadmap

This repo is built over a 3-day plan:

- **Day 1 — Solid CI:** triggers, matrix, caching, concurrency, multi-arch Docker → GHCR.
- **Day 2 — Reusable + Secure:** composite actions, reusable workflows, CodeQL,
  Dependabot, SBOM, cosign signing, SLSA provenance attestations.
- **Day 3 — Deploy:** AWS OIDC (no static keys), environments + approvals,
  ECS Fargate deployment, and a tag-driven capstone release pipeline.

See `docs/` for cheat-sheets and the troubleshooting log.
