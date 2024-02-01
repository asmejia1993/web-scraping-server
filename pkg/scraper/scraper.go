package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/asmejia1993/web-scraping-server/pkg/domain/hotel-franchises/model"
	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
	"github.com/sirupsen/logrus"
)

const RETRY_TIME = 3

type Scraper struct {
	logger       *logrus.Logger
	httpService  *Service
	sslLabsMutex sync.Mutex
}

const SSLLABS_URL = "https://api.ssllabs.com/api/v3/analyze?host=%s"

type sslLabs struct {
	Host            string     `json:"host,omitempty"`
	Port            int64      `json:"port,omitempty"`
	Protocol        string     `json:"protocol,omitempty"`
	IsPublic        bool       `json:"isPublic,omitempty"`
	Status          string     `json:"status,omitempty"`
	StatusMessage   string     `json:"statusMessage,omitempty"`
	StartTime       int64      `json:"startTime,omitempty"`
	TestTime        int64      `json:"testTime,omitempty"`
	EngineVersion   string     `json:"engineVersion,omitempty"`
	CriteriaVersion string     `json:"criteriaVersion,omitempty"`
	Endpoints       []endpoint `json:"endpoints,omitempty"`
}

type endpoint struct {
	IPAddress         string `json:"ipAddress"`
	ServerName        string `json:"serverName"`
	StatusMessage     string `json:"statusMessage"`
	Grade             string `json:"grade"`
	GradeTrustIgnored string `json:"gradeTrustIgnored"`
	HasWarnings       bool   `json:"hasWarnings"`
	IsExceptional     bool   `json:"isExceptional"`
	Progress          int64  `json:"progress"`
	Duration          int64  `json:"duration"`
	Delegation        int64  `json:"delegation"`
}

func NewScraperTask(lg *logrus.Logger) Scraper {
	return Scraper{
		logger:      lg,
		httpService: NewService(lg),
	}
}

var whoisMutex sync.Mutex

func WhoisWithMutex(domain string) (string, error) {
	whoisMutex.Lock()
	defer whoisMutex.Unlock()
	return whois.Whois(domain)
}

func (s *Scraper) InitScraping(ctx context.Context, req model.FranchiseScraper) model.SiteRes {
	arrs := strings.Split(req.Franchise.URL, ".")
	size := len(arrs)
	url := strings.Join(arrs[1:size], ".")
	res, err := WhoisWithMutex(url)
	if err != nil {
		s.logger.Errorf("error invoking whois api: %v", err)
		return model.SiteRes{}
	}
	parsed, err := whoisparser.Parse(res)
	if err != nil {
		s.logger.Errorf("error parsing whois details: %v", err)
		return model.SiteRes{}
	}
	s.sslLabsMutex.Lock()
	defer s.sslLabsMutex.Unlock()
	sslLabs := sslLabs{}
	isValid := s.validateWebsite(url, &sslLabs, ctx)
	protocol := sslLabs.Protocol

	hostNames := make([]string, 0)
	for _, v := range sslLabs.Endpoints {
		hostNames = append(hostNames, v.ServerName)
	}

	site := model.SiteRes{
		Id:          req.Id,
		Protocol:    protocol,
		Step:        len(hostNames),
		ServerNames: hostNames,
		CreatedAt:   parsed.Domain.CreatedDate,
		ExpiresAt:   parsed.Domain.ExpirationDate,
		Registrant:  parsed.Registrant.Name,
		Email:       parsed.Registrant.Email,
		IsValid:     isValid,
		Franchise: model.FranchiseReq{
			Name: req.Franchise.Name,
			URL:  req.Franchise.URL,
			Location: model.LocationReq{
				City:    req.Franchise.Location.City,
				Country: req.Franchise.Location.Country,
				Address: req.Franchise.Location.Address,
				ZipCode: req.Franchise.Location.ZipCode,
			},
		},
	}
	return site
}

func (s *Scraper) validateWebsite(url string, sLabs *sslLabs, ctx context.Context) bool {
	formattedUrl := fmt.Sprintf(SSLLABS_URL, url)
	result, err := s.httpService.CallAPIWithRetry(ctx, formattedUrl, RETRY_TIME)
	if err != nil {
		return false
	}

	err = json.Unmarshal(result.Payload, &sLabs)
	if err != nil {
		s.logger.Errorf("error unmarshalling JSON response: %v", err)
		return false
	}
	return true
}
