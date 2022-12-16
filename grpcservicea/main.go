package main

import (
	"context"
	"math/rand"
	"time"
	// "strconv"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	se "github.com/franklynobleC/microserviceTask/grpcservicea/destroyer/protos/protos/proto"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"github.com/oklog/ulid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	TargetEvent = "TARGET"
	EventName   = "targets.acquired"
	EventAction = ""
)

type server struct {
	se.UnimplementedDestroyerServiceServer
	se.DestroyerServiceServer
}

func NewServer() *server {
	return &server{}

}

type SingleEntryPayload struct {
	Message string `json:"message"`
}

func ConnectMongo() (*mongo.Collection, error) {

	// Get DB data from .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("could not Load .env file")
	}

	opts := options.Client().ApplyURI(os.Getenv("DB_URL"))

	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		log.Fatal("could not connect to mongo Db")
	}
	fmt.Print("database connected successfully FROM mongo func")

	deathstarCollection := client.Database(os.Getenv("DB_NAME")).Collection("deathstar")

	fmt.Print("database created", deathstarCollection.Database())

	return deathstarCollection, nil
}

func (srv *server) AcquireTarget(ctx context.Context, message1 *se.AcquireTargetRequest) (*se.AcquireTargetResponse, error) {

	if len(message1.String()) == 0 {
		return nil, fmt.Errorf("no messsage entered")

	}

	data := &se.AcquireTargetRequest{
		Message: message1.Message,
	}
	strings.TrimSpace(strings.ToLower(data.GetMessage()))

	//conver date to string

	v, _ := time.Now().UTC().MarshalText()
	dateString := string(v)

	//convert is to uuids autogenerate
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	uids := id.String()

	//initialise Values to proto, to return to respo
	response := &se.AcquireTargetResponse{

		Id:        uids,
		Message:   data.GetMessage(),
		Createdat: dateString,
		Updatedat: dateString,
	}

	//    convert messages to byte
	byteMessage, _ := json.Marshal(response)

	//TODO: publish to topic would UnComment this
	jst, err := JetStreamInit()
	if err != nil {
		log.Fatal("cant connect to nats service a acquired func ", err.Error())
	}
	fmt.Print("before  publishing", string(byteMessage))
	err = CreateStream(jst)
	if err != nil {
		log.Fatal("cant create stream from acquired ", err.Error())
	}

	if err != nil {
		log.Fatal("cant create stream from acquired ", err.Error())
	}
	//    jst.Publish("targets.acquired",byteMessage)
	//TODO: published bytes Message
	jst.Publish(EventName, byteMessage)
	// defer pullSub()

	return response, nil

}

// Get Single Data from all Datas

func (srv *server) GetSingleTarget(ctx context.Context, req *se.GetSingleTargetRequest) (*se.GetSingleTargetResponse, error) {

	singlecollection, err := ConnectMongo()

	if len(req.Id) == 0 {

		log.Println("please enter a  valid id")
	}
	fmt.Print(req.Id)

	// convert string id (from proto) to mongoDB ObjectId
	oId, err := primitive.ObjectIDFromHex(req.GetId())

	fmt.Print("ID Entered  all IDs", oId)
	if err != nil {
		log.Println("cant convert Ob Id", err.Error())
	}

	SinglePayLoad := &se.GetSingleTargetResponse{}

	if err != nil {
		log.Println("can not connect to mongo db")
	}
	result := singlecollection.FindOne(ctx, bson.M{"_id": oId})

	if err := result.Decode(&SinglePayLoad); err != nil {

		log.Println("can not get data")
	}

	return SinglePayLoad, nil

}

//Get All Data from Data base, all acquired targets
func (srv *server) ListAllTarget(ctx context.Context, req *se.ListAllTargetRequest) (*se.ListAllTargetResponse, error) {

	AllTargetClient, err := ConnectMongo()

	fmt.Print("database created", AllTargetClient.Database())
	fmt.Print("database created", AllTargetClient.Database())

	cur, err := AllTargetClient.Find(context.Background(), bson.M{})

	if err != nil {
		log.Println("returned error from getting data", err.Error())
	}

	TargetSlice := []*se.Data{}
	// var TagSlice2 *se.Data
	for cur.Next(context.Background()) {

		targetslicetoadd := new(se.Data)

		err = cur.Decode(&targetslicetoadd)

		if err != nil {
			log.Println("can not get result")
		}
		// log.Println(string(lt.String()), "///last Convert")
		fmt.Print("\n")
		fmt.Print("------------------------------------------------")
		fmt.Println(targetslicetoadd)

		//  ss := string(lt.String())
		TargetSlice = append(TargetSlice, targetslicetoadd)
	}
	//  fmt.Println(string(T))
	v, _ := time.Now().UTC().MarshalText()
	dateString := string(v)
	fmt.Println("FROM 	UNMARSHALLING")

	return &se.ListAllTargetResponse{
		Id:         uuid.NewString(),
		Targetname: EventName,
		Datas:      TargetSlice,
		Createdat:  dateString,
	}, nil

}

// ListMultipleTarget: acquired multiple target from database, but limited to 3
func (srv *server) ListMultipleTarget(ctx context.Context, req *se.ListMultipleTargetRequest) (*se.ListMultipleTargetResponse, error) {

	AllTargetClient, err := ConnectMongo()

	fmt.Print("database created", AllTargetClient.Database())
	fmt.Print("database created", AllTargetClient.Database())

	  limitStage := bson.D{{"$limit", 3}}
	cur, err := AllTargetClient.Aggregate(context.TODO(), mongo.Pipeline{limitStage})

	if err != nil {
		log.Println("returned error from getting data", err.Error())
	}
TargetSlice := []*se.Data{}
	
      if err := cur.All(context.TODO(), &TargetSlice); err != nil {
          fmt.Println("error from limit")
	   }  

		if err != nil {
			log.Println("can not get result")
		}
		// log.Println(string(lt.String()), "///last Convert")
		fmt.Print("\n")
		fmt.Print("------------------------------------------------")
		fmt.Println(TargetSlice)

	//  convert Data time to string
	v, _ := time.Now().UTC().MarshalText()
	dateString := string(v)
	fmt.Println("FROM 	UNMARSHALLING")

	return &se.ListMultipleTargetResponse{
		Id:         uuid.NewString(),
		Targetname: EventName,
		Datas:      TargetSlice,
		Createdat:  dateString,
	}, nil

}

//Main Entry Point  to grpc Server
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

//Create stream To Send to JetStream
func CreateStream(jetStream nats.JetStreamContext) error {

	stream, err := jetStream.StreamInfo(TargetEvent)
	// stream  not found ,create it

	if stream == nil {
		log.Printf("creating stream: %s\n", TargetEvent)

		_, err = jetStream.AddStream(
			&nats.StreamConfig{
				Name:     TargetEvent,
				Subjects: []string{EventName},
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

	log.Println("connected to Nats from service A")

	if err != nil {
		log.Println("could not publish to Jestream")
		return nil, err
	}
	log.Println("successfully published JetStream")
	return js, nil

}

//Get connection  to Nats
func ConnectToNats() (*nats.Conn, error) {

	nc, err := nats.Connect("nats://Nats:4222")

	if err != nil {
		log.Println("coudl not connect to Nats servicea", err.Error())
	}

	log.Println("connected to nats servicea", nc.ConnectedAddr())

	return nc, nil
}
