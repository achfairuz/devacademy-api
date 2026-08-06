FROM golang:1.26-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
COPY vendor ./vendor

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o /app/server ./cmd/server

FROM scratch

WORKDIR /app
COPY --from=build /app/server /app/server

EXPOSE 8080

CMD ["./server"]
