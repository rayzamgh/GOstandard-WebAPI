package service

import (
	"context"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
)

// Collection names
const (
	CollectionUsers = "users"
)

// MongoRepo implements facultative.Repository using MongoDB
type MongoRepo struct {
	client *mongo.Client
	db     *mongo.Database
}

// NewMongoRepo instantiates a new MongoRepo instance
func NewMongoRepo(url, dbName string) (*MongoRepo, error) {
	client, err := mongo.Connect(context.Background(), url)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return nil, err
	}

	return &MongoRepo{
		client: client,
		db:     client.Database(dbName),
	}, nil
}

// Collection returns collection by the given name
func (r *MongoRepo) Collection(name string) *mongo.Collection {
	return r.db.Collection(name)
}

// Database returns database :v
func (r *MongoRepo) Database() *mongo.Database {
	return r.db
}

// Close disconnects db connection
func (r *MongoRepo) Close() error {
	return r.client.Disconnect(context.Background())
}

// Client returns its *mongo.Client instance
func (r *MongoRepo) Client() *mongo.Client {
	return r.client
}

// Service implements project.Service interface
type Service struct {
	Repo *MongoRepo
}

// New instantiates new Service instance
func New(mongoURL, dbName string) (*Service, error) {
	repo, err := NewMongoRepo(mongoURL, dbName)
	if err != nil {
		return nil, err
	}

	return &Service{Repo: repo}, nil
}

// Close db connection
func (s *Service) Close() error {
	return s.Repo.Close()
}
