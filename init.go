package goodm

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

type Config struct {
	ConnectionString string

	root_options      *options.ClientOptions
	root_context      context.Context
	connection_string connstring.ConnString
}

var config Config

func Init(c Config) error {
	var err error
	config = Config{}
	config.root_context = context.Background()
	config.connection_string, err = connstring.ParseAndValidate(c.ConnectionString)
	if err != nil {
		return err
	}
	config.root_options = options.Client().ApplyURI(c.ConnectionString)
	return nil
}

func CreateConnection() (context.Context, *mongo.Client, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(config.root_context, 10*time.Second)
	client, err := mongo.Connect(ctx, config.root_options)
	if err != nil {
		return nil, nil, cancel, err
	}
	return ctx, client, cancel, nil
}

func Ping() error {
	ctx, client, cancel, err := CreateConnection()
	defer client.Disconnect(ctx)
	defer cancel()
	if err != nil {
		return err
	}

	return client.Ping(ctx, readpref.Primary())
}
