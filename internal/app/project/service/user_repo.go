package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	log "github.com/sirupsen/logrus"

	"gitlab.com/standard-go/project/internal/app/project"
)

func (r *MongoRepo) FetchIndexUser(data *project.PageRequest) ([]*project.User, int, error) {
	collection := r.Collection(CollectionUsers)

	search := ".*" + data.Search + ".*"
	searchFilter := []bson.D{
		bson.D{{"full_name", primitive.Regex{Pattern: search, Options: "i"}}},
	}

	exclusiveFilter := []bson.D{
		bson.D{{"full_name", primitive.Regex{Pattern: ".*.*", Options: "i"}}},
	}

	if len(data.Filters) > 0 {
		for _, v := range data.Filters {
			if fil := v.ToBson(); fil != nil {
				exclusiveFilter = append(exclusiveFilter, fil)
			}
		}
	}

	filter := bson.M{
		"deleted_at": nil,
		"$or":        searchFilter,
		"$and":       exclusiveFilter,
	}

	count, err := collection.Count(context.TODO(), filter, nil)

	options := options.Find()
	if len(data.Sorts) > 0 {
		for _, v := range data.Sorts {
			if sor := v.ToBson(); sor != nil {
				options.SetSort(sor)
			}
		}
	}

	if data.Paginate == 1 {
		options.SetLimit(data.PerPage)
		if data.Page >= 1 {
			options.SetSkip(data.PerPage * (data.Page - 1))
		}
	}

	fetch, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		log.Fatal(err)
		return nil, 0, errors.New("500")
	}

	defer fetch.Close(nil)
	users := make([]*project.User, 0)
	for fetch.Next(context.TODO()) {
		var elem project.User
		err := fetch.Decode(&elem)

		if err != nil {
			log.Print(err)
			return nil, 0, errors.New("500")
		}
		users = append(users, &elem)
	}

	return users, int(count), nil
}

func (r *MongoRepo) FetchShowUser(id string) (*project.User, error) {
	errNotFound := errors.New(
		fmt.Sprintf("User with ID: %s is not found", id),
	)

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errNotFound
	}

	var user *project.User
	collection := r.Collection(CollectionUsers)
	filter := bson.M{"deleted_at": nil, "_id": oid}

	err = collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, errNotFound
	}

	fmt.Printf("Found a single document: %+v\n", id)

	return user, nil
}

func (r *MongoRepo) FetchStoreUser(data *project.User) (*project.User, error) {
	collection := r.Collection(CollectionUsers)

	now := time.Now()

	data.CreatedAt = now
	data.UpdatedAt = now

	insertResult, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	data.ID = insertResult.InsertedID

	return data, nil
}

func (r *MongoRepo) FetchUpdateUser(id string, data *project.User) (*project.User, error) {
	errNotFound := errors.New(
		fmt.Sprintf("User with ID: %s is not found", id),
	)

	collection := r.Collection(CollectionUsers)
	_, err := r.FetchShowUser(id)
	if err != nil {
		log.Print(err)
		return nil, errors.New("404")
	}

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Print(err)
		return nil, errors.New("404")
	}

	filter := bson.M{
		"deleted_at": nil,
		"_id": bson.M{
			"$ne": oid,
		},
	}

	now := time.Now()

	data.UpdatedAt = now

	filter = bson.M{"deleted_at": nil, "_id": oid}

	_, err = collection.UpdateOne(context.TODO(), filter, bson.M{"$set": data})

	if err != nil {
		return nil, errNotFound
	}

	fmt.Println("Updated a single document: ", id)

	data.ID = id

	return data, nil
}

func (r *MongoRepo) FetchDestroyUser(id string) error {
	errNotFound := errors.New(
		fmt.Sprintf("User with ID: %s is not found", id),
	)

	collection := r.Collection(CollectionUsers)
	data, err := r.FetchShowUser(id)
	if err != nil {
		log.Print(err)
		return errors.New("404")
	}

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Print(err)
		return errors.New("404")
	}

	filter := bson.M{
		"deleted_at": nil,
		"_id": bson.M{
			"$ne": oid,
		},
	}

	now := time.Now()

	data.UpdatedAt = now
	data.DeletedAt = now

	filter = bson.M{"deleted_at": nil, "_id": oid}

	_, err = collection.UpdateOne(context.TODO(), filter, bson.M{"$set": data})
	if err != nil {
		return errNotFound
	}

	fmt.Println("Updated a single document: ", id)

	data.ID = id

	return nil
}
