package workerpool

import (
	"context"
	"errors"
	"sync"

	"github.com/asmejia1993/web-scraping-server/pkg/domain/hotel-franchises/model"
	"github.com/asmejia1993/web-scraping-server/pkg/scraper"
	"github.com/sirupsen/logrus"
)

var ErrWorkerBusy = errors.New("workers are busy, try again later")

type WorkerPool struct {
	workers         []*worker
	jobQueue        chan model.FranchiseScraper
	resultChan      chan model.SiteRes
	logger          *logrus.Logger
	scraper         *scraper.Scraper
	jobQueueMutex   sync.Mutex
	resultChanMutex sync.Mutex
}

type worker struct {
	id              int
	pool            *WorkerPool
	cancel          context.CancelFunc
	wg              *sync.WaitGroup
	logger          *logrus.Logger
	resultChanMutex sync.Mutex
	scraper         *scraper.Scraper
}

func NewWorkerPool(workerCount, buffer int, log *logrus.Logger, s *scraper.Scraper, ctx context.Context) *WorkerPool {
	wp := &WorkerPool{
		jobQueue:   make(chan model.FranchiseScraper, buffer),
		resultChan: make(chan model.SiteRes, buffer),
		logger:     log,
		scraper:    s,
		workers:    make([]*worker, 0),
	}
	//_, c := context.WithCancel(ctx)
	for i := 0; i < workerCount; i++ {
		w := &worker{
			id:              i,
			pool:            wp,
			wg:              &sync.WaitGroup{},
			logger:          log,
			scraper:         s,
			resultChanMutex: sync.Mutex{},
			//cancel:          &c,
		}
		wp.workers = append(wp.workers, w)
	}
	return wp
}

func (wp *WorkerPool) Start(ctx context.Context) {
	for _, w := range wp.workers {
		w.wg.Add(1)
		go w.start(ctx)
	}
}

func (wp *WorkerPool) Stop() {
	for _, w := range wp.workers {
		if w.cancel != nil {
			defer w.cancel()
			defer w.wg.Wait()
		}
	}
	if wp.jobQueue != nil {
		close(wp.jobQueue)
	}
	if wp.resultChan != nil {
		close(wp.resultChan)
	}
}

func (wp *WorkerPool) QueueTask(task model.FranchiseScraper) error {
	// Lock the mutex to synchronize access to the jobQueue field
	wp.jobQueueMutex.Lock()
	defer wp.jobQueueMutex.Unlock()

	select {
	case wp.jobQueue <- task:
		return nil
	default:
		return errors.New("worker pool is busy, try again later")
	}
}

func (wp *WorkerPool) GetResultChan() <-chan model.SiteRes {
	// Lock the mutex to synchronize access to the resultChan field
	wp.resultChanMutex.Lock()
	defer wp.resultChanMutex.Unlock()

	return wp.resultChan
}

func (w *worker) start(ctx context.Context) {
	defer w.wg.Done()
	//ctx, w.cancel = context.WithCancel(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-w.pool.jobQueue:
			if !ok {
				return
			}
			w.doWork(ctx, task)
		}
	}
}

func (w *worker) doWork(ctx context.Context, task model.FranchiseScraper) {
	defer func() {
		if r := recover(); r != nil {
			w.logger.Error("panic occurred in worker:", r)
		}
	}()
	w.logger.WithField("task", task.Franchise.URL).Info("start scraping ...")

	w.resultChanMutex.Lock()
	defer w.resultChanMutex.Unlock()
	select {
	case w.pool.resultChan <- w.scraper.InitScraping(ctx, task):
	case <-ctx.Done():
		return
	default:
		w.logger.Warn("ResultChan is full, dropping result")
	}
	w.logger.WithField("task", task.Franchise.URL).Info("work completed!")
}
