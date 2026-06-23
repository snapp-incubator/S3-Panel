# S3 Panel

[![CI](https://github.com/snapp-incubator/S3-Panel/actions/workflows/ci.yml/badge.svg)](https://github.com/snapp-incubator/S3-Panel/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/snapp-incubator/S3-Panel?sort=semver)](https://github.com/snapp-incubator/S3-Panel/releases)
[![Go](https://img.shields.io/github/go-mod/go-version/snapp-incubator/S3-Panel)](go.mod)
[![License: GPL v3](https://img.shields.io/badge/license-GPLv3-blue.svg)](LICENSE)

S3 Panel is a self-service web panel for S3-compatible object storage (Ceph RGW):
authenticate with S3 credentials and manage buckets and objects — list, upload,
download, delete and share — and view quotas.

This is a monorepo:

- **Backend** (`cmd/`, `internal/`) — a Go HTTP API in front of the Ceph RADOS Gateway.
- **Frontend** (`frontend/`) — a Vite + React + TypeScript UI, embedded into and
  served by the backend binary (`internal/web`), so there is a single image.
- **Deploy** (`deploy/helm/`) — the `s3-panel` Helm chart.

## Getting started

Run the backend with a config file:

```sh
go run ./cmd/s3-panel s3-panel --configPath=./configs/sample-config.toml
```

Configuration is TOML, loaded with [koanf](https://github.com/knadh/koanf) (defaults
< config file < `s3panel_`-prefixed environment variables). See
`configs/sample-config.toml` for the format and `frontend/README.md` for the UI.

## Deployment

A Helm chart lives in [`deploy/helm`](deploy/helm) (chart name `s3-panel`). On each
semver tag the release workflow publishes it as an OCI artifact alongside the
single image (which serves both the API and the embedded frontend):

| Artifact | Reference |
| --- | --- |
| Image (API + UI) | `ghcr.io/snapp-incubator/s3-panel:<version>` |
| Helm chart (OCI) | `oci://ghcr.io/snapp-incubator/charts/s3-panel` |

Install from the OCI registry:

```sh
helm install s3-panel oci://ghcr.io/snapp-incubator/charts/s3-panel \
  --version <version> -n snappcloud-unified-panel -f my-values.yaml
```

## CI

GitHub Actions (`.github/workflows`):

- **CI** (`ci.yml`) — **lint** (golangci-lint, Biome, `helm lint --strict`) and **test**
  (`go test`, plus chart manifest validation with kubeconform), then **build** the single
  image. The build job runs only after lint and test pass.
- **Release** (`release.yml`) — on a semver tag, builds and pushes the image and the Helm
  chart (OCI) to GHCR.

## API

The HTTP API is documented with Swagger/OpenAPI. With the server running, browse the
interactive docs at `/docs/` (e.g. <http://127.0.0.1:8080/docs/>); the generated spec lives
in [`docs/swagger.yaml`](docs/swagger.yaml).
