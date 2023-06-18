package weather_subscription

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Location struct {
	Longitude float64 `bson:"longitude"`
	Latitude  float64 `bson:"latitude"`
}

type Subscription struct {
	ID        primitive.ObjectID `bson:"_id"`
	ChatId    int64              `bson:"chat_id"`
	Username  string             `bson:"username"`
	SendAt    string             `bson:"send_at"`
	Location  Location           `bson:",inline"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type Message struct {
	ChatId    int64
	Longitude float64
	Latitude  float64
	Username  string
}

type MongoSubscriptionConnection struct {
	usersCollection mongo.Collection
}

func NewMongoSubscriptionConnection(usersCollection mongo.Collection) *MongoSubscriptionConnection {
	return &MongoSubscriptionConnection{usersCollection: usersCollection}
}

func (u *MongoSubscriptionConnection) UpdateSubscriptionWithLocation(message Message) error {
	var err error
	usernameInBson := bson.D{{"username", message.Username}}

	update := bson.D{
		{"$set", bson.D{
			{"chat_id", message.ChatId},
			{"longitude", message.Longitude},
			{"latitude", message.Latitude},
			{"updated_at", time.Now()},
		},
		},
	}

	_, err = u.usersCollection.UpdateOne(context.TODO(), usernameInBson, update)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"username": message.Username,
			"error":    err,
		}).Error(err)
		return err
	}

	return nil
}

func (u *MongoSubscriptionConnection) UpdateSubscriptionWithTime(username string, pressedButton string) error {
	usernameInBson := bson.D{{"username", username}}
	opts := options.Update().SetUpsert(true)

	subscription := bson.D{
		{"$set", bson.D{
			{"username", username},
			{"send_at", pressedButton},
			{"created_at", time.Now()},
			{"updated_at", time.Now()}}}}

	_, err := u.usersCollection.UpdateOne(context.TODO(), usernameInBson, subscription, opts)

	if err != nil {
		return err
	}

	return nil
}

func (u *MongoSubscriptionConnection) GetSubscriptionDataFromMongoDB() ([]Subscription, error) {
	var listOfSubscriptions []Subscription

	cursor, err := u.usersCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		logrus.Error(err)
		return []Subscription{}, err
	}

	logrus.Info("Users collection was found")

	for cursor.Next(context.TODO()) {

		result := Subscription{}

		if err := cursor.Decode(&result); err != nil {
			logrus.Error(err)
			return nil, err
		}

		listOfSubscriptions = append(listOfSubscriptions, result)
	}

	return listOfSubscriptions, nil
}
