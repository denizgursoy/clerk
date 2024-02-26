# Stage 1: Build stage
FROM golang:1.21-alpine AS build

# Set the working directory
WORKDIR /app

# Copy and download dependencies

# Copy the source code
COPY . .
RUN go mod tidy

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o clerk .

# Stage 2: Final stage
FROM alpine:edge

# Set the working directory
WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /app/clerk .

# Set the timezone and install CA certificates
RUN apk --no-cache add ca-certificates tzdata

# Set the entrypoint command
ENTRYPOINT ["/app/clerk"]