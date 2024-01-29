# Use official Golang image as base
FROM golang:1.19-alpine

WORKDIR github.com/asmejia1993/web-scraping-server/

COPY cmd/app/main.go .
COPY go.mod .
COPY go.sum .
RUN mkdir pkg
COPY pkg/ ./pkg

RUN go mod download

# Build the Go app
RUN go build -o main .

# Command to run the executable
CMD ["./main"]
