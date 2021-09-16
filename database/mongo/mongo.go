package mongo

import (
	"context"
	"fmt"
	"github.com/niuhuan/mirai-bot/config"
	"github.com/niuhuan/mirai-bot/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var (
	mongoClient *mongo.Client
	dbName      string
)

func init() {
	uri := fmt.Sprintf(
		"mongodb://%s:%d/%s?w=majority",
		config.Config.Database.Mongo.Hostname,
		config.Config.Database.Mongo.Port,
		config.Config.Database.Mongo.Database,
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	utils.PanicNotNil(err)
	mongoClient = client
	dbName = config.Config.Database.Mongo.Hostname
}

func Test() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	utils.PanicNotNil(mongoClient.Ping(ctx, nil))
}

func Collection(collectionName string) *mongo.Collection {
	return mongoClient.Database(dbName).Collection(collectionName)
}
