package weather_subscription

import (
	"context"
	"git.foxminded.ua/foxstudent104911/2.1about-me-bot/telegram"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

type MongoSubscriptionConnection struct {
	usersCollection mongo.Collection
}

func NewMongoSubscriptionConnection(usersCollection mongo.Collection) *MongoSubscriptionConnection {
	return &MongoSubscriptionConnection{usersCollection: usersCollection}
}

func (u *MongoSubscriptionConnection) CreateUser(callback telegram.Callback) error {
	subscription := bson.D{
		{"callback_id", callback.CallbackId},
		{"username", callback.User.UserName},
		{"send_at", callback.Button},
		{"created_at", time.Now()},
		{"updated_at", time.Now()}}

	insertRes, err := u.usersCollection.InsertOne(context.TODO(), subscription)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user": callback.User.UserName,
		})

		return err
	}

	logrus.WithFields(logrus.Fields{
		"user": callback.User.UserName,
	}).Info("user inserted", insertRes)

	return nil
}

func (u *MongoSubscriptionConnection) UpdateSubscriptionWithLocation(message telegram.Message) error {
	var err error
	usernameInBson := bson.D{{"username", message.Username.UserName}}

	update := bson.D{
		{"$set", bson.D{
			{"username", message.Username.UserName},
			{"chat_id", message.ChatId},
			{"longitude", message.Location.Longitude},
			{"latitude", message.Location.Latitude},
			{"updated_at", time.Now()},
		},
		},
	}

	_, err = u.usersCollection.UpdateOne(context.TODO(), usernameInBson, update)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"username": message.Username.UserName,
			"error":    err,
		}).Error(err)
		return err
	}

	return nil
}

func (u *MongoSubscriptionConnection) UpdateSubscriptionWithTime(callback telegram.Callback, pressedButton string) error {
	var result Subscription
	usernameInBson := bson.D{{"username", callback.User.UserName}}

	update := bson.D{
		{"$set", bson.D{
			{"callback_id", callback.CallbackId},
			{"send_at", pressedButton},
			{"updated_at", time.Now()}},
		},
	}

	err := u.usersCollection.FindOne(context.TODO(), usernameInBson).Decode(&result)

	if result.Username == "" {
		err = u.CreateUser(callback)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"user": result.Username,
			}).Error(err)
			return err
		}
	}

	_, err = u.usersCollection.UpdateByID(context.TODO(), result.ID, update)

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

		}

		listOfSubscriptions = append(listOfSubscriptions, result)
	}

	return listOfSubscriptions, nil
}
