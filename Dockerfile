FROM registry.snapp.tech/docker/golang:1.23-alpine3.20 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go env -w GOPROXY="https://repo.snapp.tech/repository/goproxy,goproxy.io,direct"
RUN go mod download
COPY . .
WORKDIR /app/cmd/snapp_object_store
RUN go build -v -o snapp_object_store ./main.go

FROM alpine:3.20
WORKDIR /app/
COPY --from=builder /app/cmd/snapp_object_store .
ENTRYPOINT ["./snapp_object_store"]
CMD ["snapp_object_store", "--configPath=./config.yaml"]
