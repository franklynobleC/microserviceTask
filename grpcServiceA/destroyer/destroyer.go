package main

import (
	"context"
	// "itoa"
	// "internal/itoa"
	// "strconv"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	se "github.com/franklynobleC/microserviceTask/grpcServiceA/destroyer/protos/protos/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nats-io/nats.go"
	// "go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
	"github.com/google/uuid"
	// "google.golang.org/protobuf/types/known/timestamppb"
)

const (
	TargetEvent = "TARGETS"
)

type server struct {
	se.UnimplementedDestroyerServiceServer
}

func NewServer() *server {
	return &server{}

}

type SingleEntryPayload struct {
	Message string `json:"message"`
}

type ResPayLoad struct {
	ID         string    `json:"id"`
	TargetName string    `json:"data,omitempty"`
	Created_on time.Time `json:"created_on"`
	//  Datas2  Datas    `json:"datas2,omitempty"`
	// ActualDataholde `r []string `json:"actual_dataholder"
}

// type DataA struct {
// 	ID         string    `json:"id,omitempty"`
// 	Message    string    `json:"message,omitempty"`
// 	Created_on time.Time `json:"created_on"`
// 	Updated_on time.Time `json:"updated_on"`
// }

func (s *server) AcquireTarget(ctx context.Context, message1 *se.DestroyerRequest) (*se.DestroyerResponse, error) {

	if len(message1.String()) == 0 {
		return nil, fmt.Errorf("no messsage entered")

	}

	data := &se.DestroyerRequest{
		Message: message1.Message,
	}
	strings.TrimSpace(strings.ToLower(data.GetMessage()))
	//  consumeWords(jst)
	// ss := &se.DestroyerResponse.UpdatedOn
	//    sss := itoa.Uitoa(uint(time.Now().Local().Unix()))
	sss, _ := time.Now().UTC().MarshalText()
	// fmt.Println(string(v))

	response := &se.DestroyerResponse{

		Id:        uuid.NewString(),
		Message:   data.GetMessage(),
		CreatedOn: string(sss),
		UpdatedOn: string(sss),
	}

	// nn = append(n, n...)

	//    convert messages to byte 
	byteMessage, err := json.Marshal(response)

	if err != nil {
		log.Println("can not marshall", err)
		//  strings(newe)
	}
	//TODO: publish to topic
	jst, err := JetStreamInit()
	if err != nil {
		log.Fatal("cant connect to nats ", err.Error())
	}

	err = CreateStream(jst)
	if err != nil {
		log.Fatal("cant create stream ", err.Error())
	}

	//TODO: published bytes Message
	jst.Publish(TargetEvent, byteMessage)

	return response, nil

}


func main() {

	grpcMux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	err := se.RegisterDestroyerServiceHandlerServer(ctx, grpcMux, NewServer())

	if err != nil {
		log.Fatal("can not register handler Server", err)
	}

	mux := http.NewServeMux()

	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", ":5000")

	if err != nil {
		log.Fatal("can not create listener", err)
	}

	log.Println("http Gateway Server is being Started", listener.Addr().String())

	err = http.Serve(listener, mux)

	if err != nil {
		log.Fatal("can not start grpc server", err)
	}
}

func CreateStream(jetStream nats.JetStreamContext) error {

	stream, err := jetStream.StreamInfo(TargetEvent)
	// stream  not found ,create it

	if stream == nil {
		log.Printf("creating stream: %s\n", TargetEvent)

		_, err = jetStream.AddStream(
			&nats.StreamConfig{
				Name: TargetEvent,
				// Storage:  nats.FileStorage,
			},
		)

		// fmt.Print(ack)
		if err != nil {
			log.Println("could not add  stream")
			return err
		}

	}
	return nil

}

func JetStreamInit() (nats.JetStreamContext, error) {

	//connect to NATS

	nc, err := ConnectToNats()
	if err != nil {
		log.Println("could not connet to nats", err)
	}

	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))

	if err != nil {
		log.Println("could not publish to Jestream")
		return nil, err
	}
	log.Println("successfully published JetStream")
	return js, nil

}

func ConnectToNats() (*nats.Conn, error) {

	nc, err := nats.Connect(os.Getenv("JESTREAM_URL"))

	if err != nil {
		log.Println("coudl not connect to Nats", err.Error())
	}

	log.Println("connected to Jetstream", nc.ConnectedAddr())

	//subscribe

	return nc, nil
}
