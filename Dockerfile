FROM golang:latest

COPY ./ ./
RUN go build -o main ./cmd/
CMD ["./main"]