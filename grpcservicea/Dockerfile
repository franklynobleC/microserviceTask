FROM golang:1.19-alpine as builder

WORKDIR  /app 

COPY go.mod ./ 
COPY go.sum ./ 
COPY ./grpcservicea/destroyer/protos/protos/proto ./

RUN go mod download

COPY  . . 

RUN go build -o docker-grpc-microservicea main.go



FROM alpine:3.13
WORKDIR /app

COPY --from=builder /app/docker-grpc-microservicea  . 

COPY ./.env .

EXPOSE 5000

CMD ["/app/docker-grpc-microservicea"]
