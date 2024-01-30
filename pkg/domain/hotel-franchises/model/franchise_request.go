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

func ConvertReqToFranchiseInfo(req CompanyReq) FranchiseInfo {
	var franchiseInfo FranchiseInfo

	// Convert OwnerReq to Owner
	owner := Owner{
		FirstName: req.Owner.FirstName,
		LastName:  req.Owner.LastName,
		Contact: Contact{
			Email: req.Owner.Contact.Email,
			Phone: req.Owner.Contact.Phone,
			Location: Location{
				City:    req.Owner.Contact.Location.City,
				Country: req.Owner.Contact.Location.Country,
				Address: req.Owner.Contact.Location.Address,
				ZipCode: req.Owner.Contact.Location.ZipCode,
			},
		},
	}

	// Convert InformationReq to Information
	information := Information{
		Name:      req.Information.Name,
		TaxNumber: req.Information.TaxNumber,
		Location: Location{
			City:    req.Information.Location.City,
			Country: req.Information.Location.Country,
			Address: req.Information.Location.Address,
			ZipCode: req.Information.Location.ZipCode,
		},
	}

	// Convert FranchiseReq to Franchise
	var franchises []Franchise
	for _, fr := range req.Franchises {
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
