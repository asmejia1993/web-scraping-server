package service

import (
	"context"
	"fmt"
	"log"
	"strings"

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
	//var results []model.FranchiseInfo
	results := make([]model.FranchiseInfo, 0)
	criteria := buildCriteria(params)
	filter, err := bson.Marshal(criteria)
	fmt.Println(criteria)
	fmt.Println(filter)
	if err != nil {
		return results, err
	}

	var decodedDoc bson.M
	err = bson.Unmarshal(filter, &decodedDoc)
	if err != nil {
		log.Fatal("Error decoding BSON data:", err)
	}

	// Print the decoded BSON document
	fmt.Println(decodedDoc)

	//fil := bson.M{"company.owner.first_name": "jorge"}
	cursor, err := f.db.Client.Database(dbName).Collection(coll).Find(context.TODO(), filter)
	if err != nil {
		return results, fmt.Errorf("error invoking find: %v", err)
	}
	defer cursor.Close(ctx)

	if err = cursor.All(context.TODO(), &results); err != nil {
		return results, fmt.Errorf("error in cursor find: %v", err)
	}
	return results, nil
}

func buildCriteria(queryParams map[string][]string) bson.M {
	criteria := bson.M{}

	for key, values := range queryParams {
		if path, ok := AcceptedQueryParams[key]; ok {
			for _, value := range values {
				// Split the path into individual fields
				keys := strings.Split(path, ".")
				currentField := criteria

				// Traverse the map to set the value at the appropriate nested level
				for i, k := range keys {
					// If it's the last key, set the value
					if i == len(keys)-1 {
						currentField[k] = value
					} else {
						// If the key doesn't exist, create a nested map
						if _, ok := currentField[k]; !ok {
							currentField[k] = bson.M{}
						}
						currentField = currentField[k].(bson.M)
					}
				}
			}
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

	//... and so on!
}
