package usecase

import (
	"context"

	"github.com/asmejia1993/web-scraping-server/pkg/domain/hotel-franchises/model"
)

func (s Service) All(ctx context.Context, params map[string][]string) ([]model.FranchiseInfo, error) {
	return s.repo.All(ctx, params)

}
