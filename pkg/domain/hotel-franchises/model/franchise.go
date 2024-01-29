package model

type Contact struct {
	Email    string   `json:"email" bson:"email"`
	Phone    string   `json:"phone" bson:"phone"`
	Location Location `json:"location" bson:"location"`
}

type Location struct {
	City    string `json:"city" bson:"city"`
	Country string `json:"country" bson:"country"`
	Address string `json:"address" bson:"address"`
	ZipCode string `json:"zip_code" bson:"zip_code"`
}

type Owner struct {
	FirstName string  `json:"first_name" bson:"first_name"`
	LastName  string  `json:"last_name" bson:"last_name"`
	Contact   Contact `json:"contact" bson:"contact"`
}

type Franchise struct {
	Name     string   `json:"name" bson:"name"`
	URL      string   `json:"url" bson:"url"`
	Location Location `json:"location" bson:"location"`
}

type Information struct {
	Name      string   `json:"name" bson:"name"`
	TaxNumber string   `json:"tax_number" bson:"tax_number"`
	Location  Location `json:"location" bson:"location"`
}

type Company struct {
	Owner       Owner       `json:"owner" bson:"owner"`
	Information Information `json:"information" bson:"information"`
	Franchises  []Franchise `json:"franchises" bson:"franchises"`
}

type FranchiseInfo struct {
	ID      string  `bson:"_id,omitempty" json:"id"`
	Company Company `json:"company" bson:"company"`
}

const Collection string = "franchises_hotel"
