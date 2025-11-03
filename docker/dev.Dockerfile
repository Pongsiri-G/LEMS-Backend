FROM golang:1.25-alpine
RUN apk add make
RUN mkdir app

COPY . /app/

WORKDIR /app

RUN go install github.com/air-verse/air@latest

CMD ["air", "-c", "/app/.air.toml"]
