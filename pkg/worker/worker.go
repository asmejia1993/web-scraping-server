package worker

import (
	"context"
	"errors"
	"sync"

	"github.com/asmejia1993/web-scraping-server/pkg/domain/hotel-franchises/model"
	"github.com/asmejia1993/web-scraping-server/pkg/scraper"
	"github.com/sirupsen/logrus"
)

var ResultChan = make(chan model.SiteRes)

type worker struct {
	workchan    chan workType
	workerCount int
	buffer      int
	wg          *sync.WaitGroup
	cancelFunc  context.CancelFunc
	logger      *logrus.Logger
	scraper     scraper.Scraper
}

type IWorker interface {
	Start(ctx context.Context, size int)
	Stop()
	QueueTask(task model.FranchiseScraper) error
}

func (w *worker) SetBuffer(size int) {
	w.buffer = size
}
func New(workerCount, buffer int, log *logrus.Logger) IWorker {
	w := worker{
		//workchan:    make(chan workType, buffer),
		workerCount: workerCount,
		buffer:      buffer,
		wg:          new(sync.WaitGroup),
		logger:      log,
		scraper:     scraper.NewScraperTask(log),
	}
	return &w
}

func (w *worker) Start(ctx context.Context, size int) {
	ctx, cancelFunc := context.WithCancel(ctx)
	w.cancelFunc = cancelFunc
	w.workchan = make(chan workType, size)

	for i := 0; i < size; i++ {
		w.wg.Add(1)
		go w.spawnWorkers(ctx)
	}
}

func (w *worker) Stop() {
	w.logger.Info("stop workers")
	close(w.workchan)
	w.cancelFunc()
	w.wg.Wait()
	w.logger.Info("all workers exited!")
}

func (w *worker) QueueTask(task model.FranchiseScraper) error {
	if len(w.workchan) >= w.buffer {
		return ErrWorkerBusy
	}
	w.workchan <- workType{Task: task}
	return nil
}

func (w *worker) spawnWorkers(ctx context.Context) {
	defer w.wg.Done()

	for work := range w.workchan {
		select {
		case <-ctx.Done():
			return
		default:
			w.doWork(ctx, work.Task)
		}
	}
}

func (w *worker) doWork(ctx context.Context, task model.FranchiseScraper) {
	w.logger.WithField("task", task.Franchise.URL).Info("start scraping ...")

	ResultChan <- w.scraper.InitScraping(task)

	w.logger.WithField("task", task).Info("work completed!")
}

type workType struct {
	Task model.FranchiseScraper
}

var (
	ErrWorkerBusy = errors.New("workers are busy, try again later")
)