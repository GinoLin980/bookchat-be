FROM golang:1.27.1 AS build
# install Swagger CLI for lazy API docs auto gen 
RUN GOBIN=/usr/local/bin go install github.com/swaggo/swag/cmd/swag@latest
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# generate API docs by swag
RUN swag init -g cmd/api/main.go --output ./docs
# disable CGO to avoid no C libraries available in the distroless image, use pure Go libraries
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /app ./cmd/api

# minimum distroless image to save artifact size(i have no money for GCP :(
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /app /app
ENTRYPOINT ["/app"]
