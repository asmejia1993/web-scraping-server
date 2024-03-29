package rest

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/asmejia1993/web-scraping-server/pkg/config"
	"github.com/asmejia1993/web-scraping-server/pkg/http/rest/handler"
	"github.com/asmejia1993/web-scraping-server/pkg/scraper"
	"github.com/asmejia1993/web-scraping-server/pkg/workerpool"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

const (
	WORKER_THREAD = 10
	BUFFER        = 1000
)

type Server struct {
	logger *logrus.Logger
	router *mux.Router
	config *config.AppConfig
	worker *workerpool.WorkerPool
}

func NewServer(ctx context.Context) (*Server, error) {

	//Initialize
	appConfig := config.LoadConfig()
	log := NewLogger()
	router := mux.NewRouter()
	scraper := scraper.NewScraperTask(log)
	worker := workerpool.NewWorkerPool(WORKER_THREAD, BUFFER, log, &scraper, ctx)

	s := Server{
		logger: log,
		router: router,
		config: appConfig,
		worker: worker,
	}
	hf := handler.NewHandler(log, appConfig, s.worker, ctx)
	handler.Register(router, &hf, ctx)
	return &s, nil
}

func (s *Server) Run(ctx context.Context) error {
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.HTTPInfo.Port),
		Handler: cors.Default().Handler(s.router),
	}

	stopServer := make(chan os.Signal, 1)
	signal.Notify(stopServer, syscall.SIGINT, syscall.SIGTERM)

	defer signal.Stop(stopServer)

	// channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		s.logger.Printf("WebScraping listening on port: %d", s.config.HTTPInfo.Port)
		serverErrors <- server.ListenAndServe()
	}(&wg)

	// shutdown.
	select {
	case err := <-serverErrors:
		return fmt.Errorf("error: starting REST API server: %w", err)
	case <-ctx.Done():
	case <-stopServer:
		s.logger.Warn("server received STOP signal")
		s.config.CloseMongoDB(ctx)
		s.worker.Stop()
		err := server.Shutdown(ctx)
		if err != nil {
			return fmt.Errorf("graceful shutdown did not complete: %w", err)
		}
		wg.Wait()
		s.logger.Info("server was shut down gracefully")
	}
	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
