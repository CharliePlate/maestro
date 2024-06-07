package maestro

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoWatcher struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
	opts       MongoWatcherOpts
}

type MongoWatcherOpts struct {
	DatabaseName   string
	CollectionName string
}

type MongoChangeEvent struct {
	FullDocument  bson.M             `bson:"fullDocument"`
	OperationType string             `bson:"operationType"`
	DocumentKey   primitive.ObjectID `bson:"documentKey"`
}

func NewMongoWatcher(client *mongo.Client, opts MongoWatcherOpts) (*MongoWatcher, error) {
	if opts.DatabaseName == "" {
		return nil, errors.New("database name is required")
	} else if opts.CollectionName == "" {
		return nil, errors.New("collection name is required")
	}

	mw := &MongoWatcher{
		client:     client,
		opts:       opts,
		database:   nil,
		collection: nil,
	}

	mw.database = client.Database(mw.opts.DatabaseName)
	mw.collection = mw.database.Collection(mw.opts.CollectionName)

	return mw, nil
}

func (mw *MongoWatcher) Watch(ctx context.Context, c chan QueueUpdateMessage) error {
	watcher, err := mw.collection.Watch(ctx, mongo.Pipeline{})
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if watcher.Next(ctx) {
				var data MongoChangeEvent
				err := watcher.Decode(&data)
				if err != nil {
					return err
				}
				var op OpType
				switch data.OperationType {
				case "insert":
					op = OpTypeInsert
				case "update":
					op = OpTypeUpdate
				case "delete":
					op = OpTypeDelete
				default:
					continue
				}

				msg := QueueUpdateMessage{
					OpType: op,
					ID:     data.DocumentKey.Hex(),
					Data:   data.FullDocument,
				}

				c <- msg
			}
		}
	}
}

// MongoAuthURL constructs a MongoDB connection string from environment variables.
//
//nolint:nosprintfhostport // Protocol prefix required for MongoDB connection string
func MongoAuthURL() (string, error) {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		return "", err
	}

	username := os.Getenv("MONGO_USERNAME")
	password := os.Getenv("MONGO_PASSWORD")
	host := os.Getenv("MONGO_HOST")
	port := os.Getenv("MONGO_PORT")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "27017"
	}

	if (username == "" && password != "") || (username != "" && password == "") {
		return "", errors.New("both username and password must be provided together")
	}

	if username != "" && password != "" {
		return fmt.Sprintf("mongodb://%s:%s@%s:%s/?connect=direct", username, password, host, port), nil
	}
	return fmt.Sprintf("mongodb://%s:%s/?connect=direct", host, port), nil
}
