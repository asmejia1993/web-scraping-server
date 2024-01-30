package handler

import (
	"context"
	"net/http"

	"github.com/asmejia1993/web-scraping-server/pkg/domain/hotel-franchises/model"
	"github.com/asmejia1993/web-scraping-server/pkg/worker"
)

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
		hf.worker.Start(r.Context(), len(req.Company.Franchises))
		for _, v := range req.Company.Franchises {
			scrapReq := model.FranchiseScraper{
				Id:        response.ID,
				Franchise: v,
			}

			if err := hf.worker.QueueTask(scrapReq); err != nil {
				hf.logger.WithError(err).Info("failed to queue task")
				if err == worker.ErrWorkerBusy {
					w.Header().Set("Retry-After", "60")
					hf.respond(w, `{"error": "workers are busy, try again later"}`, http.StatusServiceUnavailable)
					return
				}
				hf.respond(w, `{"error": "failed to queue task"}`, http.StatusInternalServerError)
				return
			}
		}

		hf.respond(w, response, http.StatusAccepted)

		go hf.receiveResultFromWorker(r.Context())
	}
}

func (hf handlerFranchises) receiveResultFromWorker(ctx context.Context) {
	for res := range hf.resultChan {
		select {
		case <-ctx.Done():
			return
		default:
			hf.logger.Infof("Received result from worker: %v", res)
		}
	}
}
