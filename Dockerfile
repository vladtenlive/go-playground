FROM golang:1.19-alpine

COPY . /app

WORKDIR /app
RUN go build -o app
CMD ["./app"]
