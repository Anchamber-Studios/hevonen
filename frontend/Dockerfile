# Use the official Go image as the base image
ARG PORT=4443
ARG HOST=localhost
ARG TLS_ENABLED=false
ARG TLS_KEY=tls/key.pem
ARG TLS_CERT=tls/cert.pem
ARG MEMBERS_SERVICE_URL=https://localhost:8443

FROM docker.io/library/golang:1.21-alpine as build

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download and install the Go dependencies
RUN go mod download

# Copy the source code into the container
COPY . .
RUN rm -rf .git
RUN rm -rf tmp

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o frontend .

# Set the entry point for the container

FROM docker.io/library/alpine:latest
ARG PORT
ARG HOST
ARG TLS_ENABLED
ARG TLS_KEY
ARG TLS_CERT
ARG MEMBERS_SERVICE_URL

USER nobody
ENV PORT=${PORT}
ENV HOST=${HOST}
ENV TLS_ENABLED=${TLS_ENABLED}
COPY --from=build /app/frontend /frontend
ENTRYPOINT ["./frontend"]
