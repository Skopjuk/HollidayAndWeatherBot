FROM golang:1.19

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY ./ ./
RUN go build -o /about-me-bot

CMD ["/about-me-bot"]