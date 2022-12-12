FROM golang:1.19-alpine as builder

WORKDIR  /app 

COPY  go.mod  ./ 
COPY  go.sum  ./ 
# COPY /grpcServiceA/destroyer/protos/protos/proto ./

RUN go mod download

# RUN go get github.com/cespare/reflex
# COPY reflex.conf /

COPY  . . 

RUN go build -o docker-grpc-microserviceb   main.go



FROM alpine:3.13
WORKDIR /app

COPY --from=builder /app/docker-grpc-microserviceb  . 
# COPY --from=builder /app/migrate.linux-amd64 ./migrate
COPY .env .

EXPOSE 7000

CMD ["/app/docker-grpc-microserviceb"]
