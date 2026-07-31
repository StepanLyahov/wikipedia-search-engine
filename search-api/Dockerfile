FROM golang:1.25-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /search-api ./cmd/search-api

FROM alpine:3.21
RUN addgroup -S app && adduser -S app -G app
COPY --from=build /search-api /usr/local/bin/search-api
USER app
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/search-api"]
