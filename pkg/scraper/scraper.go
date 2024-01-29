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
	host            string     `json:"host"`
	port            int64      `json:"port"`
	protocol        string     `json:"protocol"`
	isPublic        bool       `json:"isPublic"`
	status          string     `json:"status"`
	startTime       int64      `json:"startTime"`
	testTime        int64      `json:"testTime"`
	engineVersion   string     `json:"engineVersion"`
	criteriaVersion string     `json:"criteriaVersion"`
	endpoints       []endpoint `json:"endpoints,omitempty"`
}

type endpoint struct {
	iPAddress         string `json:"ipAddress"`
	serverName        string `json:"serverName"`
	statusMessage     string `json:"statusMessage"`
	grade             string `json:"grade"`
	gradeTrustIgnored string `json:"gradeTrustIgnored"`
	hasWarnings       bool   `json:"hasWarnings"`
	isExceptional     bool   `json:"isExceptional"`
	progress          int64  `json:"progress"`
	duration          int64  `json:"duration"`
	delegation        int64  `json:"delegation"`
}

type Domain struct {
	Name     string
	URL      string
	Location struct {
		City    string
		Country string
		Address string
		ZipCode string
	}
}

func NewScraperTask(lg *logrus.Logger) Scraper {
	return Scraper{
		logger: lg,
	}
}

func (s *Scraper) InitScraping(domains []model.Franchise) {
	for _, domain := range domains {
		s.processEachDomain(domain)
	}
}

func (s *Scraper) processEachDomain(domain model.Franchise) {
	d := strings.Split(domain.URL, ".")
	size := len(d)
	url := strings.Join(d[1:size], ".")
	res, err := whois.Whois(url)
	if err != nil {
		s.logger.Errorf("error invoking whois api: %v", err)
		return
	}
	parsed, err := whoisparser.Parse(res)
	if err != nil {
		s.logger.Errorf("error parsing whois details: %v", err)
		return
	}
	sslLabs := sslLabs{}
	isValid := validateWebsite(domain.URL, &sslLabs)
	protocol := determineProtocol(domain.URL, &sslLabs)

	s.logger.Infof("domainStatus: %v, registrant: %v - createdAt: %v - expiredAt: %v - emailContact: %s", parsed.Domain.Status, parsed.Registrant.Name, parsed.Domain.CreatedDate, parsed.Domain.ExpirationDate, parsed.Registrant.Email)
	s.logger.Infof("protocol: %v - isValid: %v\n", protocol, isValid)
}

func determineProtocol(url string, sLabs *sslLabs) string {
	return ""
}

func validateWebsite(url string, sLabs *sslLabs) bool {
	s := NewService()
	formattedUrl := fmt.Sprintf(SSLLABS_URL, url)
	result, err := s.CallAPI(formattedUrl)
	if err != nil {
		return false
	}

	err = json.Unmarshal([]byte(result.Payload), &sLabs)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
