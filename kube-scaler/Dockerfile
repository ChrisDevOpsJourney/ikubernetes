# Build stage
FROM golang:1.20 AS builder

WORKDIR /workspace

# Copy go.mod and go.sum files separately to leverage caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o kubescaler main.go

# Final stage
FROM gcr.io/distroless/static:nonroot
COPY --from=builder /workspace/kubescaler /kubescaler

USER 1001:1001

ENTRYPOINT ["/kubescaler"]
