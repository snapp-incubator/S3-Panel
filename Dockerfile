FROM registry.snapp.tech/docker/golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go env -w GOPROXY="https://repo.snapp.tech/repository/goproxy,goproxy.io,direct"
RUN go mod download
COPY . .
WORKDIR /app/cmd/s3-panel
RUN go build -v -o s3-panel ./main.go

FROM alpine:3.20
WORKDIR /app/
COPY --from=builder /app/cmd/s3-panel .
ENTRYPOINT ["./s3-panel"]
CMD ["s3-panel", "--configPath=./config.yaml"]
