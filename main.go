package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/bxcodec/faker/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	numcpu := runtime.NumCPU()
	runtime.GOMAXPROCS(numcpu)
}

const (
	DB_STAGING = "dfl-campaing"
	USERS      = "users"
)

// User ...
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

func main() {
	r := http.NewServeMux()
	h := HandlerApp{m: MongoDbClient()}

	r.HandleFunc("/users/create", h.CreaeteUser)
	r.HandleFunc("/users/all", h.ReadUsers)

	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%s", os.Getenv("PORT")),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

// Handlers
// HandlerApp
type HandlerApp struct {
	m *mongo.Client
}

// CreaeteUser ...
func (h *HandlerApp) CreaeteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		reqBody, _ := ioutil.ReadAll(r.Body)
		var user User
		json.Unmarshal(reqBody, &user)

		user.ID = primitive.NewObjectID()
		user.PaymentCard = faker.AmountWithCurrency()
		user.URL = faker.URL()
		user.Domain = faker.DomainName()

		err := h.InsertUser(r.Context(), user)
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
	}
}

// ReadUsers ...
func (h *HandlerApp) ReadUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		users, err := h.FetchAllUser(r.Context())
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		err = json.NewEncoder(w).Encode(users)
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
	}
}

// DB
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

// InsertUser ...
func (h *HandlerApp) InsertUser(
	ctx context.Context,
	payload User,
) error {
	var err error

	//Create a handle to the respective collection in the database.
	collection := h.m.Database(DB_STAGING).Collection(USERS)
	//Perform InsertOne operation & validate against the error.
	_, err = collection.InsertOne(ctx, payload)
	if err != nil {
		return err
	}

	return nil
}

// FetchAllUser ...
func (h *HandlerApp) FetchAllUser(ctx context.Context) ([]map[string]interface{}, error) {
	filter := bson.D{{}} //bson.D{{}} specifies 'all documents'

	collection := h.m.Database(DB_STAGING).Collection(USERS)
	//Perform Find operation & validate against the error.
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	users := make([]map[string]interface{}, cur.RemainingBatchLength())
	if err := cur.All(ctx, &users); err != nil {
		panic(err)
	}

	return users, nil

}
