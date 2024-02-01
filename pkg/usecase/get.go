package usecase

import (
	"context"

	"github.com/asmejia1993/web-scraping-server/pkg/domain/hotel-franchises/model"
)

func (s Service) Get(id string, ctx context.Context) model.FranchiseInfo {
	return s.repo.FindFranchisesById(id, ctx)
}
