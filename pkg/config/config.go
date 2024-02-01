package config

import (
	"context"
	"fmt"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HTTPInfo struct {
	Addr string
	Port int
}

type DBInfo struct {
	Client *mongo.Client
	DBName string
}

type AppConfig struct {
	HTTPInfo    *HTTPInfo
	MongoDBInfo *DBInfo
}

func LoadConfig() *AppConfig {
	//port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	port, err := strconv.Atoi("8003")
	if err != nil {
		log.Fatal(err)
	}
	//addr := ":" + os.Getenv("SERVER_PORT")
	addr := ":" + "8003"
	httpInfo := &HTTPInfo{
		Addr: addr,
		Port: port,
	}

	client := connectMongoDb()
	//dbName := os.Getenv("DATABASE_NAME")
	dbName := "web-scraping"

	dbInfo := &DBInfo{
		Client: client,
		DBName: dbName,
	}

	conf := AppConfig{
		MongoDBInfo: dbInfo,
		HTTPInfo:    httpInfo,
	}

	return &conf
}

func connectMongoDb() *mongo.Client {
	// mongoHost := os.Getenv("DATABASE_HOST")
	// dbName := os.Getenv("DATABASE_NAME")
	mongoHost := "mongodb://127.0.0.1:27017/"
	dbName := "web-scraping"
	url := fmt.Sprintf("%s,%s", mongoHost, dbName)

	clientOptions := options.Client().ApplyURI(url).SetMaxPoolSize(100).SetServerSelectionTimeout(10 * time.Second).SetSocketTimeout(10 * time.Second).SetHeartbeatInterval(5 * time.Second)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}

	defer cancel()

	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(ctx, nil)

	if err != nil {
		log.Fatal(err)
	}

	conn := client.NumberSessionsInProgress()
	msg := fmt.Sprintf("MongoClient connected - Sessions actives: %d", conn)
	log.Info(msg)

	return client
}

func (a *AppConfig) CloseMongoDB(ctx context.Context) {
	a.MongoDBInfo.Client.Disconnect(ctx)
	conn := a.MongoDBInfo.Client.NumberSessionsInProgress()
	msg := fmt.Sprintf("MongoClient disconnected - Sessions actives: %d", conn)
	log.Info(msg)
}
