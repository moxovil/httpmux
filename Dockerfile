FROM golang:latest

COPY ./ ./
RUN go build -o main ./web
CMD ["./main"]