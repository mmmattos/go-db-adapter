package storage

import (
	"context"
	"time"

	"go-db-adapter/main/app/domain"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
)

// MongoDB is a struct that represents a MongoDB datastore.
type MongoDB struct {
	Name       string
	client     *mongo.Client
	database   string
	collection string
}

func (m *MongoDB) AdapterName() string {
	m.Name = "MongoDB"
	return m.Name
}

func (m *MongoDB) Connect(uri, database, collection string) (*mongo.Client, error) {

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return nil, err
	}

	return m.client, nil
}

func (m *MongoDB) Close() error {
	return m.client.Disconnect(context.Background())
}

func (m *MongoDB) FindAll() ([]domain.Person, error) {
	collection := m.client.Database(m.database).Collection(m.collection)

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var Persons []domain.Person
	for cursor.Next(context.Background()) {
		var person domain.Person
		err := cursor.Decode(&person)
		if err != nil {
			return nil, err
		}
		Persons = append(Persons, person)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return Persons, nil
}

func (m *MongoDB) GetByID(id string) (domain.Person, error) {
	collection := m.client.Database(m.database).Collection(m.collection)

	var person domain.Person
	err := collection.FindOne(context.Background(), bson.M{"id": id}).Decode(&person)
	if err != nil {
		return domain.Person{}, err
	}

	return person, nil
}

func (m *MongoDB) Insert(person domain.Person) error {
	ctx := context.Background()
	collection := m.client.Database(m.database).Collection(m.collection)

	err := m.WithTransaction(ctx, func(sessCtx mongo.SessionContext) error {
		_, err := collection.InsertOne(sessCtx, person)
		return err
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoDB) GetSession() (mongo.Session, error) {
	return m.client.StartSession()
}

func (m *MongoDB) WithTransaction(ctx context.Context, fn func(sessionContext mongo.SessionContext) error) error {
	session, err := m.GetSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	dur := 5 * time.Second
	durMemAddr := &dur

	opts := options.Transaction().SetMaxCommitTime(durMemAddr)
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (any, error) {
		return fn(sessCtx), nil // Pass the session context to the provided function
	}, opts)
	if err != nil {
		return err
	}

	return nil
}

//============================== ADAPTER STARTS HERE

// MongoDBAdapter is a struct that represents a MongoDBAdapter datastore adapter.
type MongoDBAdapter struct {
	mongodb *MongoDB
}

func (a *MongoDBAdapter) AdapterName() string {
	return a.mongodb.AdapterName()
}

func (a *MongoDBAdapter) Connect(uri, database, collection string) (*mongo.Client, error) {
	return a.mongodb.Connect(uri, database, collection)
}

func (a *MongoDBAdapter) Close() error {
	return a.mongodb.Close()
}

func (a *MongoDBAdapter) GetSession() (mongo.Session, error) {
	return a.mongodb.GetSession()
}

func (a *MongoDBAdapter) WithTransaction(ctx context.Context, fn func(sessionContext mongo.SessionContext) error) error {
	return a.mongodb.WithTransaction(ctx, fn)
}

func (a *MongoDBAdapter) FindAll() ([]domain.Person, error) {
	return a.mongodb.FindAll()
}

func (a *MongoDBAdapter) GetByID(id string) (domain.Person, error) {
	return a.mongodb.GetByID(id)
}

func (a *MongoDBAdapter) Insert(person domain.Person) error {
	return a.mongodb.Insert(person)
}
