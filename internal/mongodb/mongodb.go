package mongodb

import (
	"context"
	"time"

	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Emoji is the data format stored in MongoDB
type Emoji struct {
	EmojiID string
	Vote    int
}

// Client stores the configuration
type Client struct {
	client         *mongo.Client
	collection     *mongo.Collection
	collectionName string
	databaseName   string
	uri            string
}

// UseMongoURI sets the URI of mongodb
func UseMongoURI(uri string) func(*Client) error {
	return func(c *Client) error {
		c.uri = uri
		return nil
	}
}

// UseDatabase sets the database
func UseDatabase(db string) func(*Client) error {
	return func(c *Client) error {
		c.databaseName = db
		return nil
	}
}

// UseCollection sets the collection
func UseCollection(collection string) func(*Client) error {
	return func(c *Client) error {
		c.collectionName = collection
		return nil
	}
}

// NewClient creates a client to MongoDB
func NewClient(opts ...func(*Client) error) (*Client, error) {
	c := &Client{}

	for _, f := range opts {
		if err := f(c); err != nil {
			return nil, err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.NewClient(options.Client().ApplyURI(c.uri))
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	c.client = client
	c.collection = c.client.Database(c.databaseName).Collection(c.collectionName)
	return c, nil
}

// AddEmoji adds a new emoji
func (c *Client) AddEmoji(emoji *Emoji) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := c.collection.InsertOne(ctx, emoji)
	if err != nil {
		return err
	}

	return nil
}

// GetEmoji gets an item
func (c *Client) GetEmoji(emojiID string) (*Emoji, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"emojiid": emojiID}

	var emoji Emoji
	err := c.collection.FindOne(ctx, filter).Decode(&emoji)
	if err != nil {
		return nil, err
	}

	return &emoji, nil
}

// UpdateEmoji updates info of an emoji
func (c *Client) UpdateEmoji(e *Emoji) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"emojiid": e.EmojiID}

	_, err := c.collection.UpdateOne(ctx, filter, bson.M{
		"$set": bson.M{"vote": e.Vote + 1},
	})

	if err != nil {
		return err
	}

	return nil
}
