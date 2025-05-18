FROM golang:1.23-alpine

# Enable CGO and set environment variables
ENV CGO_ENABLED=1
ENV GO111MODULE=on

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk update && apk add --no-cache gcc musl-dev sqlite sqlite-dev

# Copy all project files
COPY . .

# Initialize go module and download dependencies
RUN go mod tidy

# Build the Go binary
RUN go build -o server .

# Seed the database during build
RUN sqlite3 database/shifts.db < database/seed.sql

# Expose backend port
EXPOSE 8000

# Run the compiled app
CMD ["./server"]
