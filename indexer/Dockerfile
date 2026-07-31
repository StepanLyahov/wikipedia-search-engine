FROM golang:1.22-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /indexer ./cmd/indexer

FROM alpine:3.21
RUN addgroup -S app && adduser -S app -G app
COPY --from=build /indexer /usr/local/bin/indexer
USER app
ENTRYPOINT ["/usr/local/bin/indexer"]
