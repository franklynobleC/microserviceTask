FROM golang:1.19-alpine as builder

WORKDIR  /app 

COPY go.mod ./ 
COPY go.sum ./ 
COPY ./grpcServiceA/destroyer/protos/protos/proto ./

RUN go mod download

COPY  . . 

RUN go build -o docker-grpc-microservice  grpcServiceA/main.go





FROM alpine:3.13
WORKDIR /app

COPY --from=builder /app/docker-grpc-microservice  . 
# COPY --from=builder /app/migrate.linux-amd64 ./migrate
COPY ./grpcServiceA/.env .


EXPOSE 5000

CMD ["/app/docker-grpc-microservice"]
