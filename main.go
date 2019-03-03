package main

import (
	"log"
	"net/http"

	"github.com/emoji/internal/mongodb"
	"github.com/emoji/keys"

	"github.com/emoji/internal/votebot"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	client *mongo.Client
	emojis *mongo.Collection
)

const (
	emojiDatabase   = "emoji"
	emojiCollection = "emoji"
)

func main() {
	client, err := mongodb.NewClient(
		mongodb.UseMongoURI(keys.MongoDBLink),
		mongodb.UseDatabase(emojiDatabase),
		mongodb.UseCollection(emojiCollection),
	)

	if err != nil {
		log.Fatal(err)
	}

	server, err := votebot.NewServer(
		votebot.UseClient(client),
	)

	if err != nil {
		log.Fatal(err)
	}

	http.ListenAndServe(":5000", server.Route())
}
