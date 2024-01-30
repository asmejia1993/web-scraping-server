package scraper

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"go.uber.org/ratelimit"
)

type Service struct {
	httpClient *http.Client
	logger     *logrus.Logger
	limiter    ratelimit.Limiter
}

func NewService(log *logrus.Logger) *Service {
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	return &Service{
		httpClient: httpClient,
		logger:     log,
		limiter:    ratelimit.New(10),
	}
}

type CallResponse struct {
	Payload []byte
}

func (s *Service) CallAPI(url string) (CallResponse, error) {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatalf("Error Occurred. %+v", err)
	}

	start := time.Now()
	//s.limiter.Take()

	res, err := s.httpClient.Do(req)
	if err != nil {
		return CallResponse{}, fmt.Errorf("do request: %w", err)
	}

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return CallResponse{}, fmt.Errorf("read response body: %w", err)
	}

	// gracefully handle the bad responses
	if res.StatusCode >= 400 && res.StatusCode < 500 {
		s.logger.Errorf("external service returned bad response. Code: %d. Content: %s\n", res.StatusCode, string(content))
		return CallResponse{}, nil
	}

	spentSeconds := time.Since(start).Seconds()
	s.logger.Infof("call api duration: %v with response: %s", spentSeconds, string(content))
	return CallResponse{Payload: content}, nil
}
