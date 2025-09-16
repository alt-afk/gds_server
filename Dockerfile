FROM golang:1.24.5-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server .

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]

# docker run --network=host ^
#   -e NEO4J_URI=neo4j://host.docker.internal:7687 ^
#   -e NEO4J_USER=neo4j ^
#   -e NEO4J_PASSWORD=sample-db-password ^
#   -e NEO4J_DB=neo4j
#   <img_name>

