package main

import (
	"context"
	// "crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	// "net"
	// "net/http" 
	"os"
	"time"

	// se "github.com/franklynobleC/microserviceTask/grpcServiceA/destroyer/protos/protos/proto"
	// "github.com/grpc-ecosystem/grpc-gateway/v2/runtime" uncomment later
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
)

const (
	TargetEvent = "TARGETS"
)


type SubScribePayLoad struct {	
    ID         string    `json:"id,omitempty"`
	Message    string    `json:"message,omitempty"`
	Created_on  string `json:"created_on"`
	Updated_on  string `json:"updated_on"`

}


func main() {

	subScribeAndWrite()

// 	grpcMux := runtime.NewServeMux()
// 	ctx, cancel := context.WithCancel(context.Background())

// 	defer cancel()

// 	err := se.RegisterDestroyerServiceHandlerServer(ctx, grpcMux, NewServer())

// 	if err != nil {
// 		log.Fatal("can not register handler Server", err)
// 	}

// 	mux := http.NewServeMux()

// 	mux.Handle("/", grpcMux)

// 	listener, err := net.Listen("tcp", ":5000")

// 	if err != nil {
// 		log.Fatal("can not create listener", err)
// 	}

// 	log.Println("http Gateway Server is being Started", listener.Addr().String())

// 	err = http.Serve(listener, mux)

// 	if err != nil {
// 		log.Fatal("can not start grpc server", err)
// 	}
// }
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

///
//Connet To Mongo Db and Return db Client
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

	wordDictionary := client.Database(os.Getenv("DB_NAME")).Collection("deathstar")

	fmt.Print("database created", wordDictionary.Database())

	return wordDictionary, nil
}

//subscribe NATS to a topic and write to Database
func subScribeAndWrite() {

	// 	//TODO: connect to Database and get Database Client

	wordDictionary, err := ConnectMongo()

	if err != nil {
		log.Fatal("could not connect to mongo Db")
	}
	fmt.Print("database connected successfully")

	fmt.Print("database created", wordDictionary.Database())

	//TODO: NATS CONNECTION
	//subscribe to natsTopic
	nc, err := ConnectToNats()
	if err != nil {
		log.Println("could not connet to nats", err)
	}

	sub, err := nc.SubscribeSync(TargetEvent)

	// nc.InMsgs

	if err != nil {
		log.Print("error subscribing", err)
	}

	//wait for a  message
	//wait for this number of seconds to get the using  this time out
	msg, err := sub.NextMsg(50 * time.Second)

	if err != nil {
		log.Fatal(err.Error())
	}
	var SubM SubScribePayLoad
	//use  the response
	log.Print("from metadata", msg.Subject)

	_ = json.Unmarshal(msg.Data, &SubM)

	if err != nil {
		log.Println("ERROR UNMARSHALLING FROM SERVICE B", err)
	}

	log.Printf("Data: All Details printed %s", SubM)
	// smg := string(msg.Data)

	// msg.Metadata().
    //    SubM.Created_on = SubM.Created_on
	   
	nn := bson.D{{Key: "id", Value: string(SubM.ID)}, {Key: "message", Value: SubM.Message}, {Key: "created_on", Value: SubM.Created_on}, {Key: "updated_on", Value: SubM.Updated_on}}

	if err != nil {
		log.Print("can not unmarshal")
	}

	words, err := wordDictionary.InsertOne(context.TODO(), nn)

	if err != nil {
		log.Print("could not insert data", err.Error())
	}
	//else diplay the id of the newly inserted ID
	fmt.Println(words.InsertedID)

	fil, err := wordDictionary.Find(context.TODO(), nn)

	// Ok := fil.Next(context.TODO())

	defer fil.Close(context.Background())
	// fmt.Println(fil.Next(context.TODO()))
	//  fmt.Print(fil.Decode(fil))

	for fil.Next(context.Background()) {

		result := struct {
			m map[string]string
		}{}

		err := fil.Decode(&result)

		if err != nil {
			log.Fatal(err.Error(), "decoding data")
		}

	}

}
