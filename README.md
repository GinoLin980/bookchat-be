# BookChat Backend API

BookChat backend API written in **Go**, deployed to Google Cloud Run with minimum, non-root distroless Docker runtime.
Super small Docker image for only 10.8MB

The application entry point is `cmd/api/main.go`.

## Getting started

Choose one of the following ways to run the application:

| Method | Requirements | Recommended for |
| --- | --- | --- |
| Docker | Docker installed and running | Running the application without installing Go |
| Native Go | Go version `go 1.27.1` required by `go.mod` | Development and debugging |

Clone the repository and enter its directory:

```bash
git clone https://github.com/GinoLin980/bookchat-be.git
cd bookchat-be
```

Create a `.env` file in the project root:

```dotenv
PORT=8080
# PORT will fallback to 8080 if not set
# Add the other environment variables required by the application.
# Example include database configuration and authentication secrets.
```

If the repository includes `.env.example`, copy it instead:

```bash
cp .env.example .env
```

Replace placeholder values with your own configuration.
```env
JWT_SECRET=SECRET
PG_HOST=host
PG_PORT=5432
PG_USER=postgres
PG_PSWD=pswd
PG_DB=postgres
```

Use plain `KEY=value` entries so the file also works with Docker's
`--env-file` option. Avoid shell expressions and variable interpolation.

Do not commit `.env` or real credentials.

## Run with Docker

Go and the Swagger CLI are installed inside the build image.
You only need Docker on your machine.

### Build the image

Run this command from the project root:

```bash
docker build -t bookchat-be:local .
```

The Dockerfile:

- Downloads Go dependencies.
- Installs the Swagger CLI.
- Generates API documentation from `cmd/api/main.go`.
- Builds the application with CGO disabled to prevent calling C libraries in a distroless, which doesn't include the C libraries.
- Copies the compiled executable into a non-root distroless image.

The final image contains the runtime and application executable,
not the Go compiler or project source tree.

### Start the application

```bash
docker run --rm \
  --name bookchat-be \
  --env-file .env \
  -e PORT=8080 \
  -p 127.0.0.1:8080:8080 \
  bookchat-be:local
```

This runs in the foreground and binds the published port to your local
machine only.

Press `Ctrl+C` to stop it.

To run in the background instead:

```bash
docker run -d \
  --name bookchat-be \
  --env-file .env \
  -e PORT=8080 \
  -p 127.0.0.1:8080:8080 \
  bookchat-be:local
```

View logs:

```bash
docker logs -f bookchat-be
```

Stop and remove the background container:

```bash
docker stop bookchat-be
docker rm bookchat-be
```

If host port `8080` is already in use, change the port mapping to:

```bash
-p 127.0.0.1:8081:8080
```

The application still listens on port `8080` inside the container.
You access it through port `8081` on your machine.

The application must listen on `0.0.0.0:8080` or `:8080` inside the
container, not `127.0.0.1:8080`.

## Run with Go

Install the Go version required by `go.mod`.

Download dependencies:

```bash
go mod download
```

Generate Swagger documentation:

```bash
go run github.com/swaggo/swag/cmd/swag@latest \
  init -g cmd/api/main.go --output ./docs
```

Set the required environment variables before starting the application.

For example, on Linux or macOS:

```bash
export PORT=8080
# Export the application's other required variables here.
```

Go does not automatically load `.env` files. If the application includes
a dotenv loader, it may load `.env` itself; otherwise, configure your
shell environment.

Start the API:

```bash
go run ./cmd/api/main.go
```

To build an executable instead:

```bash
mkdir -p bin
CGO_ENABLED=0 go build \
  -trimpath \
  -ldflags="-s -w" \
  -o ./bin/api \
  ./cmd/api
```

Run it:

```bash
./bin/api
```

The executable is built for your current operating system and architecture.

## API documentation

Swagger documentation is generated from annotations in the Go source code.

Regenerate it after changing API annotations:

```bash
go run github.com/swaggo/swag/cmd/swag@latest \
  init -g cmd/api/main.go --output ./docs
```

Generated files are written to `docs/`.

Docker builds generate these files automatically.

Generating documentation does not automatically expose Swagger UI.
The application must register a Swagger route to serve it. Check the
application's route configuration for the correct URL.

## Troubleshooting and security

### Missing configuration

Docker's `--env-file .env` passes environment variables to the process.
It does not copy or mount the `.env` file into the container.

The application should accept environment variables without requiring
a physical `.env` file at startup.

### Database connection errors

Inside a container, `localhost` refers to that container, not your
computer or another container.

Configure the database hostname for your environment. When using
Docker Compose, this is typically the database service name.

Any databases or other external services required by the application
must be configured separately. This Dockerfile builds only the API.

### No shell inside the container

The distroless runtime does not include `sh`, `bash`, or a package manager.

Commands such as this will not work:

```bash
docker exec -it go-api sh
```

Use application logs for routine troubleshooting.

### Missing runtime files

The Dockerfile copies only the compiled executable into the final image.

If the application reads templates, static files, migrations, or other
files at runtime, copy them into the runtime stage or embed them in the
Go executable.

### Keep secrets out of builds

Add the following entries to `.dockerignore`:

```gitignore
.git
.env
.env.*
!.env.example
bin/
tmp/
```

Also exclude `.env` and generated binaries from Git.

### Reproducible builds

Use a Go builder version compatible with `go.mod` and verify that the
chosen Docker image tag exists.

The current Dockerfile installs `swag@latest`. For more reproducible
builds, replace `@latest` with a tested version and use the same version
in the local documentation-generation command.

The `gcr.io/distroless/...` image is a container base image, not a
requirement to deploy the application on Google Cloud.
