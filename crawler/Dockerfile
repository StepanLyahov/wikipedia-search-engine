FROM golang:1.25-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /crawler ./cmd/crawler

FROM alpine:3.21
RUN addgroup -S app && adduser -S app -G app
COPY --from=build /crawler /usr/local/bin/crawler
USER app
ENTRYPOINT ["/usr/local/bin/crawler"]
