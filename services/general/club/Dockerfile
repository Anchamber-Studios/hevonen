# Use the official Go image as the base image
ARG PORT=8443
ARG HOST=localhost
ARG TLS_ENABLED=false
ARG TLS_KEY=tls/key.pem
ARG TLS_CERT=tls/cert.pem
ARG ORY_URL=http://localhost:4000/.ory

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
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# Set the entry point for the container

FROM docker.io/library/alpine:latest
ARG PORT
ARG HOST
ARG TLS_ENABLED
ARG TLS_KEY
ARG TLS_CERT
ARG ORY_URL

USER nobody
ENV PORT=${PORT}
ENV HOST=${HOST}
ENV TLS_ENABLED=${TLS_ENABLED}
COPY --from=build /app/app /app
ENTRYPOINT ["./app"]
