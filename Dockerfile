# Globally accessible ARGS for all intermediary containers
ARG APP_NAME=demoapp

FROM golang:1.14.4 as builder

LABEL maintainer="hippeus <plachta.maciej@gmail.com>"

ARG GO_MOD_ROOT
ARG VERSION
ARG GIT_COMMIT
ARG BUILD_TIME

# use globally defined APP_NAME in this container stage
ARG APP_NAME
ENV APP_NAME=${APP_NAME}

ENV GO_MOD_ROOT=${GO_MOD_ROOT}
ENV VERSION=${VERSION}
ENV GIT_COMMIT=${GIT_COMMIT}
ENV BUILD_TIME=${BUILD_TIME}

# remove C dependency for Go and pick target architecture
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# Create a location in the container for the source code (outside GOPATH)
RUN mkdir -p /app && mkdir -p /app/bin
WORKDIR /app

# Copy the go module files first and then download the dependencies.
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

# Build the binary
RUN go build \
  # -mod=readonly \
  -ldflags "-X ${GO_MOD_ROOT}/pkg/buildinfo.version=${VERSION} \
  -X  '${GO_MOD_ROOT}/pkg/buildinfo.commit=${GIT_COMMIT}' \
  -X  '${GO_MOD_ROOT}/pkg/buildinfo.buildTime=${BUILD_TIME}'" \
  -a -o ./bin/${APP_NAME} ./cmd/${APP_NAME}

FROM alpine:3.12

# use globally defined APP_NAME in this container stage
ARG APP_NAME
RUN mkdir -p /config

COPY --from=builder /app/bin/${APP_NAME}/ /bin/app
COPY --from=builder /app/config/${APP_NAME}.yaml /config/app.yaml

# COPY local .certs directory in dev mode
# COPY .cert/tls/ .cert/tls/

EXPOSE 8080
ENTRYPOINT ["/bin/app"]
