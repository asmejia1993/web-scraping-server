package scraper

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/asmejia1993/web-scraping-server/pkg/domain/hotel-franchises/model"
	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
	"github.com/sirupsen/logrus"
)

type Scraper struct {
	logger *logrus.Logger
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
		logger: lg,
	}
}

func (s *Scraper) InitScraping(req model.FranchiseScraper) model.SiteRes {
	d := strings.Split(req.Franchise.URL, ".")
	size := len(d)
	url := strings.Join(d[1:size], ".")
	res, err := whois.Whois(url)
	if err != nil {
		s.logger.Errorf("error invoking whois api: %v", err)
		return model.SiteRes{}
	}
	parsed, err := whoisparser.Parse(res)
	if err != nil {
		s.logger.Errorf("error parsing whois details: %v", err)
		return model.SiteRes{}
	}
	sslLabs := sslLabs{}
	isValid := s.validateWebsite(req.Franchise.URL, &sslLabs)
	protocol := sslLabs.Protocol

	hostNames := make([]string, len(sslLabs.Endpoints))
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
	}
	return site
}

func (s *Scraper) validateWebsite(url string, sLabs *sslLabs) bool {
	httpService := NewService(s.logger)
	formattedUrl := fmt.Sprintf(SSLLABS_URL, url)
	result, err := httpService.CallAPI(formattedUrl)
	if err != nil {
		return false
	}

	err = json.Unmarshal(result.Payload, &sLabs)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
