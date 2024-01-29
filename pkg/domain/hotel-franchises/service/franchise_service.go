package service

import (
	"context"

	"github.com/asmejia1993/web-scraping-server/pkg/config"
	"github.com/asmejia1993/web-scraping-server/pkg/domain/hotel-franchises/model"
	"github.com/asmejia1993/web-scraping-server/pkg/domain/hotel-franchises/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type franchiseRepository struct {
	db *config.DBInfo
}

func NewFranchiseRepository(db *config.DBInfo) repository.IFranchiseRepository {
	return &franchiseRepository{db: db}
}

func (f *franchiseRepository) FindFranchisesById(id string, ctx context.Context) model.FranchiseInfo {
	var franchisesInfo model.FranchiseInfo
	objID, _ := primitive.ObjectIDFromHex(id)
	dbName := f.db.DBName
	coll := model.Collection
	res := f.db.Client.Database(dbName).Collection(coll).FindOne(ctx, bson.M{"_id": objID})
	res.Decode(&franchisesInfo)
	return franchisesInfo
}

func (f *franchiseRepository) CreateFranchisesHotel(ctx context.Context, req model.FranchiseInfo) (string, error) {
	dbName := f.db.DBName
	coll := model.Collection
	res, err := f.db.Client.Database(dbName).Collection(coll).InsertOne(ctx, req)
	if err != nil {
		return "", err
	}
	id := res.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}
