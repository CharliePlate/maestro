package database

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/charlieplate/maestro"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoProcesser struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
	opts       MongoProcessorOpts
	handler    maestro.QueueHandler
}

type MongoProcessorOpts struct {
	DatabaseName   string
	CollectionName string
}

type MongoChangeEvent struct {
	FullDocument  bson.M             `bson:"fullDocument"`
	OperationType string             `bson:"operationType"`
	DocumentKey   primitive.ObjectID `bson:"documentKey"`
}

func NewMongoWatcher(client *mongo.Client, opts MongoProcessorOpts) (*MongoProcesser, error) {
	if opts.DatabaseName == "" {
		return nil, errors.New("database name is required")
	} else if opts.CollectionName == "" {
		return nil, errors.New("collection name is required")
	}

	mw := &MongoProcesser{
		client:     client,
		opts:       opts,
		database:   nil,
		collection: nil,
	}

	mw.database = client.Database(mw.opts.DatabaseName)
	mw.collection = mw.database.Collection(mw.opts.CollectionName)

	return mw, nil
}

func (mp *MongoProcesser) Watch(ctx context.Context, c chan maestro.QueueUpdateMessage) error {
	watcher, err := mp.collection.Watch(ctx, mongo.Pipeline{})
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
				var op maestro.OpType
				switch data.OperationType {
				case "insert":
					op = maestro.OpTypeInsert
				case "update":
					op = maestro.OpTypeUpdate
				case "delete":
					op = maestro.OpTypeDelete
				default:
					continue
				}

				msg := maestro.QueueUpdateMessage{
					OpType:  op,
					MsgID:   data.DocumentKey.Hex(),
					Content: data.FullDocument,
				}

				c <- msg
			}
		}
	}
}

type MongoHandler struct {
	container maestro.Container
}

func (mh *MongoHandler) Handle(m *maestro.QueueUpdateMessage, errChan chan error) {
	switch m.OpType {
	case maestro.OpTypeInsert:
		mh.container.Push(m)
	case maestro.OpTypeUpdate:
		elem, err := mh.container.Find(m.ID())
		if err != nil {
			errChan <- err
			return
		}

		elem.SetID(m.ID())
		elem.SetData(m.Data())
	case maestro.OpTypeDelete:
		err := mh.container.Delete(m.ID())
		if err != nil {
			errChan <- err
			return
		}
	}
}

func (mh *MongoHandler) SetContainer(c maestro.Container) {
	mh.container = c
}

func (mh *MongoHandler) Container() maestro.Container {
	return mh.container
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
