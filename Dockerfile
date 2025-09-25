# Build stage
FROM golang:1.25.1 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download && go mod verify

COPY . ./

RUN go build -o /order-nest


# Runtime stage
FROM ubuntu:24.04
COPY --from=build /order-nest /order-nest
COPY ./order-nest-config.yaml order-nest-config.yaml

EXPOSE 8080

ENTRYPOINT ["/order-nest"]
CMD ["--config=order-nest-config.yaml", "serve"]