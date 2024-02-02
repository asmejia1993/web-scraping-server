package model

type Contact struct {
	Email    string   `json:"email,omitempty" bson:"email"`
	Phone    string   `json:"phone,omitempty" bson:"phone"`
	Location Location `json:"location,omitempty" bson:"location"`
}

type Location struct {
	City    string `json:"city,omitempty" bson:"city"`
	Country string `json:"country,omitempty" bson:"country"`
	Address string `json:"address,omitempty" bson:"address"`
	ZipCode string `json:"zip_code,omitempty" bson:"zip_code"`
}

type Owner struct {
	FirstName string  `json:"first_name,omitempty" bson:"first_name"`
	LastName  string  `json:"last_name,omitempty" bson:"last_name"`
	Contact   Contact `json:"contact,omitempty" bson:"contact"`
}

type Franchise struct {
	Name     string   `json:"name,omitempty" bson:"name"`
	URL      string   `json:"url,omitempty" bson:"url"`
	Location Location `json:"location,omitempty" bson:"location"`
	Site     Site     `json:"site,omitempty" bson:"site"`
}

type Site struct {
	Protocol    string   `json:"protocol,omitempty" bson:"protocol"`
	Step        int      `json:"steps,omitempty" bson:"steps"`
	ServerNames []string `json:"server_names,omitempty" bson:"server_names"`
	CreatedAt   string   `json:"created_at,omitempty" bson:"created_at"`
	ExpiresAt   string   `json:"expires_at,omitempty" bson:"expires_at"`
	Registrant  string   `json:"registrant,omitempty" bson:"registrant"`
	Email       string   `json:"email_contact,omitempty" bson:"email_contact"`
}

type Information struct {
	Name      string   `json:"name,omitempty" bson:"name"`
	TaxNumber string   `json:"tax_number,omitempty" bson:"tax_number"`
	Location  Location `json:"location,omitempty" bson:"location"`
}

type Company struct {
	Owner       Owner       `json:"owner,omitempty" bson:"owner"`
	Information Information `json:"information,omitempty" bson:"information"`
	Franchises  []Franchise `json:"franchises,omitempty" bson:"franchises"`
}

type FranchiseInfo struct {
	ID      string  `bson:"_id" json:"id,omitempty"`
	Company Company `json:"company,omitempty" bson:"company"`
}

const Collection string = "franchises_hotel"

type Franchise_Request struct {
}
