package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	numcpu := runtime.NumCPU()
	runtime.GOMAXPROCS(numcpu)
}

func main() {
	r := http.NewServeMux()
	api := New(r)
	api.routes()
	api.Start()
}

type server struct {
	db     *mongo.Client
	router *http.ServeMux
}

func New(r *http.ServeMux) *server {
	return &server{
		db:     nil,
		router: r,
	}
}

func (s *server) Start() {
	srv := &http.Server{
		Handler: s.router,
		Addr:    fmt.Sprintf(":%s", os.Getenv("PORT")),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
	log.Println("Server running on 0.0.0.0:" + os.Getenv("PORT"))
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
func (s *server) routes() {
	s.router.HandleFunc("/all", s.handleUsers())
}

func (s *server) handleUsers() http.HandlerFunc {

	var (
		init sync.Once

		DB    = "challenge"
		USERS = "users"
	)

	//Perform connection creation operation only once.
	init.Do(func() {
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
		s.db = client
	})

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
	// HandleFunc ...
	return func(w http.ResponseWriter, r *http.Request) {
		filter := bson.D{{}}

		collection := s.db.Database(DB).Collection(USERS)
		cur, err := collection.Find(r.Context(), filter)
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		users := make([]map[string]interface{}, cur.RemainingBatchLength())
		if err := cur.All(r.Context(), &users); err != nil {
			fmt.Fprintf(w, err.Error())
		}

		if r.Method == "GET" {

			err = json.NewEncoder(w).Encode(users)
			if err != nil {
				fmt.Fprintf(w, err.Error())
			}
		}
	}
}
