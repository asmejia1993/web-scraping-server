package service

import (
	"context"
	"fmt"

	"github.com/asmejia1993/web-scraping-server/pkg/config"
	"github.com/asmejia1993/web-scraping-server/pkg/domain/hotel-franchises/model"
	"github.com/asmejia1993/web-scraping-server/pkg/domain/hotel-franchises/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (f *franchiseRepository) UpSertFranchiseSite(ctx context.Context, site model.SiteRes) error {
	dbName := f.db.DBName
	coll := model.Collection

	objectId, _ := primitive.ObjectIDFromHex(site.Id)
	filter := bson.M{
		"_id": objectId,
	}

	franchise := model.Franchise{
		Name: site.Franchise.Name,
		URL:  site.Franchise.URL,
		Location: model.Location{
			City:    site.Franchise.Location.City,
			Country: site.Franchise.Location.Country,
			Address: site.Franchise.Location.Address,
			ZipCode: site.Franchise.Location.ZipCode,
		},
		Site: model.Site{
			Protocol:    site.Protocol,
			Step:        site.Step,
			ServerNames: site.ServerNames,
			CreatedAt:   site.CreatedAt,
			ExpiresAt:   site.ExpiresAt,
			Registrant:  site.Registrant,
			Email:       site.Email,
		},
	}

	update := bson.M{
		"$set": bson.M{
			"company.franchises.$[x].site": franchise.Site,
		},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.D{
				{Key: "x.name", Value: franchise.Name},
				{Key: "x.url", Value: franchise.URL},
				{Key: "x.location.city", Value: franchise.Location.City},
				{Key: "x.location.country", Value: franchise.Location.Country},
			},
		},
	}

	opts := options.Update().
		SetUpsert(true).
		SetArrayFilters(arrayFilters)

	_, err := f.db.Client.Database(dbName).Collection(coll).UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		fmt.Printf("error executing site update with: %v", err)
		return err
	}
	return nil
}

func (f *franchiseRepository) All(ctx context.Context, params map[string][]string) ([]model.FranchiseInfo, error) {
	dbName := f.db.DBName
	coll := model.Collection

	results := make([]model.FranchiseInfo, 0)

	criteria := buildCriteria(params)

	cursor, err := f.db.Client.Database(dbName).Collection(coll).Find(ctx, criteria)
	if err != nil {
		return nil, fmt.Errorf("error invoking find: %v", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var franchise model.FranchiseInfo
		if err := cursor.Decode(&franchise); err != nil {
			return nil, fmt.Errorf("error decoding document: %v", err)
		}
		results = append(results, franchise)
	}
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error in cursor: %v", err)
	}

	return results, nil
}

func buildCriteria(queryParams map[string][]string) bson.M {
	criteria := bson.M{}
	for key, path := range AcceptedQueryParams {
		if values, ok := queryParams[key]; ok && len(values) > 0 {
			criteria[path] = values[0]
		}
	}
	return criteria
}

var AcceptedQueryParams = map[string]string{
	"owner_first_name":   "company.owner.first_name",
	"owner_last_name":    "company.owner.last_name",
	"owner_email":        "company.owner.contact.email",
	"company_name":       "company.information.name",
	"franchise_location": "company.franchises.location.city",
	"franchise_name":     "company.franchises.name",
	"franchise_url":      "company.franchises.url",
	"company_tax_number": "company.information.tax_number",
}
