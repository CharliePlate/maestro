package database_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMongo_CanConnectToTestDB(t *testing.T) {
	ctx := context.Background()

	cli, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/?connect=direct"))
	require.Nil(t, err)

	err = cli.Ping(ctx, nil)
	require.Nil(t, err)
}
