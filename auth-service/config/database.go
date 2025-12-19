package config

import (
    "context"
    "fmt"
    "log"
    "os"
    "sync"
    "time"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"
)

var (
    DB             *mongo.Database
    UserCollection *mongo.Collection
    client         *mongo.Client
    once           sync.Once
)

func ConnectDB() {
    once.Do(func() {
        mongoURI := os.Getenv("MONGO_URI")
        dbName := os.Getenv("DB_NAME")

        if mongoURI == "" {
            log.Fatal("MONGO_URI environment variable is not set")
        }
        if dbName == "" {
            dbName = "auth_db"
        }

        // Create client with options
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()

        clientOptions := options.Client().
            ApplyURI(mongoURI).
            SetMaxPoolSize(100).
            SetMinPoolSize(5).
            SetMaxConnIdleTime(30 * time.Second)

        var err error
        client, err = mongo.Connect(ctx, clientOptions)
        if err != nil {
            log.Fatal("Failed to connect to MongoDB:", err)
        }

        // Check connection
        err = client.Ping(ctx, nil)
        if err != nil {
            log.Fatal("Failed to ping MongoDB:", err)
        }

        fmt.Println("✅ Connected to MongoDB!")
        
        DB = client.Database(dbName)
        UserCollection = DB.Collection("users")

        // Create indexes
        createIndexes()
    })
}

func DisconnectDB() {
    if client != nil {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        
        if err := client.Disconnect(ctx); err != nil {
            log.Println("Error disconnecting from MongoDB:", err)
        } else {
            fmt.Println("✅ Disconnected from MongoDB")
        }
    }
}

func createIndexes() {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Create unique index on email
    emailIndex := mongo.IndexModel{
        Keys:    bson.D{{Key: "email", Value: 1}},
        Options: options.Index().SetUnique(true).SetName("unique_email"),
    }

    // Create TTL index for OTP expiration
    ttlIndex := mongo.IndexModel{
        Keys:    bson.D{{Key: "otp_expires_at", Value: 1}},
        Options: options.Index().SetExpireAfterSeconds(600).SetName("otp_expiry"), // 10 minutes
    }

    // Create indexes
    _, err := UserCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{emailIndex, ttlIndex})
    if err != nil {
        log.Printf("Warning: Could not create indexes: %v\n", err)
    } else {
        fmt.Println("✅ Database indexes created")
    }
}

func GetClient() *mongo.Client {
    return client
}