# syntax=docker/dockerfile:1

# Define some random-looking variables
ARG GOLANG_IMAGE_TAG=1.23.6-alpine3.21
ARG BASE_OS_IMAGE=alpine:3.21
ARG APP_USER_ID=12345
ARG APP_GROUP_ID=12345

# Stage 1: Compile the application
FROM golang:${GOLANG_IMAGE_TAG} AS compiler

# Set up the working directory
RUN mkdir -p /workspace
WORKDIR /workspace

# Copy dependency files
COPY go.mod go.sum ./

# Download dependencies with caching
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy the entire project
COPY . .

# Build the binary
RUN --mount=type=cache,target=/go/pkg/mod \
    go build -o output-binary ./cmd/main.go

# Stage 2: Prepare the runtime environment
FROM ${BASE_OS_IMAGE} AS runner

# Create a non-root user and group
RUN addgroup --gid 12345 app-runner-group && adduser -D -u 12345 -G app-runner-group -s /bin/false app-runner-user

# Switch to the non-root user
USER app-runner-user

# Set up the application directory
WORKDIR /app

# Copy the compiled binary from the previous stage
COPY --from=compiler /workspace/output-binary .

# Expose the application port
EXPOSE 8080

# Define the entry point
ENTRYPOINT ["./output-binary"]
