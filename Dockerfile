# ===========================
# Stage 1: Build the Frontend (Vite)
# ===========================
FROM node:20-alpine AS frontend

# Set working directory
WORKDIR /MedscribeUI

# Copy package.json and yarn.lock
COPY MedscribeUI/package.json ./
COPY MedscribeUI/yarn.lock ./

# Install dependencies
RUN yarn install

# Copy the rest of the frontend code
COPY MedscribeUI/ .

# Determine build mode based on service name build argument
ARG _SERVICE_NAME
RUN if [ "$_SERVICE_NAME" = "medscribe-prod" ]; then \
    echo "Building for Production (Service Name: $_SERVICE_NAME)"; \
    yarn build --mode production; \
elif [ "$_SERVICE_NAME" = "medscribe-dev" ]; then \
    echo "Building for Development (Service Name: $_SERVICE_NAME)"; \
    yarn build --mode development; \
else \
    echo "Building for local or other environment (Service Name: $_SERVICE_NAME)"; \
    yarn build --mode test; \
fi

# ===========================
# Stage 2: Build the Go Backend
# ===========================

FROM golang:1.23 AS builder

# Set working directory inside the container
WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the entire backend source code
COPY . .

# Build a statically compiled Go binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app cmd/api/api.go


# ===========================
# Stage 3: Production Image (Final)
# ===========================

# Use a minimal base image
FROM alpine:latest

# Set working directory inside the container to /app
WORKDIR /app

# Copy the .env file into /app
COPY .env .

# Copy only the compiled Go binary from the previous stage
COPY --from=builder /app/app .
COPY --from=frontend /MedscribeUI/dist ./MedscribeUI/dist

# Ensure the binary is executable
RUN chmod +x /app/app

EXPOSE 8080

ENV ENVIRONMENT="production"
ENV GOOGLE_APPLICATION_CREDENTIALS=${GOOGLE_APPLICATION_CREDENTIALS}

CMD ["/app/app"]