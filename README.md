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
- [x] Single Binary deployment with embedded SPA
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

```sh
docker run -p 8080:8080 --rm macpla/webapp
```

### Local build

1. build ui component (React App)

```bash
./scripts/build-ui.sh
```

2. build backend with embedded SPA.

```bash
./scripts/build-local.sh
```

_optional steps_

- clean `pack2` generated files

```bash
./scripts/build-local.sh clean
```

- openApi spec generation

```bash
./scripts/openapi-gen.sh
```

## Interacting with a Service

`http://localhost:8080/` -- landing page in React
`http://localhost:8080/api/v1/demoapp/.well-known/alive?full=1` -- backend health-check

## Developer tools

[deepmap/oapi-codegen](https://github.com/deepmap/oapi-codegen) - OpenAPI 3 client and server boilerplate generator for GO

[mikefarah/yq](https://github.com/mikefarah/yq) - cli YAML preprocessor

[openapi-generator-cli](https://github.com/openapitools/openapi-generator-cli) - all purpose node based OpenAPI generator used to generate frontend client

[gobuffalo/packr/v2](github.com/gobuffalo/packr/v2) - static files to Go embedding processor
