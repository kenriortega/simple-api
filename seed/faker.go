package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/bxcodec/faker/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	numcpu := runtime.NumCPU()
	runtime.GOMAXPROCS(numcpu) // Try to use all available CPUs.
}

func main() {

	f, err := os.Create("users.json")
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	for i := 0; i < 10_001; i++ {
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
		fmt.Fprintln(f, user.ToJSON())
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("file was written line by line successfully")
}

type User struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	FirstName   string             `json:"first_name,omitempty"`
	LastName    string             `json:"last_name,omitempty"`
	Email       string             `json:"email,omitempty"`
	Phone       string             `json:"phone,omitempty"`
	JobTitle    string             `json:"job_title,omitempty"`
	Domain      string             `json:"domain,omitempty"`
	URL         string             `json:"url,omitempty"`
	PaymentCard string             `json:"payment_card,omitempty"`
}

// ToJSON ...
func (u *User) ToJSON() string {
	bytes, err := json.Marshal(u)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return string(bytes)
}
