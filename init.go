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
var client *mongo.Client

func Init(c Config) error {
	var err error
	config = Config{}
	config.root_context = context.Background()
	config.connection_string, err = connstring.ParseAndValidate(c.ConnectionString)
	if err != nil {
		return err
	}
	config.root_options = options.Client().ApplyURI(c.ConnectionString)

	client, err = mongo.Connect(context.TODO(), config.root_options)
	if err != nil {
		return err
	}

	return nil
}

func CreateContext() (context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(config.root_context, 10*time.Second)
	return ctx, cancel, nil
}

func Ping() error {
	ctx, cancel, err := CreateContext()
	defer client.Disconnect(ctx)
	defer cancel()
	if err != nil {
		return err
	}

	return client.Ping(ctx, readpref.Primary())
}
