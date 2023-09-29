##FROM golang:1.21
#FROM golang:alpine as builder
#
#LABEL maintainer="Diyorbek Abdulaxatov <dabdulaxatov036@gmail.com>"
#
#RUN apk update && apk add --no-cache git
#
## Install git
#WORKDIR /app
#
## Copy go.mod and go.sum files
#COPY go.mod go.sum ./
#
## Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed
#RUN go mod tidy
#
## Copy the source from the current directory to the working Directory inside the container
#COPY . .
#
## Build the go app
#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go
#
## Start the new stage from scratch
#FROM alpine:latest
#RUN apk --no-cache add ca-certificates
#
#WORKDIR /root/
#
## Copy the Pre-built binary file from the previous stage. Observe we also copied the .env file
#COPY --from=builder /app/main .
#COPY --from=builder /app/.env .
#
#EXPOSE 8080
#
## Command to run the executable
#CMD ["./main"]
# Step 1: Modules caching
FROM golang:alpine as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:alpine as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bin/app ./cmd/app

# Step 3: Final
FROM scratch
COPY --from=builder /app/config /config
COPY --from=builder /app/migrations /migrations
COPY --from=builder /bin/app /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["/app"]
