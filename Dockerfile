# ===========================
# Stage 1: Build the Frontend (Vite)
# ===========================
FROM node:18-alpine AS frontend

# Set working directory
WORKDIR /MedscribeUI
RUN echo "=== After WORKDIR /app  ==" && ls -la /

# Copy package.json and yarn.lock
COPY MedscribeUI/package.json ./
COPY MedscribeUI/yarn.lock ./
RUN echo "=== After copying package.json and yarn.lock ===" && ls -la /MedscribeUI

# Install dependencies
RUN yarn install
RUN echo "=== After yarn install ===" && ls -la /MedscribeUI

# Copy the rest of the frontend code
COPY MedscribeUI/ .
RUN echo "=== After copying MedscribeUI ===" && ls -la /MedscribeUI

# Build the application
RUN yarn build
RUN echo "=== After yarn build  ===" && ls -la /MedscribeUI && ls -la dist


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

# (Optional) Print directory contents to confirm successful copy
RUN echo "=== Listing all files in /app ===" && ls -al /app

# (Optional) Print the contents of a specific file
RUN echo "=== Contents of cmd/api/api.go ===" && cat cmd/api/api.go

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

CMD ["/app/app"]
