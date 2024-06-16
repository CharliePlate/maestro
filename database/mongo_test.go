package database_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMongo_CanConnectToTestDB(t *testing.T) {
	ctx := context.Background()

	connTimeout := 10 * time.Second
	connCtx, connCancel := context.WithTimeout(ctx, connTimeout)
	defer connCancel()

	cli, err := mongo.Connect(connCtx, options.Client().ApplyURI("mongodb://localhost:27017/?connect=direct"))
	require.NoError(t, err)
	defer func() {
		discErr := cli.Disconnect(ctx)
		require.NoError(t, discErr)
	}()

	pingTimeout := 2 * time.Second
	pingCtx, pingCancel := context.WithTimeout(ctx, pingTimeout)
	defer pingCancel()

	err = cli.Ping(pingCtx, nil)
	require.NoError(t, err)

	opTimeout := 5 * time.Second
	opCtx, opCancel := context.WithTimeout(ctx, opTimeout)
	defer opCancel()

	db := cli.Database("TestDB")
	coll := db.Collection("TestCollection")

	_, err = coll.InsertOne(opCtx, bson.M{"name": "pi", "value": 3.14159})
	require.NoError(t, err)

	elem := coll.FindOne(opCtx, bson.M{"name": "pi"})
	require.NotNil(t, elem)

	var result bson.M
	err = elem.Decode(&result)
	require.NoError(t, err)

	require.Equal(t, "pi", result["name"])
}
