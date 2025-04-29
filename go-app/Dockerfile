# ./go-worker/Dockerfile
FROM golang:latest as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o go-worker-app ./main.go # Ensure ./main.go is correct path to your main package file

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/go-worker-app .

# Define environment variables (optional here, but good practice)
# These will be overridden by docker-compose environment section
# ENV REDIS_HOST="localhost"
# ENV REDIS_PORT="6379"
# ENV REDIS_QUEUE_NAME="garden_actions"
# ENV POSTGRES_HOST="localhost"
# ENV POSTGRES_PORT="5432"
# ENV POSTGRES_DB="admin"
# ENV POSTGRES_USER="admin"
# ENV POSTGRES_PASSWORD="admin"


# Command to run the application when the container starts
CMD ["./go-worker-app"]
