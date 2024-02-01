package handler

import (
	"context"

	"github.com/asmejia1993/web-scraping-server/pkg/config"
	"github.com/asmejia1993/web-scraping-server/pkg/domain/hotel-franchises/model"
	fr "github.com/asmejia1993/web-scraping-server/pkg/domain/hotel-franchises/service"
	fs "github.com/asmejia1993/web-scraping-server/pkg/usecase"
	"github.com/asmejia1993/web-scraping-server/pkg/workerpool"
	"github.com/sirupsen/logrus"
)

const (
	WORKER_THREAD = 10
	BUFFER        = 100
)

type handlerFranchises struct {
	logger     *logrus.Logger
	fService   fs.Service
	workerPool *workerpool.WorkerPool
	resultChan <-chan model.SiteRes
}

func NewHandler(lg *logrus.Logger, db *config.DBInfo, worker *workerpool.WorkerPool, ctx context.Context) handlerFranchises {

	worker.Start(ctx)

	return handlerFranchises{
		logger:     lg,
		fService:   fs.NewService(fr.NewFranchiseRepository(db)),
		workerPool: worker,
		resultChan: worker.GetResultChan(),
	}
}
