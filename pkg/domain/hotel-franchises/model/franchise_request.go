package model

type ContactReq struct {
	Email    string      `json:"email"`
	Phone    string      `json:"phone"`
	Location LocationReq `json:"location"`
}

type LocationReq struct {
	City    string `json:"city"`
	Country string `json:"country"`
	Address string `json:"address"`
	ZipCode string `json:"zip_code"`
}

type OwnerReq struct {
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Contact   ContactReq `json:"contact"`
}

type FranchiseReq struct {
	Name     string      `json:"name"`
	URL      string      `json:"url"`
	Location LocationReq `json:"location"`
}

type SiteRes struct {
	Id          string   `json:"id"`
	Protocol    string   `json:"protocol"`
	Step        int      `json:"steps"`
	ServerNames []string `json:"server_names"`
	CreatedAt   string   `json:"created_at"`
	ExpiresAt   string   `json:"expires_at"`
	Registrant  string   `json:"registrant"`
	Email       string   `json:"email_contact"`
	IsValid     bool     `json:"is_valid"`
}

type InformationReq struct {
	Name      string      `json:"name"`
	TaxNumber string      `json:"tax_number"`
	Location  LocationReq `json:"location"`
}

type CompanyReq struct {
	Owner       Owner          `json:"owner"`
	Information InformationReq `json:"information"`
	Franchises  []FranchiseReq `json:"franchises"`
}

type FranchiseInfoReq struct {
	Company CompanyReq `json:"company"`
}

type FranchiseScraper struct {
	Id        string       `json:"id"`
	Franchise FranchiseReq `json:"franchise"`
}

func ConvertReqToFranchiseInfo(req FranchiseInfoReq) FranchiseInfo {
	var franchiseInfo FranchiseInfo

	// Convert OwnerReq to Owner
	owner := Owner{
		FirstName: req.Company.Owner.FirstName,
		LastName:  req.Company.Owner.LastName,
		Contact: Contact{
			Email: req.Company.Owner.Contact.Email,
			Phone: req.Company.Owner.Contact.Phone,
			Location: Location{
				City:    req.Company.Owner.Contact.Location.City,
				Country: req.Company.Owner.Contact.Location.Country,
				Address: req.Company.Owner.Contact.Location.Address,
				ZipCode: req.Company.Owner.Contact.Location.ZipCode,
			},
		},
	}

	// Convert InformationReq to Information
	information := Information{
		Name:      req.Company.Information.Name,
		TaxNumber: req.Company.Information.TaxNumber,
		Location: Location{
			City:    req.Company.Information.Location.City,
			Country: req.Company.Information.Location.Country,
			Address: req.Company.Information.Location.Address,
			ZipCode: req.Company.Information.Location.ZipCode,
		},
	}

	// Convert FranchiseReq to Franchise
	var franchises []Franchise
	for _, fr := range req.Company.Franchises {
		franchises = append(franchises, Franchise{
			Name: fr.Name,
			URL:  fr.URL,
			Location: Location{
				City:    fr.Location.City,
				Country: fr.Location.Country,
				Address: fr.Location.Address,
				ZipCode: fr.Location.ZipCode,
			},
		})
	}

	// Construct FranchiseInfo
	franchiseInfo = FranchiseInfo{
		Company: Company{
			Owner:       owner,
			Information: information,
			Franchises:  franchises,
		},
	}

	return franchiseInfo
}
