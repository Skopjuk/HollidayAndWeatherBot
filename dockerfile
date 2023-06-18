FROM golang:1.20.5

WORKDIR /app

COPY ./ ./
RUN go build -o /about-me-bot

CMD ["/about-me-bot"]