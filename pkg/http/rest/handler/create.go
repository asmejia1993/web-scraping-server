package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/asmejia1993/web-scraping-server/pkg/domain/hotel-franchises/model"
)

var ErrWorkerBusy = errors.New("workers are busy, try again later")

func (hf handlerFranchises) Create() http.HandlerFunc {

	type inserted struct {
		ID string `json:"id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req model.FranchiseInfoReq

		err := hf.decode(r, &req)
		if err != nil {
			hf.respond(w, err, http.StatusInternalServerError)
			return
		}
		fr := model.ConvertReqToFranchiseInfo(req)
		res, err := hf.fService.Create(r.Context(), fr)
		response := inserted{
			ID: res,
		}
		hf.logger.Infof("response: %v", response)
		if err != nil {
			hf.respond(w, err, http.StatusInternalServerError)
			return
		}

		//Queue new task
		for _, v := range req.Company.Franchises {
			scrapReq := model.FranchiseScraper{
				Id:        response.ID,
				Franchise: v,
			}

			if err := hf.workerPool.QueueTask(scrapReq); err != nil {
				hf.logger.WithError(err).Info("failed to queue task")
				if err == ErrWorkerBusy {
					w.Header().Set("Retry-After", "60")
					hf.respond(w, `{"error": "workers are busy, try again later"}`, http.StatusServiceUnavailable)
					return
				}
				hf.respond(w, `{"error": "failed to queue task"}`, http.StatusInternalServerError)
				return
			}
		}

		hf.respond(w, response, http.StatusAccepted)

		go func() {
			for res := range hf.workerPool.GetResultChan() {
				hf.logger.Infof("site received: %v", res)
				if err := hf.processResult(r.Context(), res); err != nil {
					hf.logger.Errorf("error processing result from worker: %v", err)
				} else {
					hf.logger.Infof("updated franchise with: %v", res)
				}
			}
		}()
	}
}

func (hf handlerFranchises) processResult(ctx context.Context, res model.SiteRes) error {
	if err := hf.fService.Upsert(ctx, res); err != nil {
		hf.logger.Errorf("error calling upsert handler: %v", err)
		return err
	}
	return nil
}
