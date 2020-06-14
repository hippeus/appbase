# Default WebApp template

Test project to build simple RESTful service based on OpenApi definitions defined in `/api` directory.

## TODOs

- [x] Configuration with Viper
- [ ] Add support for flags (cobra?)
- [x] Version tags / Build info
- [x] Containerization
- [x] OpenApi code generation
- [x] REST server using echo.Echo
- [x] Middleware for REST endpoints
- [x] Picking structured logging library
- [x] Dummy Frontend SPA
- [ ] Single Binary deployment with embedded SPA
- [ ] Persistence storage
- [ ] SearchStore capabilities

## Directory layout

| Directory      | Description                                                         |
| :------------- | :------------------------------------------------------------------ |
| `./.iac`       | "infrastructure as code"                                            |
| `./api`        | OpenAPI, gRPC contracts                                             |
| `./bin/<arch>` | Build output directory, grouped by architecture (eg: darwin, linux) |
| `./cmd`        | application main                                                    |
| `./config`     | application specific configuration                                  |
| `./internal`   | project private libraries                                           |
| `./pkg`        | project libraries                                                   |
| `./scripts`    | Shell scripts used in dev environment (env setup, automation etc)   |
| `./ui`         | frontend as Single Page Application                                 |

## Prerequisites

Minimum

- Go 1.14+
- Go Modules

Recommended

- Git (SCM)
- Docker (containerization)
- Unix based system (for scripts execution)
- Javascript with React (for UI component development)
- yarn (frontend package manager)
- developer tools binaries installed on host

## Installation

### Image build

```sh
gomodroot=$(go list -m)
buildTime=$(date -u '+%Y-%m-%d %H:%M:%S %Z') && \
docker build \
  --build-arg VERSION=$(git describe --tags --always --dirty) \
  --build-arg GIT_COMMIT=$(git rev-list -1 HEAD) \
  --build-arg BUILD_TIME=$buildTime \
  --build-arg GO_MOD_ROOT=$gomodroot \
  -t macpla/webapp .
```

### Docker run

### Local build

## Interacting with a Service

## Developer tools

[deepmap/oapi-codegen](https://github.com/deepmap/oapi-codegen) - OpenAPI 3 client and server boilerplate generator for GO

[mikefarah/yq](https://github.com/mikefarah/yq) - cli YAML preprocessor

[openapi-generator-cli](https://github.com/openapitools/openapi-generator-cli) - all purpose node based OpenAPI generator used to generate frontend client

[markbates/pkger](https://github.com/markbates/pkger) - static files to Go embedding processor
