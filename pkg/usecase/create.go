package usecase

import (
	"context"

	"github.com/asmejia1993/web-scraping-server/pkg/domain/hotel-franchises/model"
)

func (s Service) Create(ctx context.Context, req model.FranchiseInfo) (string, error) {
	return s.repo.CreateFranchisesHotel(ctx, req)
}
