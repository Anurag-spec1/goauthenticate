package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var UserCollection *mongo.Collection

func ConnectDB(){
	client, err := mongo.NewClient(options.Client().ApplyURI(GetEnv("MONGO_URI")))
	if err!=nil{
		log.Fatal(err)
	}

	ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()

	err=client.Connect(ctx)
	if err !=nil{
		log.Fatal(err)
	}

	db:=client.Database(GetEnv("DB_NAME"))
	UserCollection=db.Collection("users")

	log.Println("MongoDB connected")
}