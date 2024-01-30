package scraper

import (
	"fmt"
	"io"
	"math/rand"
	"net"
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
	netTransport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
	}
	httpClient := &http.Client{
		Timeout:   15 * time.Second,
		Transport: netTransport,
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

	start := time.Now()

	// Apply rate limiting
	s.limiter.Take()

	res, err := s.httpClient.Get(url)
	if err != nil {
		s.logger.Errorf("call external service returned: %v\n", err)
		return CallResponse{}, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer res.Body.Close()

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return CallResponse{}, fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle non-successful status codes
	if res.StatusCode >= 400 || res.StatusCode < 500 {
		s.logger.Errorf("external service returned bad response. Status Code: %d. Content: %s\n", res.StatusCode, string(content))
		return CallResponse{}, nil
	}

	duration := time.Since(start).Seconds()
	s.logger.Infof("API call duration: %v seconds, body: %v\n", duration, string(content))
	return CallResponse{Payload: content}, nil
}

func (s *Service) CallAPIWithRetry(url string, maxRetries int) (CallResponse, error) {
	var resp CallResponse
	var err error

	// Retry with exponential backoff and jitter
	err = retryWithBackoff(maxRetries, 2*time.Second, func() error {
		// Call the API
		resp, err = s.CallAPI(url)
		if err != nil {
			fmt.Printf("API call failed: %v\n", err)
		}
		return err
	})

	if err != nil {
		return CallResponse{}, fmt.Errorf("exceeded max retry attempts: %w", err)
	}

	return resp, nil
}

func retryWithBackoff(maxRetries int, initialDelay time.Duration, fn func() error) error {
	delay := initialDelay
	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}

		// Calculate the next delay using exponential backoff with jitter
		jitter := time.Duration(rand.Int63n(int64(delay / 2)))
		nextDelay := delay + jitter

		fmt.Printf("Retry attempt %d failed. Retrying in %v...\n", attempt, nextDelay)

		time.Sleep(nextDelay)
		delay *= 2 // Exponential backoff
	}
	return fmt.Errorf("exceeded max retry attempts")
}
