package handler

import (
	"github.com/asmejia1993/web-scraping-server/pkg/config"
	"github.com/asmejia1993/web-scraping-server/pkg/domain/hotel-franchises/model"
	fr "github.com/asmejia1993/web-scraping-server/pkg/domain/hotel-franchises/service"
	fs "github.com/asmejia1993/web-scraping-server/pkg/usecase"
	"github.com/asmejia1993/web-scraping-server/pkg/worker"
	"github.com/sirupsen/logrus"
)

const (
	WORKER_THREAD = 10
	BUFFER        = 100
)

type handlerFranchises struct {
	logger     *logrus.Logger
	fService   fs.Service
	worker     worker.IWorker
	resultChan chan model.SiteRes
}

func newHandler(lg *logrus.Logger, db *config.DBInfo, w worker.IWorker) handlerFranchises {
	return handlerFranchises{
		logger:     lg,
		fService:   fs.NewService(fr.NewFranchiseRepository(db)),
		worker:     w,
		resultChan: worker.ResultChan,
	}
}
