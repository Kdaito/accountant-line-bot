# build stage
FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o server ./cmd/api/main.go

# run state
FROM alpine:3.13

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/credentials.json .
COPY --from=builder /app/token.json .

ARG LINE_CHANNEL_SECRET
ARG LINE_CHANNEL_TOKEN
ARG FOLDER_ID
ARG GPT_API_URL
ARG GPT_API_KEY
ARG PORT

ENV LINE_CHANNEL_SECRET=${LINE_CHANNEL_SECRET}
ENV LINE_CHANNEL_TOKEN=${LINE_CHANNEL_TOKEN}
ENV FOLDER_ID=${FOLDER_ID}
ENV GPT_API_URL=${GPT_API_URL}
ENV GPT_API_KEY=${GPT_API_KEY}
ENV PORT=${PORT}

EXPOSE $PORT

CMD ["/app/server"]
