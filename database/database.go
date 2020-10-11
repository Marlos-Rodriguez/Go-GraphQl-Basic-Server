package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Marlos-Rodriguez/Go_MongoDB_GraphQL/graph/model"
)

type DB struct {
	client *mongo.Client
}

func Connect() *DB {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)

	return &DB{
		client: client,
	}
}

//Save ...
func (db *DB) Save(input model.NewDog) *model.Dog {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.client.Database("animals").Collection("dogs")

	res, err := collection.InsertOne(ctx, input)

	if err != nil {
		log.Fatal(err)
	}

	return &model.Dog{
		ID:        res.InsertedID.(primitive.ObjectID).Hex(),
		Name:      input.Name,
		IsGoodBoy: input.IsGoodBoy,
	}
}

func (db *DB) FindByID(ID string) *model.Dog {
	ObjID, err := primitive.ObjectIDFromHex(ID)

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.client.Database("animals").Collection("dogs")

	res := collection.FindOne(ctx, bson.M{"_id": ObjID})

	if err != nil {
		log.Fatal(err)
	}

	dog := model.Dog{}

	res.Decode(&dog)

	return &dog
}

func (db *DB) All() []*model.Dog {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection := db.client.Database("animals").Collection("dogs")

	cur, err := collection.Find(ctx, bson.D{})

	if err != nil {
		log.Fatal(err)
	}

	var dogs []*model.Dog

	for cur.Next(ctx) {
		var dog *model.Dog

		err := cur.Decode(&dog)
		if err != nil {
			log.Fatal(err)
		}

		dogs = append(dogs, dog)

	}
	return dogs
}
