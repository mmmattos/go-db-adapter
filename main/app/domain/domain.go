package domain

import (
	"database/sql"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Person struct {
	ID   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name,omitempty" bson:"name,omitempty"`
	Age  int                `json:"age,omitempty" bson:"age,omitempty"`
}

type Service struct {
	// Any dependencies or database adapters
}

func NewService() *Service {
	// Initialize any dependencies here
	return &Service{}
}

func (s *Service) CreatePerson(person Person) error {
	// Implement the logic to create a person in your database
	return nil
}

func (s *Service) GetPerson(id string) (Person, error) {
	// Implement the logic to retrieve persons from your database
	person := Person{}
	return person, nil
}

func (s *Service) GetPersons() ([]Person, error) {
	// Implement the logic to retrieve persons from your database
	return nil, nil
}

func NewPersonFromRow(row *sql.Rows) (*Person, error) {
	person := &Person{}
	if row.Next() {
		err := row.Scan(&person.ID, &person.Name, &person.Age)
		if err != nil {
			return nil, err
		}
	}
	return person, nil
}
