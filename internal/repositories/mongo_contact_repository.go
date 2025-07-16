package repositories

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoContactRepository struct {
	collection *mongo.Collection
}

func NewMongoContactRepository(db *mongo.Database) *MongoContactRepository {
	return &MongoContactRepository{
		collection: db.Collection("contacts"),
	}
}

func (r *MongoContactRepository) CreateContact(ctx context.Context, name, email, message string) error {
	contact := Contact{
		Name:      name,
		Email:     email,
		Message:   message,
		CreatedAt: time.Now(),
	}
	_, err := r.collection.InsertOne(ctx, contact)
	return err
}

func (r *MongoContactRepository) GetContacts(ctx context.Context) ([]Contact, error) {
	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var contacts []Contact
	if err := cursor.All(ctx, &contacts); err != nil {
		return nil, err
	}
	return contacts, nil
}

func (r *MongoContactRepository) GetContactByID(ctx context.Context, id string) (Contact, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Contact{}, err
	}
	var contact Contact
	if err := r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&contact); err != nil {
		return Contact{}, err
	}
	return contact, nil
}

func (r *MongoContactRepository) UpdateContact(ctx context.Context, id string, name, email, message string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "name", Value: name},
			{Key: "email", Value: email},
			{Key: "message", Value: message},
			{Key: "updated_at", Value: time.Now()},
		}},
	}
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

func (r *MongoContactRepository) DeleteContact(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}
