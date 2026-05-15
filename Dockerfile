# syntax=docker/dockerfile:1.6

FROM golang:1.21-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/interview-chatbot ./cmd/server

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata && \
    adduser -D -u 10001 app
WORKDIR /app

COPY --from=build /out/interview-chatbot /app/interview-chatbot
COPY internal/database/schema.sql    /app/internal/database/schema.sql
COPY internal/database/ai_schema.sql /app/internal/database/ai_schema.sql
COPY internal/database/seed.sql      /app/internal/database/seed.sql
COPY web /app/web

RUN mkdir -p /app/voice_answers /app/tts_cache && chown -R app:app /app
USER app

EXPOSE 8080
CMD ["/app/interview-chatbot"]