package services

import (
	"amqpBackend/models"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccountModel struct {
	MongoCollection *mongo.Collection
}

func (r *AccountModel) CreateAccount(account *models.CollAccount) (interface{}, error) {
	result, err := r.MongoCollection.InsertOne(context.Background(), account)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func (r *AccountModel) UpdateAccountById(accountId int, updatedAccount *models.CollAccount) (int64, error) {
	result, err := r.MongoCollection.UpdateOne(context.Background(),
		bson.D{{Key: "AccountID", Value: accountId}},
		bson.D{{Key: "$set", Value: updatedAccount}})
	if err != nil {
		return 0, err
	}

	return result.ModifiedCount, nil
}

func (r *AccountModel) DeleteAccountById(accountId int) (int64, error) {
	result, err := r.MongoCollection.DeleteOne(context.Background(),
		bson.D{{Key: "AccountID", Value: accountId}})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

func (r *AccountModel) FindAccountById(accountId int) (*models.CollAccount, error) {
	var account models.CollAccount

	err := r.MongoCollection.FindOne(context.Background(),
		bson.D{{Key: "AccountID", Value: "accountId"}}).Decode(&account)

	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *AccountModel) FindAllAccount() ([]models.CollAccount, error) {
	result, err := r.MongoCollection.Find(context.Background(), bson.D{})

	if err != nil {
		return nil, err
	}

	var account []models.CollAccount
	err = result.All(context.Background(), &account)

	if err != nil {
		return nil, fmt.Errorf("resulst error decode %s", err.Error())
	}

	return account, nil
}
