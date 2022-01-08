package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/bxcodec/faker/v3"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	numcpu := runtime.NumCPU()
	runtime.GOMAXPROCS(numcpu) // Try to use all available CPUs.
}

const (
	DB_STAGING = "dfl-campaing"
	USERS      = "users"
)

func main() {
	mgo := MongoDbClient()
	ctx := context.Background()

	for i := 0; i < 1000; i++ {
		user := User{
			ID:          primitive.NewObjectID(),
			FirstName:   faker.Name(),
			LastName:    faker.LastName(),
			Email:       faker.Email(),
			Phone:       faker.Phonenumber(),
			JobTitle:    faker.Word(),
			Domain:      faker.DomainName(),
			URL:         faker.URL(),
			PaymentCard: faker.AmountWithCurrency(),
		}
		err := InsertUser(ctx, mgo, user)
		if err != nil {
			log.Fatal(err)
		}
	}
}

type User struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	FirstName   string             `json:"first_name,omitempty"`
	LastName    string             `json:"last_name,omitempty"`
	Email       string             `json:"email,omitempty"`
	Phone       string             `json:"phone,omitempty"`
	JobTitle    string             `json:"job_title,omitempty"`
	Domain      string             `json:"domain,omitempty"`
	URL         string             `json:"url,omitempty"`
	PaymentCard string             `json:"payment_card,omitempty"`
}

//Used to execute client creation procedure only once.
var mongoOnce sync.Once

// MongoDbClient return an unique instance of mongodb server
func MongoDbClient() *mongo.Client {
	var clientInstance *mongo.Client
	//Perform connection creation operation only once.
	mongoOnce.Do(func() {
		// Set client options
		clientOptions := options.Client().ApplyURI(
			os.Getenv("DATABASE_URL"),
		)
		// Connect to MongoDB
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Fatal(fmt.Sprintf("mongodb: %v", err))
			panic(err)
		}
		// Check the connection
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			log.Fatal(fmt.Sprintf("mongodb: %v", err))
			panic(err)
		}
		clientInstance = client
	})
	return clientInstance
}

func InsertUser(
	ctx context.Context,
	m *mongo.Client,
	payload User,
) error {
	var err error

	//Create a handle to the respective collection in the database.
	collection := m.Database(DB_STAGING).Collection(USERS)
	//Perform InsertOne operation & validate against the error.
	_, err = collection.InsertOne(ctx, payload)
	if err != nil {
		return err
	}

	return nil
}
