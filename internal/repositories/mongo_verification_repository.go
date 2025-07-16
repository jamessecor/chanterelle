package repositories

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoVerificationRepository struct {
	collection *mongo.Collection
}

func NewMongoVerificationRepository(db *mongo.Database) *MongoVerificationRepository {
	return &MongoVerificationRepository{
		collection: db.Collection("verification_codes"),
	}
}

func (r *MongoVerificationRepository) CreateVerificationCode(ctx context.Context, email string, code string, expiry time.Duration) error {
	verificationCode := VerificationCode{
		ID:        primitive.NewObjectID().Hex(),
		Code:      code,
		Email:     email,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(expiry),
	}
	_, err := r.collection.InsertOne(ctx, verificationCode)
	return err
}

func (r *MongoVerificationRepository) GetCodeByEmail(ctx context.Context, email string) (*VerificationCode, error) {
	// Find the latest verification code for the email
	filter := bson.D{{Key: "email", Value: email}}
	sort := bson.D{{Key: "created_at", Value: -1}}

	var verificationCode VerificationCode
	if err := r.collection.FindOne(ctx, filter, options.FindOne().SetSort(sort)).Decode(&verificationCode); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("verification code not found")
		}
		return nil, err
	}

	return &verificationCode, nil
}

func (r *MongoVerificationRepository) DeleteCodeByEmail(ctx context.Context, email string) error {
	filter := bson.D{{Key: "email", Value: email}}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

func (r *MongoVerificationRepository) DeleteExpiredCodes(ctx context.Context) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{"expires_at": bson.M{"$lt": time.Now()}})
	return err
}
