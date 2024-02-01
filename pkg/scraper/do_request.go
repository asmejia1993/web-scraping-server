package scraper

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type Service struct {
	httpClient *http.Client
	logger     *logrus.Logger
	limiter    *rate.Limiter
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
		limiter:    rate.NewLimiter(10, 1),
	}
}

type CallResponse struct {
	Payload []byte
}

func (s *Service) CallAPI(ctx context.Context, url string) (CallResponse, error) {
	start := time.Now()

	response := CallResponse{}
	// Apply rate limiting
	if err := s.limiter.Wait(ctx); err != nil {
		return response, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return response, err
	}

	res, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Errorf("call external service returned: %v\n", err)
		return response, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer res.Body.Close()

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return response, fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle non-successful status codes
	if res.StatusCode >= 400 && res.StatusCode < 500 {
		s.logger.Errorf("external service returned bad response. Status Code: %d. Content: %s\n", res.StatusCode, string(content))
		return response, errors.New("bad response")
	}

	duration := time.Since(start).Seconds()
	s.logger.Infof("API call duration: %v seconds, body: %v\n", duration, string(content))
	return CallResponse{Payload: content}, nil
}

func (s *Service) CallAPIWithRetry(ctx context.Context, url string, maxRetries int) (CallResponse, error) {
	var resp CallResponse
	var err error

	// Retry with exponential backoff and jitter
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 0 // Retry indefinitely
	b.MaxInterval = 10 * time.Second
	b.RandomizationFactor = 0.5

	err = backoff.RetryNotify(func() error {
		// Call the API
		resp, err = s.CallAPI(ctx, url)
		if err != nil {
			s.logger.Printf("API call failed: %v\n", err)
		}
		return err
	}, b, func(err error, t time.Duration) {
		s.logger.Printf("Retry attempt %d failed. Retrying in %v...\n", 3, t)
	})

	if err != nil {
		return resp, fmt.Errorf("exceeded max retry attempts: %w", err)
	}

	return resp, nil
}
