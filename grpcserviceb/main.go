package main

import (
	"context"
	// "sync"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	TargetEvent = "TARGET"
	EventName   = "targets.acquired"
)

//payload to subscribe to
type SubScribePayLoad struct {
	ID        string `json:"id"`
	Message   string `json:"message"`
	Createdat string `json:"createdat"`
	Updatedat string `json:"updatedat"`
}

func main() {
	for {
	 ConnectToNats()

	subScribeAndWrite()
}
	

}

func ConnectToNats() (*nats.Conn, error) {

	nc, err := nats.Connect("nats://Nats:4222")

	if err != nil {
		log.Println("coudl not connect to Nats", err.Error())
	}

	log.Println("connected to Nats from Service B", nc.ConnectedAddr())

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

	// wg := sync.WaitGroup{}
	// wg.Add(1)
	//TODO: connect to Database and get Database Client\

	deathstarCollection, err := ConnectMongo()

	if err != nil {
		log.Fatal("could not connect to mongo Db")
	}
	fmt.Print("database connected successfully")

	fmt.Print("database created", deathstarCollection.Database())

	//TODO: NATS CONNECTION
	//subscribe to natsTopic
	nc, err := ConnectToNats()
	// this would wait for incomming message from Name "targets.acquired" in order as  they arrive

	sub, _ := nc.SubscribeSync(EventName)
	// wg.Done()
	fmt.Print("events Delivered")

	if err != nil {
		log.Println("could not connet to nats", err)
	} 
	for {

	msg, err := sub.NextMsg(100 * time.Second)

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

	target, err := deathstarCollection.InsertOne(context.TODO(), nn)

	if err != nil {
		log.Print("could not insert data", err.Error())
	}
	//else diplay the id of the newly inserted ID
	fmt.Println(target.InsertedID)
}

}
