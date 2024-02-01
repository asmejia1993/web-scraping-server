package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

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

		//Create franchise size
		// key: id.franchise.size
		//value: n
		key := fmt.Sprintf("%s.franchises.size", response.ID)
		value := len(fr.Company.Franchises)
		hf.logger.Infof("creating franchise size in redis with key: %s, value: %d", key, value)
		hf.fService.SetKey(r.Context(), key, value)

		//Create franchise processed
		// key: id.franchise.processed
		//value: 0
		key = fmt.Sprintf("%s.franchises.processed", response.ID)
		value = 0
		hf.logger.Infof("creating franchise processed in redis with key: %s, value: %d", key, value)
		hf.fService.SetKey(r.Context(), key, value)

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
	//updating franchise processed
	err := hf.fService.Upsert(ctx, res)
	if err != nil {
		hf.logger.Errorf("error calling upsert handler: %v", err)
		return err
	}
	key := fmt.Sprintf("%s.franchises.processed", res.Id)
	val, err := hf.fService.GetKey(context.TODO(), key)
	if err != nil {
		hf.logger.Errorf("error updating cache with key: %s, error: %v", key, err)
	} else {
		count, _ := strconv.Atoi(val)
		hf.fService.SetKey(context.TODO(), key, count+1)
	}
	return nil
}
