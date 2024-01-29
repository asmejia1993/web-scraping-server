package handler

import (
	"github.com/asmejia1993/web-scraping-server/pkg/config"
	fr "github.com/asmejia1993/web-scraping-server/pkg/domain/hotel-franchises/service"
	"github.com/asmejia1993/web-scraping-server/pkg/scraper"
	fs "github.com/asmejia1993/web-scraping-server/pkg/usecase"
	"github.com/sirupsen/logrus"
)

type handlerFranchises struct {
	logger   *logrus.Logger
	fService fs.Service
	sc       scraper.Scraper
}

func newHandler(lg *logrus.Logger, db *config.DBInfo) handlerFranchises {
	return handlerFranchises{
		logger:   lg,
		fService: fs.NewService(fr.NewFranchiseRepository(db)),
		sc:       scraper.NewScraperTask(lg),
	}
}
