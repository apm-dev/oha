FROM golang:1.20-alpine AS build

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o oha

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/oha /app
COPY --from=build /app/migrations /app/migrations

EXPOSE 8000

CMD ./oha