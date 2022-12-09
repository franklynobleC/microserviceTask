package main

import (
	"context"
	// "strings"
	// "crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	// "net"
	// "net/http"
	"os"
	// "time"

	// se "github.com/franklynobleC/microserviceTask/grpcServiceA/destroyer/protos/protos/proto"
	// "github.com/grpc-ecosystem/grpc-gateway/v2/runtime" uncomment later
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "google.golang.org/protobuf/types/known/timestamppb"
	// "go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
)

const (
	TargetEvent = "TARGET"
)

type SubScribePayLoad struct {
	ID        string `json:"id"`
	Message   string `json:"message"`
	Createdat string `json:"createdat"`
	Updatedat string `json:"updatedat"`
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
// func subScribeAndWrite() {

// 	// 	//TODO: connect to Database and get Database Client

// 	deathstarCollection, err := ConnectMongo()

// 	if err != nil {
// 		log.Fatal("could not connect to mongo Db")
// 	}
// 	fmt.Print("database connected successfully")

// 	fmt.Print("database created", deathstarCollection.Database())

// 	//TODO: NATS CONNECTION
// 	//subscribe to natsTopic
// 	nc, err := ConnectToNats()
// 	if err != nil {
// 		log.Println("could not connet to nats", err)
// 	}

// 	sub, err := nc.SubscribeSync(Target)

// 	// nc.InMsgs

// 	if err != nil {
// 		log.Print("error subscribing", err)
// 	}

// 	//wait for a  message
// 	//wait for this number of seconds to get the using  this time out
// 	msg, err := sub.NextMsg(50 * time.Second)

// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}
// 	var SubM SubScribePayLoad
// 	//use  the response
// 	log.Print("from metadata", msg.Subject)

// 	_ = json.Unmarshal(msg.Data, &SubM)

// 	if err != nil {
// 		log.Println("ERROR UNMARSHALLING FROM SERVICE B", err)
// 	}

// 	log.Println("Data: All Details printed", SubM)
// 	// smg := string(msg.Data)

// 	// msg.Metadata().
// 	//    SubM.Created_on = SubM.Created_on

// 	// enter valur to write to datatbase
// 	// id, _ := primitive.ObjectIDFromHex(SubM.Updated_on.String())
// 	// id1, _ := primitive.ObjectIDFromHex(SubM.Created_on.String())
// 	nn := bson.D{{Key: "id", Value: (SubM.ID)}, {Key: "message", Value: SubM.Message}, {Key: "created_on", Value: SubM.Created_on}, {Key: "updated_on", Value: SubM.Updated_on}}

// 	if err != nil {
// 		log.Print("can not unmarshal")
// 	}

// 	words, err := deathstarCollection.InsertOne(context.TODO(), nn)

// 	if err != nil {
// 		log.Print("could not insert data", err.Error())
// 	}
// 	//else diplay the id of the newly inserted ID
// 	fmt.Println(words.InsertedID)

// 	fil, err := deathstarCollection.Find(context.TODO(), nn)

// 	defer fil.Close(context.Background())

// 	for fil.Next(context.Background()) {

// 		result := struct {
// 			m map[string]string
// 		}{}

// 		err := fil.Decode(&result)

// 		if err != nil {
// 			log.Fatal(err.Error(), "decoding data")
// 		}

// 	}

// }

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

	 nc.Subscribe("events.targets", func(msg *nats.Msg) {

	 
        fmt.Print("events Delivered")
		// fmt.Print(string(msg.Data))
		//  fmt.Print(msg.Data)
	
	if err != nil {
		log.Println("could not connet to nats", err)
	}


	
	
	// fmt.Print(sub.Subject)
	//     fmt.Print(sub.ConsumerInfo())

	// /.InMsgs
	if err != nil {
		log.Print("error subscribing", err)
	}

	//wait for a  message
	//wait for this number of seconds to get the using  this time out
	// msg, err := sub.NextMsg(50 * time.Second)

	if err != nil {
		log.Fatal(err.Error())
	}
	var SubM SubScribePayLoad
	//use  the response
	log.Print("from metadata", msg.Subject)
	fmt.Print("Before marshaling", msg.Data)
      
	err = json.Unmarshal(msg.Data, &SubM)
	fmt.Print("After UMarshalling", SubM)
	if err != nil {
		log.Println("ERROR UNMARSHALLING FROM SERVICE B", err.Error())
	}


	log.Printf("Data: All Details printed %s", SubM)

	fmt.Print("Before Writing To Db", SubM)
	nn := bson.D{{Key: "id", Value: SubM.ID}, {Key: "message", Value: SubM.Message}, {Key: "createdat", Value: SubM.Createdat}, {Key: "updatedat", Value: SubM.Updatedat}}

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
})

}
