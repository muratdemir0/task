package store

import (
	"context"
	"github.com/muratdemir0/faceit-task/internal/config"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserStore struct {
	client *mongo.Client
	config *config.Mongo
}
type User struct {
	ID        string `json:"id" bson:"_id"`
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	Nickname  string `json:"nickname" bson:"nickname"`
	Password  string `json:"password" json:"password"`
	Email     string `json:"email" bson:"email"`
	Country   string `json:"country" bson:"country"`
}
type ListCriteria struct {
	Country string `json:"country" bson:"country"`
	Page    int64  `json:"page" bson:"page"`
	Limit   int64  `json:"limit" bson:"limit"`
}

func NewUserStore(client *mongo.Client, config *config.Mongo) UserStore {
	return UserStore{client: client, config: config}
}

func (s UserStore) Create(ctx context.Context, cu *User) error {
	_, err := s.client.
		Database(s.config.Name).
		Collection(s.config.Collections.Users).
		InsertOne(ctx, cu)

	if err != nil {
		return errors.Wrap(err, "failed to create user")
	}
	return nil
}
func (s UserStore) Update(ctx context.Context, u *User) error {
	_, err := s.client.
		Database(s.config.Name).
		Collection(s.config.Collections.Users).
		UpdateOne(ctx, u.ID, u)

	if err != nil {
		return errors.Wrap(err, "failed to update user")
	}
	return nil
}

func (s UserStore) Delete(ctx context.Context, userID string) error {
	_, err := s.client.
		Database(s.config.Name).
		Collection(s.config.Collections.Users).
		DeleteOne(ctx, bson.D{{"_id", userID}})
	if err != nil {
		return errors.Wrap(err, "failed to delete user")
	}
	return nil
}
func (s UserStore) Get(ctx context.Context, userID string) (User, error) {
	user := User{}
	err := s.client.
		Database(s.config.Name).
		Collection(s.config.Collections.Users).
		FindOne(ctx, bson.D{{"_id", userID}}).
		Decode(user)
	if err != nil {
		return User{}, errors.Wrap(err, "failed to get user")
	}
	return user, nil
}
func (s UserStore) List(ctx context.Context, criteria ListCriteria) ([]User, error) {
	var users []User
	var filter bson.D
	var o *options.FindOptions
	if criteria.Country != "" {
		filter = bson.D{{"country", bson.D{{"$eq", criteria.Country}}}}
	} else {
		filter = bson.D{}
	}
	if criteria.Page > 0 && criteria.Limit > 0 {
		skip := criteria.Page*criteria.Limit - criteria.Limit
		o = &options.FindOptions{Limit: &criteria.Limit, Skip: &skip}
	} else {
		o = &options.FindOptions{}
	}
	cursor, err := s.client.
		Database(s.config.Name).
		Collection(s.config.Collections.Users).
		Find(ctx, filter, o)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list users")
	}
	defer cursor.Close(ctx)
	cursorErr := cursor.All(ctx, &users)
	if cursorErr != nil {
		return nil, errors.Wrap(cursorErr, "failed to get users from cursor")
	}

	return users, nil
}
