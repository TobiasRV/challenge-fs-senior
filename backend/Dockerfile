FROM golang:1.25.5-alpine3.23 AS builder
WORKDIR /backend
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
WORKDIR /backend/cmd/api
RUN go build

FROM alpine
WORKDIR /
COPY --from=builder /backend/cmd/api/api .
EXPOSE 8080
CMD ["./api"]
